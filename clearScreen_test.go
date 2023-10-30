package wmenu

import (
	"testing"
)

func TestClear(t *testing.T) {
	Clear()
}

// Used specifically per OS
//
//nolint:unused
func clearOs(os string) {
	value, ok := clear[os]
	if ok {
		value()
	}
}
