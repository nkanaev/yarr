//go:build windows
// +build windows

package platform

import (
	"os"
	"os/exec"
	"path/filepath"
)

func open(input string) *exec.Cmd {
	rundll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
	return exec.Command(rundll32, "url.dll,FileProtocolHandler", input)
}
