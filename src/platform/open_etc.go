//go:build linux || freebsd || openbsd

package platform

import "os/exec"

func open(input string) *exec.Cmd {
	return exec.Command("xdg-open", input)
}
