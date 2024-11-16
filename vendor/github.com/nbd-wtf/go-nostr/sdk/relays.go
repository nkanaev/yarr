package sdk

import (
	"github.com/nbd-wtf/go-nostr"
)

type RelayList = GenericList[Relay]

type Relay struct {
	URL    string
	Inbox  bool
	Outbox bool
}

func (r Relay) Value() string { return r.URL }

func parseRelayFromKind10002(tag nostr.Tag) (rl Relay, ok bool) {
	if u := tag.Value(); u != "" && tag[0] == "r" {
		if !nostr.IsValidRelayURL(u) {
			return rl, false
		}
		u := nostr.NormalizeURL(u)

		relay := Relay{
			URL: u,
		}

		if len(tag) == 2 {
			relay.Inbox = true
			relay.Outbox = true
		} else if tag[2] == "write" {
			relay.Outbox = true
		} else if tag[2] == "read" {
			relay.Inbox = true
		}

		return relay, true
	}

	return rl, false
}
