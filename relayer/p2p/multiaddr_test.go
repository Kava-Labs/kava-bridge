package p2p_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestParseMultiAddrSlice(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name        string
		addrStrings []string
		errArgs     errArgs
	}{
		{
			"valid - empty",
			[]string{},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - transport and peer ID",
			[]string{
				"/ip4/192.168.1.24/tcp/8765/p2p/16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9",
				"/ip4/127.0.0.1/tcp/8765/p2p/16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9",
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - transport but no peer ID",
			[]string{
				"/ip4/192.168.1.24/tcp/8765",
			},
			errArgs{
				expectPass: false,
				contains:   peer.ErrInvalidAddr.Error(),
			},
		},
		{
			"invalid - peer ID but no transport",
			[]string{
				"/p2p/16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9",
			},
			errArgs{
				expectPass: false,
				contains:   "no transport multiaddr found in peer info",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addrInfos, err := p2p.ParseMultiAddrSlice(tc.addrStrings)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
				require.Equal(t, len(tc.addrStrings), len(addrInfos))
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}

}
