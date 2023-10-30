package wmenu

import (
	"testing"
)

func TestClear(t *testing.T) {
	Clear()
}

func clearOs(os string) {
	value, ok := clear[os]
	if ok {
		value()
	}
}
