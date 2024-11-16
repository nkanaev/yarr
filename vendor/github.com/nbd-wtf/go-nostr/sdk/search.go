package sdk

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

func (sys *System) SearchUsers(ctx context.Context, query string) []ProfileMetadata {
	limit := 10
	profiles := make([]ProfileMetadata, 0, limit*len(sys.UserSearchRelays.URLs))

	for ie := range sys.Pool.SubManyEose(ctx, sys.UserSearchRelays.URLs, nostr.Filters{
		{
			Search: query,
			Limit:  limit,
		},
	}, nostr.WithLabel("search")) {
		m, _ := ParseMetadata(ie.Event)
		profiles = append(profiles, m)
	}

	return profiles
}
