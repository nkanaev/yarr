package main

import (
	"github.com/nkanaev/yarr/storage"
	//"github.com/nkanaev/yarr/worker"
	"github.com/nkanaev/yarr/server"
	"log"
)

func main() {
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(store)
	/*
	folder := store.CreateFolder("foo")
	store.RenameFolder(folder.Id, "bar")
	store.ToggleFolderExpanded(folder.Id, false)
	log.Print(store.ListFolders())
	*/
	/*
	feed := store.CreateFeed(
		"title", "description", "link", "feedlink", "icon", 1)
	store.RenameFeed(feed.Id, "newtitle")
	log.Print(store.ListFeeds())
	*/
	/*;
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
	log.Print(store.ListItems())
	*/
	/*
	log.Print(worker.FindFeeds("https://horriblesubs.info/"))
	log.Print(worker.FindFeeds("http://daringfireball.net/"))
	*/
	srv := server.New()
	srv.ListenAndServe()
}
