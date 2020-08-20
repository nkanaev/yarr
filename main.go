package main

import (
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
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

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		logger.Fatal("Failed to create app config dir: ", err)
	}

	db, err := storage.New(storageFile, logger)
	if err != nil {
		logger.Fatal("Failed to initialise database: ", err)
	}

	srv := server.New(db, logger)
	srv.Start("127.0.0.1:8000")
}
