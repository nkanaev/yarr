// +build macos windows

package platform

import (
	"github.com/getlantern/systray"
	"github.com/nkanaev/yarr/server"
	"github.com/skratchdot/open-golang/open"
)

func Start(s *server.Handler) {
	systrayOnReady := func() {
		systray.SetIcon(server.Icon)

		menuOpen := systray.AddMenuItem("Open", "")
		systray.AddSeparator()
		menuQuit := systray.AddMenuItem("Quit", "")

		go func() {
			for {
				select {
				case <-menuOpen.ClickedCh:
					open.Run("http://" + s.Addr)
				case <-menuQuit.ClickedCh:
					systray.Quit()
				}
			}
		}()

		s.Start()
	}
	systray.Run(systrayOnReady, nil)
}
