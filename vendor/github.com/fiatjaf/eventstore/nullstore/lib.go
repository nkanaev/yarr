package nullstore

import (
	"context"

	"github.com/fiatjaf/eventstore"
	"github.com/nbd-wtf/go-nostr"
)

var _ eventstore.Store = NullStore{}

type NullStore struct{}

func (b NullStore) Init() error {
	return nil
}

func (b NullStore) Close() {}

func (b NullStore) DeleteEvent(ctx context.Context, evt *nostr.Event) error {
	return nil
}

func (b NullStore) QueryEvents(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
	ch := make(chan *nostr.Event)
	close(ch)
	return ch, nil
}

func (b NullStore) SaveEvent(ctx context.Context, evt *nostr.Event) error {
	return nil
}
