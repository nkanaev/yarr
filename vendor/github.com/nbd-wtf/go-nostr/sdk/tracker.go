package sdk

import (
	"net/url"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/sdk/hints"
)

func (sys *System) TrackEventHints(ie nostr.RelayEvent) {
	if IsVirtualRelay(ie.Relay.URL) {
		return
	}
	if ie.Kind < 30000 && ie.Kind >= 20000 {
		return
	}

	switch ie.Kind {
	case nostr.KindRelayListMetadata:
		for _, tag := range ie.Tags {
			if len(tag) < 2 || tag[0] != "r" {
				continue
			}
			if len(tag) == 2 || (tag[2] == "" || tag[2] == "write") {
				sys.Hints.Save(ie.PubKey, tag[1], hints.LastInRelayList, ie.CreatedAt)
			}
		}
	case nostr.KindFollowList:
		sys.Hints.Save(ie.PubKey, ie.Relay.URL, hints.MostRecentEventFetched, ie.CreatedAt)

		for _, tag := range ie.Tags {
			if len(tag) < 3 {
				continue
			}
			if IsVirtualRelay(tag[2]) {
				continue
			}
			if p, err := url.Parse(tag[2]); err != nil || (p.Scheme != "wss" && p.Scheme != "ws") {
				continue
			}
			if tag[0] == "p" && nostr.IsValidPublicKey(tag[1]) {
				sys.Hints.Save(tag[1], tag[2], hints.LastInTag, ie.CreatedAt)
			}
		}
	case nostr.KindTextNote:
		sys.Hints.Save(ie.PubKey, ie.Relay.URL, hints.MostRecentEventFetched, ie.CreatedAt)

		for _, tag := range ie.Tags {
			if len(tag) < 3 {
				continue
			}
			if IsVirtualRelay(tag[2]) {
				continue
			}
			if p, err := url.Parse(tag[2]); err != nil || (p.Scheme != "wss" && p.Scheme != "ws") {
				continue
			}
			if tag[0] == "p" && nostr.IsValidPublicKey(tag[1]) {
				sys.Hints.Save(tag[1], tag[2], hints.LastInTag, ie.CreatedAt)
			}
		}

		for _, ref := range ParseReferences(ie.Event) {
			if ref.Profile != nil {
				for _, relay := range ref.Profile.Relays {
					if IsVirtualRelay(relay) {
						continue
					}
					if p, err := url.Parse(relay); err != nil || (p.Scheme != "wss" && p.Scheme != "ws") {
						continue
					}
					if nostr.IsValidPublicKey(ref.Profile.PublicKey) {
						sys.Hints.Save(ref.Profile.PublicKey, relay, hints.LastInNprofile, ie.CreatedAt)
					}
				}
			} else if ref.Event != nil && nostr.IsValidPublicKey(ref.Event.Author) {
				for _, relay := range ref.Event.Relays {
					if IsVirtualRelay(relay) {
						continue
					}
					if p, err := url.Parse(relay); err != nil || (p.Scheme != "wss" && p.Scheme != "ws") {
						continue
					}
					sys.Hints.Save(ref.Event.Author, relay, hints.LastInNevent, ie.CreatedAt)
				}
			}
		}
	}
}
