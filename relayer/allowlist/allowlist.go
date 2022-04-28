package allowlist

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/control"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
)

// PeerIDAllowList is a libp2p config option that configures an allowlist of
// Peer IDs. Any peer ID not in the allowlist will be rejected for both
// incoming and outgoing connections.
func PeerIDAllowList(peerIDs []peer.ID) config.Option {
	return libp2p.ConnectionGater(NewAllowList(peerIDs))
}

// AllowListConnectionGater is a connmgr.ConnectionGater implementation that
// only allows connections to and from the specified peer IDs.
type AllowListConnectionGater struct {
	// PeerIDs is a map of allowed peer.ID, converted from a slice for
	// easier and constant time lookups.
	PeerIDs map[peer.ID]struct{}
}

var _ connmgr.ConnectionGater = (*AllowListConnectionGater)(nil)

// NewAllowList returns a new AllowList.
func NewAllowList(peerIDs []peer.ID) *AllowListConnectionGater {
	peerIDsMap := make(map[peer.ID]struct{})
	for _, addr := range peerIDs {
		peerIDsMap[addr] = struct{}{}
	}

	return &AllowListConnectionGater{
		PeerIDs: peerIDsMap,
	}
}

// IsPeerIDAllowed returns true if the peer ID is contained in the allowed list.
func (allowlist *AllowListConnectionGater) IsPeerIDAllowed(peerID peer.ID) bool {
	_, ok := allowlist.PeerIDs[peerID]
	return ok
}

// InterceptPeerDial tests whether we're permitted to Dial the specified peer.
//
// This is called by the network.Network implementation when dialling a peer.
func (allowlist *AllowListConnectionGater) InterceptPeerDial(peerID peer.ID) bool {
	return allowlist.IsPeerIDAllowed(peerID)
}

// InterceptAddrDial tests whether we're permitted to dial the specified
// multiaddr for the given peer.
//
// This is called by the network.Network implementation after it has
// resolved the peer's addrs, and prior to dialling each.
func (allowlist *AllowListConnectionGater) InterceptAddrDial(peerID peer.ID, addr ma.Multiaddr) bool {
	return allowlist.IsPeerIDAllowed(peerID)
}

// InterceptAccept tests whether an incipient inbound connection is allowed.
//
// This is called by the upgrader, or by the transport directly (e.g. QUIC,
// Bluetooth), straight after it has accepted a connection from its socket.
func (allowlist *AllowListConnectionGater) InterceptAccept(network.ConnMultiaddrs) bool {
	// Peer connection IPs are not handled by relayer and should be configured
	// in the machine firewall instead.
	return true
}

// InterceptSecured tests whether a given connection, now authenticated,
// is allowed.
//
// This is called by the upgrader, after it has performed the security
// handshake, and before it negotiates the muxer, or by the directly by the
// transport, at the exact same checkpoint.
func (allowlist *AllowListConnectionGater) InterceptSecured(
	direction network.Direction,
	peerID peer.ID,
	multiAddr network.ConnMultiaddrs,
) bool {
	return allowlist.IsPeerIDAllowed(peerID)
}

// InterceptUpgraded tests whether a fully capable connection is allowed.
//
// At this point, the connection a multiplexer has been selected.
// When rejecting a connection, the gater can return a DisconnectReason.
// Refer to the godoc on the ConnectionGater type for more information.
//
// NOTE: the go-libp2p implementation currently IGNORES the disconnect reason.
func (allowlist *AllowListConnectionGater) InterceptUpgraded(network.Conn) (bool, control.DisconnectReason) {
	// A zero value stands for "no reason" / NA.
	return true, 0
}
