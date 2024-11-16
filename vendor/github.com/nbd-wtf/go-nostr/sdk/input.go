package sdk

import (
	"context"
	"encoding/hex"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip05"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// InputToProfile turns any npub/nprofile/hex/nip05 input into a ProfilePointer (or nil).
func InputToProfile(ctx context.Context, input string) *nostr.ProfilePointer {
	// handle if it is a hex string
	if len(input) == 64 {
		if _, err := hex.DecodeString(input); err == nil {
			return &nostr.ProfilePointer{PublicKey: input}
		}
	}

	// handle nip19 codes, if that's the case
	prefix, data, _ := nip19.Decode(input)
	switch prefix {
	case "npub":
		input = data.(string)
		return &nostr.ProfilePointer{PublicKey: input}
	case "nprofile":
		pp := data.(nostr.ProfilePointer)
		return &pp
	}

	// handle nip05 ids, if that's the case
	pp, _ := nip05.QueryIdentifier(ctx, input)
	if pp != nil {
		return pp
	}

	return nil
}

// InputToEventPointer turns any note/nevent/hex input into a EventPointer (or nil).
func InputToEventPointer(input string) *nostr.EventPointer {
	// handle if it is a hex string
	if len(input) == 64 {
		if _, err := hex.DecodeString(input); err == nil {
			return &nostr.EventPointer{ID: input}
		}
	}

	// handle nip19 codes, if that's the case
	prefix, data, _ := nip19.Decode(input)
	switch prefix {
	case "note":
		if input, ok := data.(string); ok {
			return &nostr.EventPointer{ID: input}
		}
	case "nevent":
		if ep, ok := data.(nostr.EventPointer); ok {
			return &ep
		}
	}

	// handle nip05 ids, if that's the case
	return nil
}
