package session

import (
	"context"
	"fmt"
	"math/big"

	tss_common "github.com/binance-chain/tss-lib/common"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

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
) (*SigningSession, error) {
	session, err := NewSigningSession(
		broadcaster,
		txHash,
		msgToSign,
		threshold,
		currentPeerID,
		peerIDs,
	)
	if err != nil {
		return nil, err
	}

	s.sessions[txHash] = session

	return session, nil
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

type SigningSession struct {
	broadcaster *broadcast.Broadcaster
	TxHash      eth_common.Hash
	MsgToSign   *big.Int
	threshold   int

	currentPeerID peer.ID
	peerIDs       peer.IDSlice

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
) (*SigningSession, error) {
	state := NewPickingLeaderState()

	session := &SigningSession{
		broadcaster:   broadcaster,
		TxHash:        txHash,
		MsgToSign:     msgToSign,
		threshold:     threshold,
		currentPeerID: currentPeerID,
		peerIDs:       peerIDs,
		state:         state,
	}

	// Generate a random session ID part

	// Pick leader
	leaderPeerID, err := GetLeader(session.TxHash, session.peerIDs, state.leaderOffset)
	if err != nil {
		return nil, err
	}

	// Leader is the current peer, transition to LeaderWaitingForCandidatesState
	// No broadcast necessary, other peers know the leader too.
	if leaderPeerID == session.currentPeerID {
		session.state, err = NewLeaderWaitingForCandidatesState()
		if err != nil {
			return nil, err
		}

		return session, nil
	}

	// Not leader - send the session ID part to leader and transition to
	// CandidateWaitingForLeaderState
	newState, err := NewCandidateWaitingForLeaderState()
	if err != nil {
		return nil, err
	}

	session.state = newState

	msg := types.NewJoinSigningSessionMessage(
		currentPeerID,
		session.TxHash,
		newState.localPart,
	)

	err = session.broadcaster.BroadcastMessage(
		context.Background(),
		&msg,                    // join signing session message
		[]peer.ID{leaderPeerID}, // send to leader
		30,                      // ttl seconds
	)

	if err != nil {
		return nil, fmt.Errorf("failed to broadcast JoinSigningSessionMessage: %w", err)
	}

	return session, nil
}

// -----------------------------------------------------------------------------
// FSM

func (s *SigningSession) Update(event SigningSessionEvent) error {
	switch ev := event.(type) {
	case *AddCandidateEvent:
		return s.UpdateAddCandidateEvent(ev)
	case *StartSignerEvent:
		return s.UpdateStartSignerEvent(ev)
	case *AddSigningPartEvent:
	default:
		return fmt.Errorf("unexpected event type: %T", event)
	}

	return nil
}

func (s *SigningSession) UpdateAddCandidateEvent(
	ev *AddCandidateEvent,
) error {
	state, ok := s.state.(*LeaderWaitingForCandidatesState)
	if !ok {
		return fmt.Errorf("unexpected state type: %T", s.state)
	}

	signingMsg := ev.joinMsg.GetJoinSigningSessionMessage()
	if signingMsg == nil {
		return fmt.Errorf("unexpected join session type: %T", ev.joinMsg)
	}

	// Update state
	state.joinMsgs = append(state.joinMsgs, ev.joinMsg)

	if len(state.joinMsgs) <= state.threshold {
		// Do nothing more, wait for more candidates
		return nil
	}

	// Greater than threshold, create aggregate session ID and pick peers to
	// participate in the signing session
	aggSessionID, participantPeerIDs, err := state.joinMsgs.GetSessionID(state.threshold)
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
	s.state = NewSigningState(state.transport)

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
		// TODO: Do we need a mutex here?
		// This runs concurrently but in the signing state, there will be no
		// other state changes other than from here.

		select {
		case sig := <-newState.outputChan:
			s.state = NewDoneState(sig)
		case err := <-newState.errChan:
			s.state = NewErrorState(err)
		}
	}()

	return nil
}

func (s *SigningSession) AddSigningPartEvent(
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

func (s *SigningSession) TryGetSignature() (tss_common.SignatureData, error) {
	state, ok := s.state.(*DoneState)
	if !ok {
		return tss_common.SignatureData{}, fmt.Errorf("unexpected state type: %T", s.state)
	}

	return state.signature, nil
}
