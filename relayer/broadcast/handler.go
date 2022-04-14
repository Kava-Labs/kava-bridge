package broadcast

import "context"

type BroadcastHandler interface {
	HandleRawMessage(ctx context.Context, msg *RPC) error
}
