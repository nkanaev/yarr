package eventstore

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

type RelayWrapper struct {
	Store
}

var _ nostr.RelayStore = (*RelayWrapper)(nil)

func (w RelayWrapper) Publish(ctx context.Context, evt nostr.Event) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if 20000 <= evt.Kind && evt.Kind < 30000 {
		// do not store ephemeral events
		return nil
	} else if evt.Kind == 0 || evt.Kind == 3 || (10000 <= evt.Kind && evt.Kind < 20000) {
		// replaceable event, delete before storing
		ch, err := w.Store.QueryEvents(ctx, nostr.Filter{Authors: []string{evt.PubKey}, Kinds: []int{evt.Kind}})
		if err != nil {
			return fmt.Errorf("failed to query before replacing: %w", err)
		}
		isNewer := true
		for previous := range ch {
			if previous == nil {
				continue
			}
			if isOlder(previous, &evt) {
				if err := w.Store.DeleteEvent(ctx, previous); err != nil {
					return fmt.Errorf("failed to delete event for replacing: %w", err)
				}
			} else {
				// already, newer event is stored.
				isNewer = false
				break
			}
		}
		if !isNewer {
			return nil
		}
	} else if 30000 <= evt.Kind && evt.Kind < 40000 {
		// parameterized replaceable event, delete before storing
		d := evt.Tags.GetFirst([]string{"d", ""})
		if d == nil {
			return fmt.Errorf("failed to add event missing d tag for parameterized replacing")
		}
		ch, err := w.Store.QueryEvents(ctx, nostr.Filter{Authors: []string{evt.PubKey}, Kinds: []int{evt.Kind}, Tags: nostr.TagMap{"d": []string{d.Value()}}})
		if err != nil {
			return fmt.Errorf("failed to query before parameterized replacing: %w", err)
		}
		isNewer := true
		for previous := range ch {
			if previous == nil {
				continue
			}

			if !isOlder(previous, &evt) {
				if err := w.Store.DeleteEvent(ctx, previous); err != nil {
					return fmt.Errorf("failed to delete event for parameterized replacing: %w", err)
				}
			} else {
				// already, newer event is stored.
				isNewer = false
				break
			}
		}
		if !isNewer {
			return nil
		}
	}

	if err := w.SaveEvent(ctx, &evt); err != nil && err != ErrDupEvent {
		return fmt.Errorf("failed to save: %w", err)
	}

	return nil
}

func (w RelayWrapper) QuerySync(ctx context.Context, filter nostr.Filter) ([]*nostr.Event, error) {
	ch, err := w.Store.QueryEvents(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	n := filter.Limit
	if n == 0 {
		n = 500
	}

	results := make([]*nostr.Event, 0, n)
	for evt := range ch {
		results = append(results, evt)
	}

	return results, nil
}
