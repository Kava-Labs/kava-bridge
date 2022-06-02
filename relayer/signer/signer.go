package signer

import (
	"fmt"
	"math/big"

	logging "github.com/ipfs/go-log/v2"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"

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
	node         *p2p.Node
	partyIDStore *mp_tss.PartyIDStore
	sessions     *session.SessionStore
	key          keygen.LocalPartySaveData
	threshold    int
}

// NewSigner returns a new Signer.
func NewSigner(
	node *p2p.Node,
	moniker string,
	key keygen.LocalPartySaveData,
	threshold int,
) *Signer {
	return &Signer{
		node:         node,
		partyIDStore: mp_tss.NewPartyIDStore(),
		sessions:     session.NewSessionStore(),
		key:          key,
		threshold:    threshold,
	}
}

// Start starts the signer.
func (s *Signer) Start() error {
	// Connect to peers

	// Start listening for messages

	return nil
}

// SignMessage signs a message with a corresponding txHash. This creates a
// signing session.
func (s *Signer) SignMessage(
	txHash eth_common.Hash,
	msgHash *big.Int,
) (tss_common.SignatureData, error) {
	// Check if already signed
	_, found := s.sessions.Signing.GetSessionFromTxHash(txHash)
	if found {
		return tss_common.SignatureData{}, fmt.Errorf("signing session already exists for txHash %v", txHash)
	}

	// Create new signing session
	session := s.sessions.Signing.NewSession(txHash)

	return session.WaitForSignature()
}

func (s *Signer) HandleBroadcastMessage(broadcastMsg types.BroadcastMessage) {
	payload, err := broadcastMsg.UnpackPayload()
	if err != nil {
		log.Errorf("failed to unpack received broadcast message: %w", err)

		return
	}

	switch msg := payload.(type) {
	case *mp_tss_types.JoinSessionMessage:
		switch sessionMsg := msg.Session.(type) {
		case *mp_tss_types.JoinSessionMessage_JoinKeygenSessionMessage:
			panic("unimplemented")
		case *mp_tss_types.JoinSessionMessage_JoinResharingSessionMessage:
			panic("unimplemented")
		case *mp_tss_types.JoinSessionMessage_JoinSigningSessionMessage:
			// Only the leader receives this message
			txHash := sessionMsg.JoinSigningSessionMessage.GetTxHash()

			session, found := s.sessions.Signing.GetSessionFromTxHash(txHash)

			// TODO: Possible race condition if another peer sends join message
			// before leader creates their own session, this message would be
			// discarded
			if !found {
				log.Errorf("received JoinSigningSessionMessage for unknown txHash %v", txHash)

				return
			}

			// Add potential participant, once we get enough to join, pick
			// actual participants.
			sessionIDPart := sessionMsg.JoinSigningSessionMessage.GetPeerSessionIDPart()
			if err := session.AddPotentialParticipant(msg.PeerID, sessionIDPart); err != nil {
				log.Errorf("failed to add peer to signing session: %v", err)

				return
			}

		default:
			panic("unknown session type")
		}

	case *mp_tss_types.SigningPartyStartMessage:
		// Start signing sessions
		sessionID := msg.GetSessionID()
		session, found := s.sessions.Signing.GetSessionFromID(sessionID)
		if !found {
			log.Errorf("received SigningPartyStartMessage for unknown sessionID %v", sessionID)

			return
		}

		transport := NewSessionTransport(sessionID)
		if err := session.StartSigner(s.key, transport); err != nil {
			log.Errorf("failed to start signer: %w", err)

			return
		}
	case *mp_tss_types.SigningPartMessage:
		sessionID := msg.GetSessionID()

		session, found := s.sessions.Signing.GetSessionFromID(sessionID)
		if !found {
			log.Errorf("signing session not found for session id %v", sessionID)

			return
		}

		// Update session
		if err := session.AddSigningPart(
			broadcastMsg.From, // TODO: We are missing the From *tss.PartyID... add to broadcaster
			msg.Data,
			msg.IsBroadcast,
		); err != nil {
			log.Errorf("failed to add signing part: %w", err)

			return
		}
	default:
		panic("unhandled message type")
	}
}

// SessionTransport is a transport for a specific session.
type SessionTransport struct {
	sessionID mp_tss_types.AggregateSigningSessionID

	recvChan chan mp_tss.ReceivedPartyState
}

var _ mp_tss.Transporter = (*SessionTransport)(nil)

func (mt *SessionTransport) Send(data []byte, routing *tss.MessageRouting, isResharing bool) error {
	return nil
}

func (mt *SessionTransport) Receive() <-chan mp_tss.ReceivedPartyState {
	return mt.recvChan
}

func NewSessionTransport(sessionID mp_tss_types.AggregateSigningSessionID) mp_tss.Transporter {
	return &SessionTransport{
		sessionID: sessionID,
		recvChan:  make(chan mp_tss.ReceivedPartyState, 1),
	}
}
