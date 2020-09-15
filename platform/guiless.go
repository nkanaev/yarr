// +build !windows,!macos

package platform

import (
	"github.com/nkanaev/yarr/server"
)

func Start(s *server.Handler) {
	s.Start()
}
