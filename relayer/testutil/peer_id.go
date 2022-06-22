package testutil

import "github.com/libp2p/go-libp2p-core/peer"

var TestPeerIDs = []peer.ID{
	MustDecodePeerID("16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9"),
	MustDecodePeerID("16Uiu2HAkwC5w1fC4xLL3hWjD6PGuk2qzGgsWdXfNeqMi8xDn2AT7"),
	MustDecodePeerID("16Uiu2HAmTdEddBdw1JVs5tHhqQGaFPkqq64TwppmL2G8fYbZeZei"),
	MustDecodePeerID("16Uiu2HAm48kushm1oCczim8L6adoCXV9A3npssbAFMM6Cgw3pYcS"),
	MustDecodePeerID("16Uiu2HAmR3sVyPwLtkemjCoU1XbZcjxcfRHzpopydvEEFbHKGiq4"),
}

func MustDecodePeerID(s string) peer.ID {
	p, err := peer.Decode(s)
	if err != nil {
		panic(err)
	}
	return p
}
