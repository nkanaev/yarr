package cache_memory

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoCache[V any] struct {
	Cache *ristretto.Cache[string, V]
}

func New32[V any](max int64) *RistrettoCache[V] {
	cache, _ := ristretto.NewCache(&ristretto.Config[string, V]{
		NumCounters: max * 10,
		MaxCost:     max,
		BufferItems: 64,
		KeyToHash:   func(key string) (uint64, uint64) { return h32(key), 0 },
	})
	return &RistrettoCache[V]{Cache: cache}
}

func (s RistrettoCache[V]) Get(k string) (v V, ok bool) { return s.Cache.Get(k) }
func (s RistrettoCache[V]) Delete(k string)             { s.Cache.Del(k) }
func (s RistrettoCache[V]) Set(k string, v V) bool      { return s.Cache.Set(k, v, 1) }
func (s RistrettoCache[V]) SetWithTTL(k string, v V, d time.Duration) bool {
	return s.Cache.SetWithTTL(k, v, 1, d)
}

func h32(key string) uint64 {
	// we get an event id or pubkey as hex,
	// so just extract the last 8 bytes from it and turn them into a uint64
	return shortUint64(key)
}

func shortUint64(idOrPubkey string) uint64 {
	length := len(idOrPubkey)
	if length < 8 {
		return 0
	}
	b, err := hex.DecodeString(idOrPubkey[length-8:])
	if err != nil {
		return 0
	}
	return uint64(binary.BigEndian.Uint32(b))
}
