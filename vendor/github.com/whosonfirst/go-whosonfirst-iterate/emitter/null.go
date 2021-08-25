package emitter

import (
	"context"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "null", NewNullEmitter)
}

type NullEmitter struct {
	Emitter
}

func NewNullEmitter(ctx context.Context, uri string) (Emitter, error) {

	idx := &NullEmitter{}
	return idx, nil
}

func (idx *NullEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {
	return nil
}
