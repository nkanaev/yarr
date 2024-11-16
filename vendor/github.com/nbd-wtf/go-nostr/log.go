package nostr

import (
	"log"
	"os"
)

var (
	// call SetOutput on InfoLogger to enable info logging
	InfoLogger = log.New(os.Stderr, "[go-nostr][info] ", log.LstdFlags)

	// call SetOutput on DebugLogger to enable debug logging
	DebugLogger = log.New(os.Stderr, "[go-nostr][debug] ", log.LstdFlags)
)
