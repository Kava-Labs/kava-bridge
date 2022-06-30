package signer

import (
	"context"
	"fmt"
	"math/big"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/zap"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"

	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	mp_tss_types "github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/kava-labs/kava-bridge/relayer/session"
)

var log = logging.Logger("signer")

// Signer is a multi-party signer that handles messages between multiple peers
// for keygen, signing, and resharing.
type Signer struct {
	Node         *p2p.Node
	broadcaster  *broadcast.P2PBroadcaster
	partyIDStore *mp_tss.PartyIDStore
	sessions     *session.SessionStore
	key          keygen.LocalPartySaveData

	partyID   *tss.PartyID
	threshold int

	logger *zap.SugaredLogger
}

// NewSigner returns a new Signer.
func NewSigner(
	node *p2p.Node,
	moniker string,
	partyID *tss.PartyID,
	key keygen.LocalPartySaveData,
	threshold int,
	broadcasterOptions ...broadcast.BroadcasterOption,
) (*Signer, error) {
	signer := &Signer{
		Node: node,
		// Broadcaster set later, broadcaster requires signer struct
		broadcaster:  nil,
		partyIDStore: mp_tss.NewPartyIDStore(),
		sessions:     session.NewSessionStore(),
		key:          key,
		partyID:      partyID,
		threshold:    threshold,
		logger:       log.Named(node.Host.ID().ShortString()),
	}

	broadcaster, err := broadcast.NewBroadcaster(
		context.Background(),
		node.Host,
		append(broadcasterOptions, broadcast.WithHandler(NewBroadcastHandler(signer)))...,
	)
	if err != nil {
		return nil, err
	}

	signer.broadcaster = broadcaster

	var peerMetas []*mp_tss.PeerMetadata
	for i, peerID := range node.PeerList {
		// TODO: Use configured monikers that should reside along with peer IDs
		// to keep track of them easier.
		pMoniker := fmt.Sprintf("%q", rune('A'+i))
		peerMetas = append(peerMetas, mp_tss.NewPeerMetadata(peerID, pMoniker))
	}

	if err := signer.partyIDStore.AddPeers(peerMetas); err != nil {
		return nil, err
	}

	signer.logger.Infof("signer initialized with partyIDStore: %v", signer.partyIDStore)

	return signer, nil
}

// SignMessage signs a message with a corresponding txHash. This creates a
// signing session.
func (s *Signer) SignMessage(
	ctx context.Context,
	txHash eth_common.Hash,
	msgHash *big.Int,
) (*tss_common.SignatureData, error) {
	// Check if there's already a session for the txHash
	_, found := s.sessions.Signing.GetSessionFromTxHash(txHash)
	if found {
		return nil, fmt.Errorf("signing session already exists for txHash %v", txHash)
	}

	// Create new signing session
	_, resultChan, err := s.sessions.Signing.NewSession(
		ctx,
		s.broadcaster,
		txHash,
		msgHash,
		s.threshold,
		s.Node.Host.ID(),
		s.Node.PeerList,
		s.partyID,
		s.partyIDStore,
		s.key,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signing session: %w", err)
	}

	select {
	case res := <-resultChan:
		if res.Err != nil {
			return nil, fmt.Errorf("failed to sign message: %w", res.Err)
		}

		return res.Signature, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *Signer) handleBroadcastMessage(broadcastMsg types.BroadcastMessage) {
	payload, err := broadcastMsg.UnpackPayload()
	if err != nil {
		s.logger.Errorf("failed to unpack received broadcast message: %w", err)

		return
	}

	switch payload := payload.(type) {
	case *mp_tss_types.JoinSessionMessage:
		switch sessionMsg := payload.Session.(type) {
		case *mp_tss_types.JoinSessionMessage_JoinKeygenSessionMessage:
			panic("unimplemented")
		case *mp_tss_types.JoinSessionMessage_JoinResharingSessionMessage:
			panic("unimplemented")
		case *mp_tss_types.JoinSessionMessage_JoinSigningSessionMessage:
			s.handleJoinSigningSessionMessage(broadcastMsg, payload, sessionMsg)
		default:
			panic("unknown session type")
		}

	case *mp_tss_types.SigningPartyStartMessage:
		s.handleSigningPartyStartMessage(broadcastMsg, payload)
	case *mp_tss_types.SigningPartMessage:
		s.handleSigningPartMessage(broadcastMsg, payload)
	default:
		panic("unhandled message type")
	}
}

func (s *Signer) handleJoinSigningSessionMessage(
	broadcastMsg types.BroadcastMessage,
	payload *mp_tss_types.JoinSessionMessage,
	sessionMsg *mp_tss_types.JoinSessionMessage_JoinSigningSessionMessage,
) {
	// Only the leader receives this message
	txHash := sessionMsg.JoinSigningSessionMessage.GetTxHash()

	sess, found := s.sessions.Signing.GetSessionFromTxHash(txHash)

	// TODO: Possible race condition if another peer sends join message
	// before leader creates their own session, this message would be
	// discarded
	if !found {
		s.logger.Errorf("received JoinSigningSessionMessage for unknown txHash %v", txHash)

		return
	}

	// Add potential participant, once we get enough to join, pick
	// actual participants.
	event := session.NewAddCandidateEvent(s.partyID, *payload)
	if err := sess.Update(event); err != nil {
		s.logger.Errorf("failed to add peer to signing session: %v", err)

		return
	}
}

func (s *Signer) handleSigningPartyStartMessage(
	broadcastMsg types.BroadcastMessage,
	payload *mp_tss_types.SigningPartyStartMessage,
) {
	// Start signing sessions
	txHash := payload.GetTxHash()
	sessionID := payload.GetSessionID()
	sess, found := s.sessions.Signing.GetSessionFromTxHash(txHash)
	if !found {
		s.logger.Errorf("received SigningPartyStartMessage for unknown txHash %v", txHash)

		return
	}

	s.sessions.Signing.SetSessionID(txHash, sessionID)

	transport := session.NewSessionTransport(
		s.broadcaster,
		sessionID,
		s.partyIDStore,
		payload.ParticipatingPeerIDs,
	)

	params, err := s.GetParams(payload.ParticipatingPeerIDs)
	if err != nil {
		s.logger.Error(err)
		return
	}

	// Start signer event
	event := session.NewStartSignerEvent(
		params,                       // tss parameters
		transport,                    // broadcast transport
		payload.ParticipatingPeerIDs, // list of participating peer IDs in session
	)
	if err := sess.Update(event); err != nil {
		s.logger.Errorf("failed to start signer: %v", err)

		return
	}
}

func (s *Signer) handleSigningPartMessage(
	broadcastMsg types.BroadcastMessage,
	payload *mp_tss_types.SigningPartMessage,
) {
	sessionID := payload.GetSessionID()

	sess, found := s.sessions.Signing.GetSessionFromID(sessionID)
	if !found {
		s.logger.Errorf("signing session not found for session id %v", sessionID)

		return
	}

	fromPartyID, err := s.partyIDStore.GetPartyID(broadcastMsg.From)
	if err != nil {
		s.logger.Errorf("failed to get party id for peer %v: %v", broadcastMsg.From, err)

		return
	}

	// Update session with signing part
	event := session.NewAddSigningPartEvent(
		fromPartyID,
		payload.Data,
		payload.IsBroadcast,
	)

	if err := sess.Update(event); err != nil {
		s.logger.Errorf("failed to add signing part: %s", err)

		return
	}
}

func (s *Signer) GetParams(participants []peer.ID) (*tss.Parameters, error) {
	partyIDs, err := s.partyIDStore.GetPartyIDs(participants)
	if err != nil {
		return nil, fmt.Errorf("failed to get party IDs: %w", err)
	}

	params := mp_tss.CreateParams(partyIDs, s.partyID, s.threshold)

	return params, nil
}
