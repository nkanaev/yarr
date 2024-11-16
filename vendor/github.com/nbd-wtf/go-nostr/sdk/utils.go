package sdk

import (
	"strings"
	"sync"
	"time"
)

var (
	_dtnmtoah     map[string]time.Time
	_dtnmtoahLock sync.Mutex
)

// IsVirtualRelay returns true if the given normalized relay URL shouldn't be considered for outbox-model calculations.
func IsVirtualRelay(url string) bool {
	if len(url) < 6 {
		// this is just invalid
		return true
	}

	if strings.HasPrefix(url, "wss://feeds.nostr.band") ||
		strings.HasPrefix(url, "wss://filter.nostr.wine") ||
		strings.HasPrefix(url, "wss://cache") {
		return true
	}

	return false
}
