package sdk

import (
	"context"

	"github.com/fiatjaf/eventstore"
	"github.com/fiatjaf/eventstore/nullstore"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/sdk/cache"
	cache_memory "github.com/nbd-wtf/go-nostr/sdk/cache/memory"
	"github.com/nbd-wtf/go-nostr/sdk/hints"
	memory_hints "github.com/nbd-wtf/go-nostr/sdk/hints/memory"
)

type System struct {
	RelayListCache   cache.Cache32[RelayList]
	FollowListCache  cache.Cache32[FollowList]
	MetadataCache    cache.Cache32[ProfileMetadata]
	Hints            hints.HintsDB
	Pool             *nostr.SimplePool
	RelayListRelays  *RelayStream
	FollowListRelays *RelayStream
	MetadataRelays   *RelayStream
	FallbackRelays   *RelayStream
	JustIDRelays     *RelayStream
	UserSearchRelays *RelayStream
	NoteSearchRelays *RelayStream
	Store            eventstore.Store

	StoreRelay nostr.RelayStore

	replaceableLoaders   map[int]*dataloader.Loader[string, *nostr.Event]
	outboxShortTermCache cache.Cache32[[]string]
}

type SystemModifier func(sys *System)

type RelayStream struct {
	URLs   []string
	serial int
}

func NewRelayStream(urls ...string) *RelayStream {
	return &RelayStream{URLs: urls, serial: -1}
}

func (rs *RelayStream) Next() string {
	rs.serial++
	return rs.URLs[rs.serial%len(rs.URLs)]
}

func NewSystem(mods ...SystemModifier) *System {
	sys := &System{
		RelayListCache:   cache_memory.New32[RelayList](1000),
		FollowListCache:  cache_memory.New32[FollowList](1000),
		MetadataCache:    cache_memory.New32[ProfileMetadata](1000),
		RelayListRelays:  NewRelayStream("wss://purplepag.es", "wss://user.kindpag.es", "wss://relay.nos.social"),
		FollowListRelays: NewRelayStream("wss://purplepag.es", "wss://user.kindpag.es", "wss://relay.nos.social"),
		MetadataRelays:   NewRelayStream("wss://purplepag.es", "wss://user.kindpag.es", "wss://relay.nos.social"),
		FallbackRelays: NewRelayStream(
			"wss://relay.damus.io",
			"wss://nostr.mom",
			"wss://nos.lol",
			"wss://mostr.pub",
			"wss://relay.nostr.band",
		),
		JustIDRelays: NewRelayStream(
			"wss://cache2.primal.net/v1",
			"wss://relay.noswhere.com",
			"wss://relay.nostr.band",
		),
		UserSearchRelays: NewRelayStream(
			"wss://search.nos.today",
			"wss://nostr.wine",
			"wss://relay.nostr.band",
		),
		NoteSearchRelays: NewRelayStream(
			"wss://nostr.wine",
			"wss://relay.nostr.band",
			"wss://search.nos.today",
		),
		Hints: memory_hints.NewHintDB(),

		outboxShortTermCache: cache_memory.New32[[]string](1000),
	}

	sys.Pool = nostr.NewSimplePool(context.Background(),
		nostr.WithEventMiddleware(sys.TrackEventHints),
		nostr.WithPenaltyBox(),
	)

	for _, mod := range mods {
		mod(sys)
	}

	if sys.Store == nil {
		sys.Store = &nullstore.NullStore{}
		sys.Store.Init()
	}
	sys.StoreRelay = eventstore.RelayWrapper{Store: sys.Store}

	sys.initializeDataloaders()

	return sys
}

func (sys *System) Close() {}

func WithHintsDB(hdb hints.HintsDB) SystemModifier {
	return func(sys *System) {
		sys.Hints = hdb
	}
}

func WithRelayListRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.RelayListRelays.URLs = list
	}
}

func WithMetadataRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.MetadataRelays.URLs = list
	}
}

func WithFollowListRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.FollowListRelays.URLs = list
	}
}

func WithFallbackRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.FallbackRelays.URLs = list
	}
}

func WithJustIDRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.JustIDRelays.URLs = list
	}
}

func WithUserSearchRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.UserSearchRelays.URLs = list
	}
}

func WithNoteSearchRelays(list []string) SystemModifier {
	return func(sys *System) {
		sys.NoteSearchRelays.URLs = list
	}
}

func WithStore(store eventstore.Store) SystemModifier {
	return func(sys *System) {
		sys.Store = store
	}
}

func WithRelayListCache(cache cache.Cache32[RelayList]) SystemModifier {
	return func(sys *System) {
		sys.RelayListCache = cache
	}
}

func WithFollowListCache(cache cache.Cache32[FollowList]) SystemModifier {
	return func(sys *System) {
		sys.FollowListCache = cache
	}
}

func WithMetadataCache(cache cache.Cache32[ProfileMetadata]) SystemModifier {
	return func(sys *System) {
		sys.MetadataCache = cache
	}
}
