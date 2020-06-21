package main

import (
	"github.com/nkanaev/yarr/storage"
	"log"
)

func main() {
	storage, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	/*
	folder := storage.CreateFolder("foo")
	storage.RenameFolder(folder.Id, "bar")
	storage.ToggleFolderExpanded(folder.Id, false)
	*/
}
