package worker

import (
	"context"
	"errors"
	"time"

	"github.com/nkanaev/yarr/src/client"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

type feedIDContextKey struct{}

func contextSetFeedID(ctx context.Context, feedID int64) context.Context {
	return context.WithValue(ctx, feedIDContextKey{}, feedID)
}

func contextGetFeedID(ctx context.Context) (int64, error) {
	feedID, ok := ctx.Value(feedIDContextKey{}).(int64)
	if !ok {
		return 0, errors.New("feed id missing from context")
	}
	return feedID, nil
}

type feedStateStore struct {
	db storage.Storage
}

func NewFeedStateStore(db storage.Storage) client.StateStore {
	return &feedStateStore{db: db}
}

func (s *feedStateStore) Get(ctx context.Context) (client.State, error) {
	feedID, err := contextGetFeedID(ctx)
	if err != nil {
		return client.State{}, err
	}
	state, err := s.db.GetFeedState(feedID)
	if err != nil {
		return client.State{}, err
	}
	if state == nil {
		return client.State{}, nil
	}
	return client.State{
		ETag:         state.HTTPEtag,
		LastModified: state.HTTPLastModified,
	}, nil
}

func (s *feedStateStore) Put(ctx context.Context, state client.State) error {
	feedID, err := contextGetFeedID(ctx)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = s.db.UpdateFeedState(feedID, model.UpdateFeedStateParams{
		HTTPLastModified: &state.LastModified,
		HTTPEtag:         &state.ETag,
		LastRefreshed:    &now,
	})
	return err
}
