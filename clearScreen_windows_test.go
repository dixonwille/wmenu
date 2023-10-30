//go:build windows
// +build windows

package wmenu

import "testing"

func TestClearWindows(t *testing.T) {
	clearOs("windows")
}
