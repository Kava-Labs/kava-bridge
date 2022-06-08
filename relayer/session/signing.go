package session

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	logging "github.com/ipfs/go-log/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("SigningSession")
var tracer = otel.Tracer("SigningSession")

// SigningSessionResult is the result of a signing session.
type SigningSessionResult struct {
	Signature *tss_common.SignatureData
	Err       *tss.Error
}

// SigningSession is a session for signing a message, consisting of all of the
// states that are required to do so for 1 transaction.
type SigningSession struct {
	mu *sync.Mutex

	broadcaster  *broadcast.Broadcaster
	TxHash       eth_common.Hash
	MsgToSign    *big.Int
	threshold    int
	partyIDStore *mp_tss.PartyIDStore
	sessionStore *SigningSessionStore // signing session store to update self

	currentPeerID peer.ID
	peerIDs       peer.IDSlice

	resultChan chan SigningSessionResult

	logger *zap.SugaredLogger

	// FSM
	state SigningSessionState

	// Context for the tracing span that lives for the entire duration of the
	// signing session
	context context.Context
	span    trace.Span
}

func NewSigningSession(
	ctx context.Context,
	broadcaster *broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	threshold int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
	partyIDStore *mp_tss.PartyIDStore,
	sessionStore *SigningSessionStore,
) (*SigningSession, <-chan SigningSessionResult, error) {
	tracer := otel.Tracer("NewSigningSession")
	ctx, span := tracer.Start(ctx, "new signing session")

	_, subSpan := tracer.Start(ctx, "picking leader")
	defer subSpan.End()

	resultChan := make(chan SigningSessionResult, 1)

	subSpan.AddEvent("Transition to PickingLeaderState")
	state := NewPickingLeaderState()

	logger := log.Named(txHash.String())

	session := &SigningSession{
		mu:           &sync.Mutex{},
		broadcaster:  broadcaster,
		TxHash:       txHash,
		MsgToSign:    msgToSign,
		threshold:    threshold,
		partyIDStore: partyIDStore,
		sessionStore: sessionStore,

		currentPeerID: currentPeerID,
		peerIDs:       peerIDs,
		resultChan:    resultChan,
		logger:        logger,
		state:         state,
		context:       ctx,
		span:          span,
	}

	// Generate a random session ID part

	// Pick leader
	leaderPeerID, err := GetLeader(session.TxHash, session.peerIDs, state.leaderOffset)
	if err != nil {
		return nil, nil, err
	}

	logger.Debugw("Leader picked", "leaderPeerID", leaderPeerID, "isLeader", currentPeerID == leaderPeerID)

	// Leader is the current peer, transition to LeaderWaitingForCandidatesState
	// No broadcast necessary, other peers know the leader too.
	if leaderPeerID == session.currentPeerID {
		logger.Debugw(
			"IS leader, transition to LeaderWaitingForCandidatesState",
			"peerID",
			leaderPeerID,
		)
		newState, err := NewLeaderWaitingForCandidatesState()
		if err != nil {
			return nil, nil, err
		}

		subSpan.AddEvent("Transition to LeaderWaitingForCandidateState")
		session.state = newState

		return session, resultChan, nil
	}

	// Not leader - send the session ID part to leader and transition to
	// CandidateWaitingForLeaderState
	newState, err := NewCandidateWaitingForLeaderState()
	if err != nil {
		return nil, nil, err
	}

	logger.Debugw("NOT leader, transition to CandidateWaitingForLeaderState")
	subSpan.AddEvent("Transition to CandidateWaitingForLeaderState")
	session.state = newState

	msg := types.NewJoinSigningSessionMessage(
		currentPeerID,
		session.TxHash,
		newState.localPart,
	)

	logger.Debugw("Sending JoinSigningSessionMessage to leader", "leaderPeerID", leaderPeerID)
	err = session.broadcaster.BroadcastMessage(
		ctx,
		&msg,                    // join signing session message
		[]peer.ID{leaderPeerID}, // send to leader
		30,                      // ttl seconds
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to broadcast JoinSigningSessionMessage: %w", err)
	}

	return session, resultChan, nil
}

// -----------------------------------------------------------------------------
// FSM

func (s *SigningSession) Update(event SigningSessionEvent) error {
	// Synchronous updates to prevent race conditions
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debugw(
		"SigningSessionEvent received",
		"peerID", s.currentPeerID,
		"type", event.EventType(),
	)

	switch ev := event.(type) {
	case *AddCandidateEvent:
		return s.UpdateAddCandidateEvent(ev)
	case *StartSignerEvent:
		return s.UpdateStartSignerEvent(ev)
	case *AddSigningPartEvent:
		return s.UpdateAddSigningPartEvent(ev)
	default:
		return fmt.Errorf("unexpected event type: %T", event)
	}
}

func (s *SigningSession) UpdateAddCandidateEvent(
	ev *AddCandidateEvent,
) error {
	state, ok := s.state.(*LeaderWaitingForCandidatesState)
	if !ok {
		s.logger.Errorw(
			"invalid state for AddCandidateEvent",
			"current state", s.state.State(),
			"from", ev.joinMsg.PeerID,
			"currentPeerID", s.currentPeerID,
		)

		return fmt.Errorf(
			"unexpected event %T for state: is currently %T, but event applies to %T",
			ev,
			s.state,
			&LeaderWaitingForCandidatesState{},
		)
	}

	_, span := tracer.Start(s.context, "add candidate to leader")
	defer span.End()

	signingMsg := ev.joinMsg.GetJoinSigningSessionMessage()
	if signingMsg == nil {
		return fmt.Errorf("unexpected join session type: %T", ev.joinMsg)
	}

	// Update state -- requires mutex as there may be concurrent updates
	state.joinMsgsLock.Lock()
	state.joinMsgs = append(state.joinMsgs, ev.joinMsg)
	state.joinMsgsLock.Unlock()

	s.logger.Debugw(
		"Added candidate to leader list",
		"candidate peerID", ev.joinMsg.PeerID.ShortString(),
	)

	// Waits for t, not t+1 because the leader is also a candidate in the next step.
	if len(state.joinMsgs) < s.threshold {
		// Do nothing more, wait for more candidates
		span.AddEvent(
			"Wait for more candidates",
			trace.WithAttributes(
				attribute.Int("current candidates", len(state.joinMsgs)),
				attribute.Int("threshold", s.threshold),
			),
		)

		s.logger.Debugw(
			"Wait for more candidates",
			"current candidates", len(state.joinMsgs),
			"threshold", s.threshold,
		)

		return nil
	}

	span.AddEvent(
		"Adding self (leader) to candidates",
		trace.WithAttributes(
			attribute.String("candidate peerID", s.currentPeerID.ShortString()),
			attribute.String("candidate session part", state.localPart.String()),
		),
	)

	s.logger.Debugw(
		"Adding self (leader) to candidates",
		"candidate peerID", s.currentPeerID.ShortString(),
	)

	// Add self to the list of candidates
	// TODO: why does this cause duplicate peer error when GetSessionID()?
	// Possible race condition when there are two AddCandidateEvent at once?

	selfJoinMsg := types.NewJoinSigningSessionMessage(
		s.currentPeerID,
		s.TxHash,
		state.localPart,
	)

	state.joinMsgsLock.Lock()
	state.joinMsgs = append(state.joinMsgs, selfJoinMsg)
	state.joinMsgsLock.Unlock()

	// Greater than threshold, create aggregate session ID and pick peers to
	// participate in the signing session
	aggSessionID, participantPeerIDs, err := state.joinMsgs.GetSessionID(s.threshold)
	if err != nil {
		return fmt.Errorf("failed to create session ID and pick participants: %w", err)
	}

	span.AddEvent(
		"create aggregate session ID and pick participants",
		trace.WithAttributes(
			attribute.String("aggregate session ID", aggSessionID.String()),
		),
	)

	s.logger.Debugw(
		"aggregate session ID created",
		"participants", participantPeerIDs,
	)

	// Set the aggregate session ID for the session in the "parent" store
	// TODO: Probably refactor this so that we don't need to pass the parent
	// store to the session?
	s.sessionStore.SetSessionID(s.TxHash, aggSessionID)

	// Broadcast StartSignerEvent with picked peers
	msg := types.NewSigningPartyStartMessage(
		s.TxHash,
		aggSessionID,
		participantPeerIDs,
	)
	s.broadcaster.BroadcastMessage(
		context.Background(),
		msg,
		s.peerIDs, // All peers
		30,
	)

	span.AddEvent("Transition leader to SigningState")

	// Transition to signing state
	transport := NewSessionTransport(
		s.broadcaster,
		aggSessionID,
		s.partyIDStore,
		participantPeerIDs,
	)

	// TODO: Not necessarily a participant, either make leader always a participant
	// or check if leader is a participant and transition accordingly.
	s.state = NewSigningState(transport)

	return nil
}

func (s *SigningSession) UpdateStartSignerEvent(
	ev *StartSignerEvent,
) error {
	// Only candidates receive this event, leader should not receive it and
	// will transition by itself.
	_, ok := s.state.(*CandidateWaitingForLeaderState)
	if !ok {
		return fmt.Errorf(
			"unexpected event %T for state: is currently %T, but event applies to %T",
			ev,
			s.state,
			&CandidateWaitingForLeaderState{},
		)
	}

	_, span := tracer.Start(s.context, "start signing",
		trace.WithAttributes(
			attribute.String("event", "StartSignerEvent"),
			attribute.String("peerID", s.currentPeerID.ShortString()),
		))

	s.logger.Debugw(
		"StartSignerEvent received",
		"peerID", s.currentPeerID,
		"participantPeerIDs", ev.participants,
	)

	isParticipant := false
	for _, participant := range ev.participants {
		if participant == s.currentPeerID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		span.AddEvent(
			"Peer is not a participant",
			trace.WithAttributes(
				attribute.String("peerID", s.currentPeerID.ShortString()),
			),
		)
		span.End()

		// TODO: transition to DoneNotParticipantState
		return fmt.Errorf("current peer is not in the list of participants")
	}

	span.AddEvent(
		"Peer is a participant",
		trace.WithAttributes(
			attribute.String("peerID", s.currentPeerID.ShortString()),
		),
	)

	newState := NewSigningState(ev.transport)

	newState.outputChan, newState.errChan = mp_tss.RunSign(
		s.MsgToSign,
		ev.tssParams,
		ev.key,
		ev.transport,
	)

	// Transition CandidatesWaitingForLeaderState -> signing state
	s.state = newState

	// Monitor output and error channels
	go func() {
		// TODO: Do we need a state mutex?
		// This runs concurrently but in the signing state, there will be no
		// other state transitions other than from here.

		select {
		case sig := <-newState.outputChan:
			s.logger.Infow(
				"done signing message",
				"txHash", s.TxHash.String(),
				"signature", sig.String(),
			)

			span.AddEvent(
				"done signing",
				trace.WithAttributes(
					attribute.String("txHash", s.TxHash.String()),
					attribute.String("Signature", sig.String()),
				),
			)

			s.state = NewDoneState(sig)
			s.resultChan <- SigningSessionResult{
				Signature: &sig,
				Err:       nil,
			}
		case err := <-newState.errChan:
			s.logger.Errorw(
				"error signing message",
				"txHash", s.TxHash.String(),
				"error", err.Error(),
			)

			span.RecordError(
				err,
				trace.WithAttributes(
					attribute.String("txHash", s.TxHash.String()),
					attribute.String("culprits", s.currentPeerID.ShortString()),
				),
			)

			s.state = NewErrorState(err)
			s.resultChan <- SigningSessionResult{
				Signature: nil,
				Err:       err,
			}
		case <-s.context.Done():
			s.logger.Debugw("signing session context done, no longer waiting for output")
			// TODO: Transition to err? or nah
		}

		span.End()   // local start span
		s.span.End() // entire session span
	}()

	return nil
}

func (s *SigningSession) UpdateAddSigningPartEvent(
	ev *AddSigningPartEvent,
) error {
	state, ok := s.state.(*SigningState)
	if !ok {
		return fmt.Errorf(
			"unexpected event %T for state: is currently %T, but event applies to %T",
			ev,
			s.state,
			&SigningState{},
		)
	}

	_, span := tracer.Start(s.context, "add signing part",
		trace.WithAttributes(
			attribute.String("event", "AddSigningPartEvent"),
			attribute.String("peerID", s.currentPeerID.ShortString()),
			attribute.String("From", ev.From.String()),
			attribute.Bool("IsBroadcast", ev.IsBroadcast),
		))
	defer span.End()

	// Receive signing part to transport
	state.transport.Receive() <- mp_tss.NewReceivedPartyState(
		ev.Data,
		ev.From,
		ev.IsBroadcast,
	)

	return nil
}

// TryGetSignature returns the signature or error if the session is done
func (s *SigningSession) TryGetSignature() (
	signature tss_common.SignatureData,
	done bool,
	err *tss.Error,
) {
	switch state := s.state.(type) {
	case *DoneState:
		// sig, done, no error
		return state.signature, true, nil
	case *ErrorState:
		// no sig, done, error
		return tss_common.SignatureData{}, true, state.err
	default:
		// no sig, not done, no error
		return tss_common.SignatureData{}, false, nil
	}
}
