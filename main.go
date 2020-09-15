package main

import (
	"flag"
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
	"github.com/nkanaev/yarr/platform"
	"log"
	"os"
	"path/filepath"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	configPath, err := os.UserConfigDir()
	if err != nil {
		logger.Fatal("Failed to get config dir: ", err)
	}
	storagePath := filepath.Join(configPath, "yarr")
	storageFile := filepath.Join(storagePath, "storage.db")

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		logger.Fatal("Failed to create app config dir: ", err)
	}

	db, err := storage.New(storageFile, logger)
	if err != nil {
		logger.Fatal("Failed to initialise database: ", err)
	}

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:7070", "address to run server on")
	flag.Parse()

	srv := server.New(db, logger, addr)
	platform.Start(srv)
}
