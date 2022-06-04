package session

import (
	"context"
	"fmt"
	"math/big"

	logging "github.com/ipfs/go-log/v2"
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

// tx_hash -> session ID
type TssSessions map[eth_common.Hash]*SigningSession

// session -> tx_hash to get a session from session ID instead of tx_hash
type SessionIDToTxHash map[string]eth_common.Hash

// SigningSessionStore keeps track of signing sessions.
type SigningSessionStore struct {
	sessions          TssSessions
	sessionIDToTxHash SessionIDToTxHash
}

func NewSigningSessionStore() *SigningSessionStore {
	return &SigningSessionStore{
		sessions:          make(TssSessions),
		sessionIDToTxHash: make(SessionIDToTxHash),
	}
}

func (s *SigningSessionStore) NewSession(
	broadcaster *broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	threshold int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
) (*SigningSession, <-chan SigningSessionResult, error) {
	session, resultChan, err := NewSigningSession(
		broadcaster,
		txHash,
		msgToSign,
		threshold,
		currentPeerID,
		peerIDs,
	)
	if err != nil {
		return nil, nil, err
	}

	s.sessions[txHash] = session

	return session, resultChan, nil
}

func (s *SigningSessionStore) GetSessionFromTxHash(txHash eth_common.Hash) (*SigningSession, bool) {
	session, ok := s.sessions[txHash]
	return session, ok
}

func (s *SigningSessionStore) GetSessionFromID(
	sessID types.AggregateSigningSessionID,
) (*SigningSession, bool) {
	txHash, ok := s.sessionIDToTxHash[sessID.String()]
	if !ok {
		return nil, false
	}

	return s.GetSessionFromTxHash(txHash)
}

// SigningSessionResult is the result of a signing session.
type SigningSessionResult struct {
	Signature *tss_common.SignatureData
	Err       *tss.Error
}

type SigningSession struct {
	broadcaster *broadcast.Broadcaster
	TxHash      eth_common.Hash
	MsgToSign   *big.Int
	threshold   int

	currentPeerID peer.ID
	peerIDs       peer.IDSlice

	resultChan chan SigningSessionResult

	logger *zap.SugaredLogger

	// FSM
	state SigningSessionState
}

func NewSigningSession(
	broadcaster *broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	threshold int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
) (*SigningSession, <-chan SigningSessionResult, error) {
	resultChan := make(chan SigningSessionResult, 1)
	state := NewPickingLeaderState()

	logger := log.Named(txHash.String())

	session := &SigningSession{
		broadcaster:   broadcaster,
		TxHash:        txHash,
		MsgToSign:     msgToSign,
		threshold:     threshold,
		currentPeerID: currentPeerID,
		peerIDs:       peerIDs,
		resultChan:    resultChan,
		logger:        logger,
		state:         state,
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
	session.state = newState

	msg := types.NewJoinSigningSessionMessage(
		currentPeerID,
		session.TxHash,
		newState.localPart,
	)

	logger.Debugw("Sending JoinSigningSessionMessage to leader", "leaderPeerID", leaderPeerID)
	err = session.broadcaster.BroadcastMessage(
		context.Background(),
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
	s.logger.Debugw(
		"SigningSessionEvent received",
		"peerID", s.currentPeerID,
		"type", event.EventType(),
		"event", event,
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
			"unexpected state type: got %T, expected %T",
			s.state,
			&LeaderWaitingForCandidatesState{},
		)
	}

	signingMsg := ev.joinMsg.GetJoinSigningSessionMessage()
	if signingMsg == nil {
		return fmt.Errorf("unexpected join session type: %T", ev.joinMsg)
	}

	// Update state
	state.joinMsgs = append(state.joinMsgs, ev.joinMsg)

	if len(state.joinMsgs) <= s.threshold {
		// Do nothing more, wait for more candidates
		return nil
	}

	// Greater than threshold, create aggregate session ID and pick peers to
	// participate in the signing session
	aggSessionID, participantPeerIDs, err := state.joinMsgs.GetSessionID(s.threshold)
	if err != nil {
		return fmt.Errorf("failed to create session ID and pick participants: %w", err)
	}

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

	// Transition to signing state
	transport := NewSessionTransport(s.broadcaster, aggSessionID)
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
		return fmt.Errorf("unexpected state type: %T", s.state)
	}

	isParticipant := false
	for _, participant := range ev.participants {
		if participant == s.currentPeerID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		// TODO: transition to DoneNotParticipantState
		return fmt.Errorf("current peer is not in the list of participants")
	}

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
			s.state = NewDoneState(sig)
			s.resultChan <- SigningSessionResult{
				Signature: &sig,
			}
		case err := <-newState.errChan:
			s.state = NewErrorState(err)
			s.resultChan <- SigningSessionResult{
				Err: err,
			}
		}
	}()

	return nil
}

func (s *SigningSession) UpdateAddSigningPartEvent(
	ev *AddSigningPartEvent,
) error {
	state, ok := s.state.(*SigningState)
	if !ok {
		return fmt.Errorf("unexpected state type: got %T, expected %T", s.state, &SigningState{})
	}

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
