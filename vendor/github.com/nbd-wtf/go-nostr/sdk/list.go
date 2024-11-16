package sdk

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/sdk/cache"
)

type GenericList[I TagItemWithValue] struct {
	PubKey string       `json:"-"` // must always be set otherwise things will break
	Event  *nostr.Event `json:"-"` // may be empty if a contact list event wasn't found

	Items []I
}

type TagItemWithValue interface {
	Value() string
}

var (
	genericListMutexes = [24]sync.Mutex{}
	valueWasJustCached = [24]bool{}
)

func fetchGenericList[I TagItemWithValue](
	sys *System,
	ctx context.Context,
	pubkey string,
	kind int,
	parseTag func(nostr.Tag) (I, bool),
	cache cache.Cache32[GenericList[I]],
	skipFetch bool,
) (fl GenericList[I], fromInternal bool) {
	// we have 24 mutexes, so we can load up to 24 lists at the same time, but if we do the same exact
	// call that will do it only once, the subsequent ones will wait for a result to be cached
	// and then return it from cache -- 13 is an arbitrary index for the pubkey
	lockIdx := (int(pubkey[13]) + kind) % 24
	genericListMutexes[lockIdx].Lock()

	if valueWasJustCached[lockIdx] {
		// this ensures the cache has had time to commit the values
		// so we don't repeat a fetch immediately after the other
		valueWasJustCached[lockIdx] = false
		time.Sleep(time.Millisecond * 10)
	}

	defer genericListMutexes[lockIdx].Unlock()

	if v, ok := cache.Get(pubkey); ok {
		return v, true
	}

	v := GenericList[I]{PubKey: pubkey}

	events, _ := sys.StoreRelay.QuerySync(ctx, nostr.Filter{Kinds: []int{kind}, Authors: []string{pubkey}})
	if len(events) != 0 {
		items := parseItemsFromEventTags(events[0], parseTag)
		v.Event = events[0]
		v.Items = items
		cache.SetWithTTL(pubkey, v, time.Hour*6)
		valueWasJustCached[lockIdx] = true
		return v, true
	}

	if !skipFetch {
		thunk := sys.replaceableLoaders[kind].Load(ctx, pubkey)
		evt, err := thunk()
		if err == nil {
			items := parseItemsFromEventTags(evt, parseTag)
			v.Items = items
			sys.StoreRelay.Publish(ctx, *evt)
		}
		cache.SetWithTTL(pubkey, v, time.Hour*6)
		valueWasJustCached[lockIdx] = true
	}

	return v, false
}

func parseItemsFromEventTags[I TagItemWithValue](
	evt *nostr.Event,
	parseTag func(nostr.Tag) (I, bool),
) []I {
	result := make([]I, 0, len(evt.Tags))
	for _, tag := range evt.Tags {
		item, ok := parseTag(tag)
		if ok {
			// check if this already exists before adding
			if slices.IndexFunc(result, func(i I) bool { return i.Value() == item.Value() }) == -1 {
				result = append(result, item)
			}
		}
	}
	return result
}
