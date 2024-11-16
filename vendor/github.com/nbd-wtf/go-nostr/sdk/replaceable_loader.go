package sdk

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/nbd-wtf/go-nostr"
)

type EventResult dataloader.Result[*nostr.Event]

func (sys *System) initializeDataloaders() {
	sys.replaceableLoaders = make(map[int]*dataloader.Loader[string, *nostr.Event])
	for _, kind := range []int{0, 3, 10000, 10001, 10002, 10003, 10004, 10005, 10006, 10007, 10015, 10030} {
		sys.replaceableLoaders[kind] = sys.createReplaceableDataloader(kind)
	}
}

func (sys *System) createReplaceableDataloader(kind int) *dataloader.Loader[string, *nostr.Event] {
	return dataloader.NewBatchedLoader(
		func(_ context.Context, pubkeys []string) []*dataloader.Result[*nostr.Event] {
			return sys.batchLoadReplaceableEvents(kind, pubkeys)
		},
		dataloader.WithBatchCapacity[string, *nostr.Event](60),
		dataloader.WithClearCacheOnBatch[string, *nostr.Event](),
		dataloader.WithWait[string, *nostr.Event](time.Millisecond*350),
	)
}

func (sys *System) batchLoadReplaceableEvents(
	kind int,
	pubkeys []string,
) []*dataloader.Result[*nostr.Event] {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	batchSize := len(pubkeys)
	results := make([]*dataloader.Result[*nostr.Event], batchSize)
	keyPositions := make(map[string]int)          // { [pubkey]: slice_index }
	relayFilters := make(map[string]nostr.Filter) // { [relayUrl]: filter }

	wg := sync.WaitGroup{}
	wg.Add(len(pubkeys))
	cm := sync.Mutex{}

	for i, pubkey := range pubkeys {
		// build batched queries for the external relays
		keyPositions[pubkey] = i // this is to help us know where to save the result later

		go func(i int, pubkey string) {
			defer wg.Done()

			// if we're attempting this query with a short key (last 8 characters), stop here
			if len(pubkey) != 64 {
				results[i] = &dataloader.Result[*nostr.Event]{
					Error: fmt.Errorf("won't proceed to query relays with a shortened key (%d)", kind),
				}
				return
			}

			// save attempts here so we don't try the same failed query over and over
			if doItNow := doThisNotMoreThanOnceAnHour("repl:" + strconv.Itoa(kind) + pubkey); !doItNow {
				results[i] = &dataloader.Result[*nostr.Event]{
					Error: fmt.Errorf("last attempt failed, waiting more to try again"),
				}
				return
			}

			// gather relays we'll use for this pubkey
			relays := sys.determineRelaysToQuery(ctx, pubkey, kind)

			// by default we will return an error (this will be overwritten when we find an event)
			results[i] = &dataloader.Result[*nostr.Event]{
				Error: fmt.Errorf("couldn't find a kind %d event anywhere %v", kind, relays),
			}

			cm.Lock()
			for _, relay := range relays {
				// each relay will have a custom filter
				filter, ok := relayFilters[relay]
				if !ok {
					filter = nostr.Filter{
						Kinds:   []int{kind},
						Authors: make([]string, 0, batchSize-i /* this and all pubkeys after this can be added */),
					}
				}
				filter.Authors = append(filter.Authors, pubkey)
				relayFilters[relay] = filter
			}
			cm.Unlock()
		}(i, pubkey)
	}

	// query all relays with the prepared filters
	wg.Wait()
	multiSubs := sys.batchReplaceableRelayQueries(ctx, relayFilters)
	for {
		select {
		case evt, more := <-multiSubs:
			if !more {
				return results
			}

			// insert this event at the desired position
			pos := keyPositions[evt.PubKey] // @unchecked: it must succeed because it must be a key we passed
			if results[pos].Data == nil || results[pos].Data.CreatedAt < evt.CreatedAt {
				results[pos] = &dataloader.Result[*nostr.Event]{Data: evt}
			}
		case <-ctx.Done():
			return results
		}
	}
}

func (sys *System) determineRelaysToQuery(ctx context.Context, pubkey string, kind int) []string {
	relays := make([]string, 0, 10)

	// search in specific relays for user
	if kind == 10002 {
		// prevent infinite loops by jumping directly to this
		relays = sys.Hints.TopN(pubkey, 3)
	} else if kind == 0 {
		// leave room for one hardcoded relay because people are stupid
		relays = sys.FetchOutboxRelays(ctx, pubkey, 2)
	} else {
		relays = sys.FetchOutboxRelays(ctx, pubkey, 3)
	}

	// use a different set of extra relays depending on the kind
	for i := 0; i < 3-len(relays); i++ {
		var next string
		switch kind {
		case 0:
			next = sys.MetadataRelays.Next()
		case 3:
			next = sys.FollowListRelays.Next()
		case 10002:
			next = sys.RelayListRelays.Next()
		default:
			next = sys.FallbackRelays.Next()
		}

		if !slices.Contains(relays, next) {
			relays = append(relays, next)
		}
	}

	return relays
}

// batchReplaceableRelayQueries subscribes to multiple relays using a different filter for each and returns
// a single channel with all results. it closes on EOSE or when all the expected events were returned.
//
// the number of expected events is given by the number of pubkeys in the .Authors filter field.
// because of that, batchReplaceableRelayQueries is only suitable for querying replaceable events -- and
// care must be taken to not include the same pubkey more than once in the filter .Authors array.
func (sys *System) batchReplaceableRelayQueries(
	ctx context.Context,
	relayFilters map[string]nostr.Filter,
) <-chan *nostr.Event {
	all := make(chan *nostr.Event)

	wg := sync.WaitGroup{}
	wg.Add(len(relayFilters))
	for url, filter := range relayFilters {
		go func(url string, filter nostr.Filter) {
			defer wg.Done()
			n := len(filter.Authors)

			ctx, cancel := context.WithTimeout(ctx, time.Millisecond*450+time.Millisecond*50*time.Duration(n))
			defer cancel()

			received := 0
			for ie := range sys.Pool.SubManyEose(ctx, []string{url}, nostr.Filters{filter}, nostr.WithLabel("repl")) {
				all <- ie.Event
				received++
				if received >= n {
					// we got all events we asked for, unless the relay is shitty and sent us two from the same
					return
				}
			}
		}(url, filter)
	}

	go func() {
		wg.Wait()
		close(all)
	}()

	return all
}
