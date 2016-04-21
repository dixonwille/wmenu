package wmenu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() {

}
func TestNewMenu(t *testing.T) {
	assert := assert.New(t)
	menu := NewMenu("Testing this menu")
	assert.Equal("Testing this menu", menu.question)
	assert.Nil(menu.defaultFunction)
	assert.Nil(menu.options)
	assert.Equal(" ", menu.multiSeperator)
	assert.Nil(menu.multiFunction)
	assert.False(menu.loopOnInvalid)
	assert.False(menu.clear)
	assert.NotNil(menu.ui)
}

func TestClearOnMenuRun(t *testing.T) {
	menu := NewMenu("Testing")
	menu.ClearOnMenuRun()
	assert.True(t, menu.clear)
}

func TestSetSeperator(t *testing.T) {
	menu := NewMenu("Testing")
	menu.SetSeperator(",")
	assert.Equal(t, ",", menu.multiSeperator)
}

func TestLoopOnInvalid(t *testing.T) {
	menu := NewMenu("Testing")
	menu.LoopOnInvalid()
	assert.True(t, menu.loopOnInvalid)
}

func TestOption(t *testing.T) {
	assert := assert.New(t)
	menu := NewMenu("Testing")
	function := func() {
		fmt.Println("This is only a test.")
	}
	menu.Option("Option", true, function)
	require.True(t, len(menu.options) > 0)
	assert.Equal(0, menu.options[0].ID)
	assert.Equal("Option", menu.options[0].Text)
	assert.True(menu.options[0].isDefault)
	assert.NotNil(menu.options[0].function)
}

func TestAction(t *testing.T) {
	menu := NewMenu("Testing")
	function := func(opt Option) {
		fmt.Println(opt)
	}
	menu.Action(function)
	assert.NotNil(t, menu.defaultFunction)
}

func TestMultipleAction(t *testing.T) {
	menu := NewMenu("Testing")
	function := func(opts []Option) {
		fmt.Println(opts)
	}
	menu.MultipleAction(function)
	assert.NotNil(t, menu.multiFunction)
}
