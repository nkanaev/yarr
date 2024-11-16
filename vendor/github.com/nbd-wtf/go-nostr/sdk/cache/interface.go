package cache

import "time"

type Cache32[V any] interface {
	Get(k string) (v V, ok bool)
	Delete(k string)
	Set(k string, v V) bool
	SetWithTTL(k string, v V, d time.Duration) bool
}
