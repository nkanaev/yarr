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
	var addr, storageFile string
	flag.StringVar(&addr, "addr", "127.0.0.1:7070", "address to run server on")
	flag.StringVar(&storageFile, "db", "", "storage file path")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	configPath, err := os.UserConfigDir()
	if err != nil {
		logger.Fatal("Failed to get config dir: ", err)
	}

	if storageFile == "" {
		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			logger.Fatal("Failed to create app config dir: ", err)
		}
		storageFile = filepath.Join(storagePath, "storage.db")
	}

	db, err := storage.New(storageFile, logger)
	if err != nil {
		logger.Fatal("Failed to initialise database: ", err)
	}

	srv := server.New(db, logger, addr)
	platform.Start(srv)
}
