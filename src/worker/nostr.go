package worker

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/github-tijlxyz/goldmark-nostr"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip05"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/nbd-wtf/go-nostr/sdk"
	"github.com/nkanaev/yarr/src/parser"
	"github.com/yuin/goldmark"
)

var (
	nostrSdk *sdk.System
)

func initializeNostr() {
	nostrSdk = sdk.NewSystem(
		sdk.WithRelayListRelays([]string{
			"wss://nos.lol", "wss://nostr.mom", "wss://nostr.bitcoiner.social", "wss://relay.damus.io", "wss://nostr-pub.wellorder.net"}, // some standard relays
		),
	)
}

// Main function for checking if the url is a nostr url
func isItNostr(ctx context.Context, url string) (bool, *sdk.ProfileMetadata) {
	if nostrSdk == nil {
		initializeNostr()
	}

	// check for nostr url prefixes
	if strings.HasPrefix(url, "nostr://") {
		url = url[8:]
	} else if strings.HasPrefix(url, "nostr:") {
		url = url[6:]
	} else {
		// only accept nostr: or nostr:// urls for now
		return false, nil
	}

	// check for npub or nprofile
	if prefix, _, err := nip19.Decode(url); err == nil {
		if prefix == "nprofile" || prefix == "npub" {
			profile, err := nostrSdk.FetchProfileFromInput(ctx, url)
			if err != nil {
				return false, nil
			}
			return true, &profile
		}
	}

	// only do nip05 check when nostr prefix
	if nip05.IsValidIdentifier(url) {
		profile, err := nostrSdk.FetchProfileFromInput(ctx, url)
		if err != nil {
			return false, nil
		}
		return true, &profile
	}

	return false, nil
}

// Load the feed and items
func discoverNostr(candidateUrl string) (bool, *DiscoverResult) {
	ctx := context.Background()

	yes, profile := isItNostr(ctx, candidateUrl)
	if yes {

		nprofile := profile.Nprofile(ctx, nostrSdk, 3)

		// get some feed items
		_, items, err := nostrListItems(candidateUrl)
		if err != nil {
			items = []parser.Item{}
		}

		return true, &DiscoverResult{
			FeedLink: fmt.Sprintf("nostr:%s", nprofile),
			Feed: &parser.Feed{
				Title:   profile.Name,
				SiteURL: fmt.Sprintf("nostr:%s", nprofile),
				Items:   items,
			},
			Sources: []FeedSource{},
		}
	}
	return false, nil
}

// Load the nostr favicon
func nostrLookForFavicon(link string) (bool, *[]byte, error) {
	ctx := context.Background()
	yes, profile := isItNostr(ctx, link)
	if yes {
		favicon, err := getIcons([]string{profile.Picture})
		return true, favicon, err

	}
	return false, nil, nil
}

// Update nostr feed items
func nostrListItems(f string) (bool, []parser.Item, error) {
	ctx := context.Background()

	if yes, profile := isItNostr(ctx, f); yes {
		relays := nostrSdk.FetchOutboxRelays(ctx, profile.PubKey, 3)
		evchan := nostrSdk.Pool.SubManyEose(ctx, relays, nostr.Filters{
			{
				Authors: []string{profile.PubKey},
				Kinds:   []int{nostr.KindArticle},
				Limit:   32,
			},
		})
		feedItems := []parser.Item{}
		for event := range evchan {

			publishedAt := event.CreatedAt.Time()
			if publishedAtTag := event.Tags.GetFirst([]string{"published_at"}); publishedAtTag != nil && len(*publishedAtTag) >= 2 {
				i, err := strconv.ParseInt((*publishedAtTag)[1], 10, 64)
				if err != nil {
					publishedAt = time.Unix(i, 0)
				}
			}

			naddr, err := nip19.EncodeEntity(event.PubKey, event.Kind, event.Tags.GetD(), relays)
			if err != nil {
				continue
			}

			title := ""
			titleTag := event.Tags.GetFirst([]string{"title"})
			if titleTag != nil && len(*titleTag) >= 2 {
				title = (*titleTag)[1]
			} else {
				continue
			}

			image := ""
			imageTag := event.Tags.GetFirst([]string{"image"})
			if imageTag != nil && len(*imageTag) >= 2 {
				image = (*imageTag)[1]
			}

			// format content from markdown to html
			md := goldmark.New(goldmark.WithExtensions(extension.NewNostr()))
			var buf bytes.Buffer
			if err := md.Convert([]byte(event.Content), &buf); err != nil {
				continue
			}

			feedItems = append(feedItems, parser.Item{
				GUID:     fmt.Sprintf("nostr:%s:%s", event.PubKey, event.Tags.GetD()),
				Date:     publishedAt,
				URL:      fmt.Sprintf("nostr:%s", naddr),
				Content:  buf.String(),
				Title:    title,
				ImageURL: image,
			})

		}

		return true, feedItems, nil
	}

	return false, nil, nil
}
