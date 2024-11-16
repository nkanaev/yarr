package nostr

import (
	"context"
	"errors"
	"slices"
)

type RelayStore interface {
	Publish(context.Context, Event) error
	QueryEvents(context.Context, Filter) (chan *Event, error)
	QuerySync(context.Context, Filter) ([]*Event, error)
}

var (
	_ RelayStore = (*Relay)(nil)
	_ RelayStore = (*MultiStore)(nil)
)

type MultiStore []RelayStore

func (multi MultiStore) Publish(ctx context.Context, event Event) error {
	errs := make([]error, len(multi))
	for i, s := range multi {
		errs[i] = s.Publish(ctx, event)
	}
	return errors.Join(errs...)
}

func (multi MultiStore) QueryEvents(ctx context.Context, filter Filter) (chan *Event, error) {
	multich := make(chan *Event)

	errs := make([]error, len(multi))
	var good bool
	for i, s := range multi {
		ch, err := s.QueryEvents(ctx, filter)
		errs[i] = err
		if err == nil {
			good = true
			go func(ch chan *Event) {
				for evt := range ch {
					multich <- evt
				}
			}(ch)
		}
	}

	if good {
		return multich, nil
	} else {
		return nil, errors.Join(errs...)
	}
}

func (multi MultiStore) QuerySync(ctx context.Context, filter Filter) ([]*Event, error) {
	errs := make([]error, len(multi))
	events := make([]*Event, 0, max(filter.Limit, 250))
	for i, s := range multi {
		res, err := s.QuerySync(ctx, filter)
		errs[i] = err
		events = append(events, res...)
	}
	slices.SortFunc(events, func(a, b *Event) int {
		if b.CreatedAt > a.CreatedAt {
			return 1
		} else if b.CreatedAt < a.CreatedAt {
			return -1
		}
		return 0
	})
	return events, errors.Join(errs...)
}
