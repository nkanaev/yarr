package main

import (
	"github.com/nkanaev/yarr/storage"
	"log"
)

func main() {
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(store)
	/*
	folder := storage.CreateFolder("foo")
	storage.RenameFolder(folder.Id, "bar")
	storage.ToggleFolderExpanded(folder.Id, false)
	*/
	/*
	feed := storage.CreateFeed(
		"title", "description", "link", "feedlink", "icon", 1)
	storage.RenameFeed(feed.Id, "newtitle")
	*/
	/*
	items := make([]storage.Item, 3, 3)
	items = append(items, storage.Item{
		Id: "id",
		FeedId: 0,
		Title: "title",
		Link: "link",
		Description: "description",
		Content: "content",
		Author: "author",
		Date: 1,
		DateUpdated: 1,
		Status: storage.UNREAD,
		Image: "image",
	})
	items = append(items, storage.Item{
		Id: "id2",
		FeedId: 0,
		Title: "title",
		Link: "link",
		Description: "description",
		Content: "content",
		Author: "author",
		Date: 1,
		DateUpdated: 50,
		Status: storage.UNREAD,
		Image: "image",
	})
	items = append(items, storage.Item{
		Id: "id",
		FeedId: 0,
		Title: "title",
		Link: "link",
		Description: "description",
		Content: "content",
		Author: "author",
		Date: 1,
		DateUpdated: 100,
		Status: storage.UNREAD,
		Image: "image",
	})
	log.Print(store.CreateItems(items))
	*/
}
