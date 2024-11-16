package hints

import "github.com/nbd-wtf/go-nostr"

type HintsDB interface {
	TopN(pubkey string, n int) []string
	Save(pubkey string, relay string, key HintKey, score nostr.Timestamp)
	PrintScores()
}
