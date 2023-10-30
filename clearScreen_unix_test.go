//go:build darwin && linux
// +build darwin,linux

package wmenu

import (
	"os"
	"testing"
)

func init() {
	// Terminal is not set in many CI environments
	if os.Getenv("TERMINAL") == "" {
		os.Setenv("TERMINAL", "xterm")
	}
}

func TestClearLinux(t *testing.T) {
	clearOs("linux")
}

func TestClearDarwin(t *testing.T) {
	clearOs("darwin")
}
