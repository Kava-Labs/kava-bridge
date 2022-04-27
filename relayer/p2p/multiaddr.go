package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ParseMultiAddrSlice parses a slice of multiaddrs strings into a slice of addr info
// Requires a valid multiaddr string *and* a valid peer ID.
func ParseMultiAddrSlice(addrStrings []string) ([]*peer.AddrInfo, error) {
	var addrs []*peer.AddrInfo

	for _, s := range addrStrings {
		peerAddr, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			return nil, fmt.Errorf("could not parse peer multiaddr %s: %w", s, err)
		}

		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			return nil, fmt.Errorf("could not parse addrinfo from Multiaddr: %w", err)
		}

		if len(peerAddrInfo.Addrs) == 0 {
			return nil, fmt.Errorf("no transport multiaddr found in peer info %s", s)
		}

		addrs = append(addrs, peerAddrInfo)
	}

	return addrs, nil
}

// PeerIDsFromAddrInfos returns a slice of peer IDs from the given slice of
// peer.AddrInfo.
func PeerIDsFromAddrInfos(addrInfos []*peer.AddrInfo) []peer.ID {
	var peerIDs []peer.ID

	for _, addrInfo := range addrInfos {
		peerIDs = append(peerIDs, addrInfo.ID)
	}

	return peerIDs
}
