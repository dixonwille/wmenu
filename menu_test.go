package wmenu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var newMenuCases = []string{"Testing this menu.", "", "!@#$%^&*()"}
var setSeperatorCases = []string{"", ".", ",", "~"}
var optionCases = []struct {
	name     string
	def      bool
	function func()
}{
	{"Options", true, func() { fmt.Println("testing option") }},
	{"", false, nil},
}
var actionCases = []func(Option){
	func(opt Option) { fmt.Println(opt) },
	nil,
}
var multipleActionCases = []func([]Option){
	func(opts []Option) { fmt.Println(opts) },
	nil,
}

func TestNewMenu(t *testing.T) {
	assert := assert.New(t)
	for _, c := range newMenuCases {
		menu := NewMenu(c)
		assert.Equal(c, menu.question)
		assert.Nil(menu.defaultFunction)
		assert.Nil(menu.options)
		assert.Equal(" ", menu.multiSeperator)
		assert.Nil(menu.multiFunction)
		assert.False(menu.loopOnInvalid)
		assert.False(menu.clear)
		assert.NotNil(menu.ui)
	}
}

func TestClearOnMenuRun(t *testing.T) {
	menu := NewMenu("Testing")
	menu.ClearOnMenuRun()
	assert.True(t, menu.clear)
}

func TestSetSeperator(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range setSeperatorCases {
		menu.SetSeperator(c)
		assert.Equal(t, c, menu.multiSeperator)
	}
}

func TestLoopOnInvalid(t *testing.T) {
	menu := NewMenu("Testing")
	menu.LoopOnInvalid()
	assert.True(t, menu.loopOnInvalid)
}

func TestOption(t *testing.T) {
	assert := assert.New(t)
	menu := NewMenu("Testing")
	for i, c := range optionCases {
		menu.Option(c.name, c.def, c.function)
		require.Equal(t, i+1, len(menu.options))
		assert.Equal(i, menu.options[i].ID)
		assert.Equal(c.name, menu.options[i].Text)
		assert.Equal(c.def, menu.options[i].isDefault)
		if c.function != nil {
			assert.NotNil(menu.options[i].function)
		} else {
			assert.Nil(menu.options[i].function)
		}
	}
}

func TestAction(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range actionCases {
		menu.Action(c)
		if c != nil {
			assert.NotNil(t, menu.defaultFunction)
		} else {
			assert.Nil(t, menu.defaultFunction)
		}
	}
}

func TestMultipleAction(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range multipleActionCases {
		menu.MultipleAction(c)
		if c != nil {
			assert.NotNil(t, menu.multiFunction)
		} else {
			assert.Nil(t, menu.multiFunction)
		}
	}
}
