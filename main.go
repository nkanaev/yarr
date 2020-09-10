package main

import (
	"github.com/getlantern/systray"
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
	"github.com/skratchdot/open-golang/open"
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

	addr := "127.0.0.1:7070"

	systrayOnReady := func() {
		systray.SetIcon(server.Icon)

		menuOpen := systray.AddMenuItem("Open", "")
		systray.AddSeparator()
		menuQuit := systray.AddMenuItem("Quit", "")

		go func() {
			for {
				select {
				case <-menuOpen.ClickedCh:
					open.Run("http://" + addr)
				case <-menuQuit.ClickedCh:
					systray.Quit()
				}
			}
		}()
		srv := server.New(db, logger)
		srv.Start(addr)
	}
	systray.Run(systrayOnReady, nil)
}
