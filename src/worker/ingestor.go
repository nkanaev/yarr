package worker

import (
	"context"
	"errors"
	"io"

	"github.com/nkanaev/yarr/src/parser/opml"
	"github.com/nkanaev/yarr/src/server"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

var (
	_ server.FeedScheduler = (*Worker)(nil)
	_ server.FeedIngestor  = (*Worker)(nil)
)

func (w *Worker) AddFeed(ctx context.Context, store storage.Storage, link server.FeedLink) (server.AddFeedResult, error) {
	result, err := DiscoverFeed(w.client, link.URL)
	if err != nil {
		return server.AddFeedResult{}, err
	}

	if len(result.Sources) > 0 {
		choices := make([]server.FeedLink, len(result.Sources))
		for i, source := range result.Sources {
			choices[i] = server.FeedLink{
				URL:           source.URL,
				Title:         source.Title,
				TitleOverride: source.TitleOverride,
			}
		}
		return server.AddFeedResult{Choices: choices}, nil
	}

	if result.Feed == nil {
		return server.AddFeedResult{}, errors.New("no feed found")
	}

	title := result.Feed.Title
	if link.TitleOverride != "" {
		title = link.TitleOverride
	} else if link.Title != "" {
		title = link.Title
	}

	feed := store.CreateFeed(model.CreateFeedParams{
		Title:    title,
		Link:     result.Feed.SiteURL,
		FeedLink: result.FeedLink,
	})
	items := ConvertItems(result.Feed.Items, *feed)
	if len(items) > 0 {
		store.CreateItems(items)
	}
	return server.AddFeedResult{Feed: feed}, nil
}

func (w *Worker) AddOPML(ctx context.Context, store storage.Storage, file io.Reader) error {
	doc, err := opml.Parse(file)
	if err != nil {
		return err
	}

	for _, f := range doc.Feeds {
		store.CreateFeed(model.CreateFeedParams{
			Title:    f.Title,
			Link:     f.SiteUrl,
			FeedLink: f.FeedUrl,
		})
	}
	for _, f := range doc.Folders {
		folder := store.CreateFolder(f.Title)
		for _, ff := range f.AllFeeds() {
			store.CreateFeed(model.CreateFeedParams{
				Title:    ff.Title,
				Link:     ff.SiteUrl,
				FeedLink: ff.FeedUrl,
				FolderID: &folder.Id,
			})
		}
	}

	w.RefreshFeeds()
	return nil
}
