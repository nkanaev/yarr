package eventstore

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

// Store is a persistence layer for nostr events handled by a relay.
type Store interface {
	// Init is called at the very beginning by [Server.Start], after [Relay.Init],
	// allowing a storage to initialize its internal resources.
	Init() error

	// Close must be called after you're done using the store, to free up resources and so on.
	Close()

	// QueryEvents is invoked upon a client's REQ as described in NIP-01.
	// it should return a channel with the events as they're recovered from a database.
	// the channel should be closed after the events are all delivered.
	QueryEvents(context.Context, nostr.Filter) (chan *nostr.Event, error)
	// DeleteEvent is used to handle deletion events, as per NIP-09.
	DeleteEvent(context.Context, *nostr.Event) error
	// SaveEvent is called once Relay.AcceptEvent reports true.
	SaveEvent(context.Context, *nostr.Event) error
}

type Counter interface {
	CountEvents(context.Context, nostr.Filter) (int64, error)
}
