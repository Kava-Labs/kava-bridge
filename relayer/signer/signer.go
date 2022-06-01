package signer

import (
	"fmt"
	"math/big"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/kava-labs/kava-bridge/relayer/session/signing"
)

// Signer is a multi-party signer that handles messages between multiple peers
// for keygen, signing, and resharing.
type Signer struct {
	node                *p2p.Node
	partyIDStore        *mp_tss.PartyIDStore
	signingSessionStore *signing.SigningSessionStore
	key                 keygen.LocalPartySaveData
	threshold           int
}

// NewSigner returns a new Signer.
func NewSigner(
	node *p2p.Node,
	moniker string,
	key keygen.LocalPartySaveData,
	threshold int,
) *Signer {
	return &Signer{
		node:                node,
		partyIDStore:        mp_tss.NewPartyIDStore(),
		signingSessionStore: signing.NewSigningSessionStore(),
		key:                 key,
		threshold:           threshold,
	}
}

// Start starts the signer.
func (s *Signer) Start() error {
	// Connect to peers

	// Start listening for messages

	return nil
}

// SignMessage signs a message with a corresponding txHash.
func (s *Signer) SignMessage(
	txHash eth_common.Hash,
	msgHash *big.Int,
) (tss_common.SignatureData, error) {
	// Check if already signed
	_, found := s.signingSessionStore.GetSessionFromTxHash(txHash)
	if found {
		return tss_common.SignatureData{}, fmt.Errorf("signing session already exists for txHash %v", txHash)
	}

	// Start signing session
	s.signingSessionStore.NewSession(txHash)

	return tss_common.SignatureData{}, nil
}
