package sdk

import (
	"context"
	"net/url"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

type FollowList = GenericList[Follow]

type Follow struct {
	Pubkey  string
	Relay   string
	Petname string
}

func (f Follow) Value() string { return f.Pubkey }

func (sys *System) FetchFollowList(ctx context.Context, pubkey string) FollowList {
	fl, _ := fetchGenericList(sys, ctx, pubkey, 3, parseFollow, sys.FollowListCache, false)
	return fl
}

func parseFollow(tag nostr.Tag) (fw Follow, ok bool) {
	if len(tag) < 2 {
		return fw, false
	}
	if tag[0] != "p" {
		return fw, false
	}

	fw.Pubkey = tag[1]
	if !nostr.IsValidPublicKey(fw.Pubkey) {
		return fw, false
	}

	if len(tag) > 2 {
		if _, err := url.Parse(tag[2]); err == nil {
			fw.Relay = nostr.NormalizeURL(tag[2])
		}
		if len(tag) > 3 {
			fw.Petname = strings.TrimSpace(tag[3])
		}
		return fw, true
	}

	return fw, false
}
