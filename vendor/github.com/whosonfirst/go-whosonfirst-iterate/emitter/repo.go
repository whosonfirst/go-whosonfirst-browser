package emitter

import (
	"context"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "repo", NewRepoEmitter)
}

type RepoEmitter struct {
	Emitter
	emitter Emitter
}

func NewRepoEmitter(ctx context.Context, uri string) (Emitter, error) {

	directory_idx, err := NewDirectoryEmitter(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &RepoEmitter{
		emitter: directory_idx,
	}

	return idx, nil
}

func (idx *RepoEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	data_path := filepath.Join(abs_path, "data")
	return idx.emitter.WalkURI(ctx, index_cb, data_path)
}
