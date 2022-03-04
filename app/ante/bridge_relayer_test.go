package ante_test

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.AnteHandler = (&MockAnteHandler{}).AnteHandle
)

type MockAnteHandler struct {
	WasCalled bool
}

func (mah *MockAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	mah.WasCalled = true
	return ctx, nil
}
