// +build !windows,!macos

package platform

import (
	"github.com/einschmidt/yarr/server"
)

func Start(s *server.Handler) {
	s.Start()
}
