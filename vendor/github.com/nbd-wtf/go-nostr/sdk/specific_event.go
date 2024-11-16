package sdk

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// FetchSpecificEvent tries to get a specific event from a NIP-19 code using whatever means necessary.
func (sys *System) FetchSpecificEvent(
	ctx context.Context,
	code string,
	withRelays bool,
) (event *nostr.Event, successRelays []string, err error) {
	// this is for deciding what relays will go on nevent and nprofile later
	priorityRelays := make([]string, 0, 8)

	prefix, data, err := nip19.Decode(code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode %w", err)
	}

	author := ""

	var filter nostr.Filter
	relays := make([]string, 0, 10)
	fallback := make([]string, 0, 10)
	successRelays = make([]string, 0, 10)

	switch v := data.(type) {
	case nostr.EventPointer:
		author = v.Author
		filter.IDs = []string{v.ID}
		relays = append(relays, v.Relays...)
		relays = appendUnique(relays, sys.FallbackRelays.Next())
		fallback = append(fallback, sys.JustIDRelays.URLs...)
		fallback = appendUnique(fallback, sys.FallbackRelays.Next())
		for _, r := range v.Relays {
			priorityRelays = append(priorityRelays, r)
		}
	case nostr.EntityPointer:
		author = v.PublicKey
		filter.Authors = []string{v.PublicKey}
		filter.Tags = nostr.TagMap{"d": []string{v.Identifier}}
		filter.Kinds = []int{v.Kind}
		relays = append(relays, v.Relays...)
		relays = appendUnique(relays, sys.FallbackRelays.Next())
		fallback = append(fallback, sys.FallbackRelays.Next(), sys.FallbackRelays.Next())
	case string:
		if prefix == "note" {
			filter.IDs = []string{v}
			relays = append(relays, sys.JustIDRelays.Next(), sys.JustIDRelays.Next())
			fallback = appendUnique(fallback,
				sys.FallbackRelays.Next(), sys.JustIDRelays.Next(), sys.FallbackRelays.Next())
		}
	}

	// try to fetch in our internal eventstore first
	if res, _ := sys.StoreRelay.QuerySync(ctx, filter); len(res) != 0 {
		evt := res[0]
		return evt, nil, nil
	}

	if author != "" {
		// fetch relays for author
		authorRelays := sys.FetchOutboxRelays(ctx, author, 3)
		relays = appendUnique(relays, authorRelays...)
		priorityRelays = appendUnique(priorityRelays, authorRelays...)
	}

	var result *nostr.Event
	fetchProfileOnce := sync.Once{}

attempts:
	for _, attempt := range []struct {
		label          string
		relays         []string
		slowWithRelays bool
	}{
		{
			label:  "fetch-" + prefix,
			relays: relays,
			// set this to true if the caller wants relays, so we won't return immediately
			//   but will instead wait a little while to see if more relays respond
			slowWithRelays: withRelays,
		},
		{
			label:          "fetchf-" + prefix,
			relays:         fallback,
			slowWithRelays: false,
		},
	} {
		// actually fetch the event here
		countdown := 6.0
		subManyCtx := ctx
		subMany := sys.Pool.SubManyEose
		if attempt.slowWithRelays {
			subMany = sys.Pool.SubManyEoseNonUnique
		}

		if attempt.slowWithRelays {
			// keep track of where we have actually found the event so we can show that
			var cancel context.CancelFunc
			subManyCtx, cancel = context.WithTimeout(ctx, time.Second*6)
			defer cancel()

			go func() {
				for {
					time.Sleep(100 * time.Millisecond)
					if countdown <= 0 {
						cancel()
						break
					}
					countdown -= 0.1
				}
			}()
		}

		for ie := range subMany(
			subManyCtx,
			attempt.relays,
			nostr.Filters{filter},
			nostr.WithLabel(attempt.label),
		) {
			fetchProfileOnce.Do(func() {
				go sys.FetchProfileMetadata(ctx, ie.PubKey)
			})

			successRelays = append(successRelays, ie.Relay.URL)
			if result == nil || ie.CreatedAt > result.CreatedAt {
				result = ie.Event
			}

			if !attempt.slowWithRelays {
				break attempts
			}

			countdown = min(countdown-0.5, 1)
		}
	}

	if result == nil {
		return nil, nil, fmt.Errorf("couldn't find this %s", prefix)
	}

	// save stuff in cache and in internal store
	sys.StoreRelay.Publish(ctx, *result)

	// put priority relays first so they get used in nevent and nprofile
	slices.SortFunc(successRelays, func(a, b string) int {
		vpa := slices.Contains(priorityRelays, a)
		vpb := slices.Contains(priorityRelays, b)
		if vpa == vpb {
			return 1
		}
		if vpa && !vpb {
			return 1
		}
		return -1
	})

	return result, successRelays, nil
}
