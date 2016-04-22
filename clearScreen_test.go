package wmenu

import "testing"

func TestClear(t *testing.T) {
	Clear()
}

func TestClearLinux(t *testing.T) {
	clearOs("linux")
}

func TestClearDarwin(t *testing.T) {
	clearOs("darwin")
}

func TestClearWindows(t *testing.T) {
	clearOs("windows")
}

func clearOs(os string) {
	value, ok := clear[os]
	if ok {
		value()
	}
}
