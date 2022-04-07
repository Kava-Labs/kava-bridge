package p2p_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalPrivateKey_Bytes(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name     string
		keyBytes []byte
		errArgs  errArgs
	}{
		{
			"valid - secp256k1",
			MustDecodeHexString("080212205847a025649cdafe118e613cbc3440dd08d4bae1a0fdee018b22cb38f9a2e862"),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - ecdsa",
			MustDecodeHexString("0803127930770201010420eeda242d02a2b8738f9546bffd36383420c6dd79173c42c2b8d3d09730e16fc7a00a06082a8648ce3d030107a14403420004170205387129fc78a60755a843f366ebff33cfcfd828a186a2f6b0df6f6fd1222b7a0eb9fa14864a392397b1d5751b079d0499f544461c6d7952dc758c3abe90"),
			errArgs{
				expectPass: false,
				contains:   "invalid key type ECDSA",
			},
		},
		{
			"invalid - empty",
			[]byte{},
			errArgs{
				expectPass: false,
				contains:   "could not decode private key",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := p2p.UnmarshalPrivateKey(tc.keyBytes)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestUnmarshalPrivateKey_EqualMarshal(t *testing.T) {
	secpPrivKey, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	require.NoError(t, err)

	ecdsaPrivKey, _, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	require.NoError(t, err)

	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		key     crypto.PrivKey
		errArgs errArgs
	}{
		{
			"valid - secp256k1",
			secpPrivKey,
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - ecdsa",
			ecdsaPrivKey,
			errArgs{
				expectPass: false,
				contains:   "invalid key type ECDSA",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			keyBytes, err := crypto.MarshalPrivateKey(tc.key)
			require.NoError(t, err)

			t.Log(hex.EncodeToString(keyBytes))

			unmarshaledPrivKey, err := p2p.UnmarshalPrivateKey(keyBytes)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
				require.Equal(t, tc.key, unmarshaledPrivKey)
				require.True(t, tc.key.Equals(unmarshaledPrivKey))
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func MustDecodeHexString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}
