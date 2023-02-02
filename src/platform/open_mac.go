//go:build darwin
// +build darwin

package platform

import "os/exec"

func open(input string) *exec.Cmd {
	return exec.Command("open", input)
}
