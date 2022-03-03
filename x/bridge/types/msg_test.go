package types_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgBridgeERC20FromEthereum(t *testing.T) {
	type args struct {
		relayer              string
		ethereumERC20Address string
		amount               sdk.Int
		receiver             string
		sequence             sdk.Int
	}
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{
			"valid",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - sequence 0 when overflow",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(0),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - empty relayer",
			args{
				"",
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: false,
				contains:   "empty address string is not allowed: invalid address",
			},
		},
		{
			"invalid - erc20 hex address length",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756C",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: false,
				contains:   "ethereum ERC20 address is not a valid hex address: invalid address",
			},
		},
		{
			"invalid - receiver hex address length",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF1",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: false,
				contains:   "receiver address is not a valid hex address: invalid address",
			},
		},
		{
			"invalid - negative amount",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(-1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: false,
				contains:   "amount must be positive non-zero",
			},
		},
		{
			"invalid - zero amount",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(0),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			},
			errArgs{
				expectPass: false,
				contains:   "amount must be positive non-zero",
			},
		},
		{
			"invalid - negative sequence",
			args{
				sdk.AccAddress("hi").String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(-123),
			},
			errArgs{
				expectPass: false,
				contains:   "sequence is negative",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgBridgeERC20FromEthereum(
				tc.args.relayer,
				tc.args.ethereumERC20Address,
				tc.args.amount,
				tc.args.receiver,
				tc.args.sequence,
			)
			err := msg.ValidateBasic()

			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestMsgBridgeERC20FromEthereumSigners(t *testing.T) {
	relayer := sdk.AccAddress("hi")

	msg := types.NewMsgBridgeERC20FromEthereum(
		relayer.String(),
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		sdk.NewInt(1234),
		"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
		sdk.NewInt(1),
	)

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, relayer, signers[0])
}

func TestMsgBridgeERC20FromEthereumSigners_Invalid(t *testing.T) {
	msg := types.NewMsgBridgeERC20FromEthereum(
		"not a valid address",
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		sdk.NewInt(1234),
		"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
		sdk.NewInt(1),
	)

	require.Panics(t, func() {
		msg.GetSigners()
	})
}
