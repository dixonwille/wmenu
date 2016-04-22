package wmenu

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/dixonwille/wlog"
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
var addColorCases = []struct {
	opt wlog.Color
	que wlog.Color
	res wlog.Color
	err wlog.Color
}{
	{wlog.None, wlog.None, wlog.None, wlog.None},
	{wlog.Red, wlog.Green, wlog.Blue, wlog.Black},
	{wlog.Yellow, wlog.White, wlog.Cyan, wlog.Magenta},
	{wlog.BrightBlack, wlog.BrightBlue, wlog.BrightCyan, wlog.BrightGreen},
	{wlog.BrightRed, wlog.BrightWhite, wlog.BrightYellow, wlog.BrightMagenta},
}

func Example_simple() {
	reader := strings.NewReader("1\r\n") //Simulates the user typing "1" and hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	actFunc := func(opt Option) {
		fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.Option("Option 0", true, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", true, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 1 has an id of 1.
}

func Example_simpleDefault() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Fprint(os.Stdout, "Option 0 was chosen.")
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", true, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 0 was chosen.
}

func Example_multiple() {
	reader := strings.NewReader("1,2\r\n") //Simulates the user typing "1,2" and hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	multiFunc := func(opts []Option) {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.MultipleAction(multiFunc)
	menu.SetSeparator(",")
	menu.Option("Option 0", true, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", true, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 1 has an id of 1.
	//Option 2 has an id of 2.
}

func Example_multipleDefault() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	multiFunc := func(opts []Option) {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.MultipleAction(multiFunc)
	menu.Option("Option 0", true, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", true, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 0 has an id of 0.
	//Option 2 has an id of 2.
}

func Example_errorNoResponse() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", false, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", false, nil)
	err := menu.Run()
	if err != nil {
		if IsNoResponseErr(err) {
			fmt.Println("We caught the err: " + err.Error())
		} else {
			log.Fatal(err)
		}
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//We caught the err: No response
}

func Example_errorInvalid() {
	reader := strings.NewReader("3\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", false, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", false, nil)
	err := menu.Run()
	if err != nil {
		if IsInvalidErr(err) {
			fmt.Println("We caught the err: " + err.Error())
		} else {
			log.Fatal(err)
		}
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//We caught the err: Invalid response: 3
}

func Example_errorTooMany() {
	reader := strings.NewReader("1 2\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", false, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", false, nil)
	err := menu.Run()
	if err != nil {
		if IsTooManyErr(err) {
			fmt.Println("We caught the err: " + err.Error())
		} else {
			log.Fatal(err)
		}
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//We caught the err: Too many responses
}

func Example_errorDuplicate() {
	reader := strings.NewReader("1 1\r\n") //Simulates the user hitting the [enter] key
	optFunc := func() {
		fmt.Println("Option 0 was chosen.")
	}
	multiFunc := func(opts []Option) {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.MultipleAction(multiFunc)
	menu.Option("Option 0", false, optFunc)
	menu.Option("Option 1", false, nil)
	menu.Option("Option 2", false, nil)
	err := menu.Run()
	if err != nil {
		if IsDuplicateErr(err) {
			fmt.Println("We caught the err: " + err.Error())
		} else {
			log.Fatal(err)
		}
	}
	//Output:
	//0) Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//We caught the err: Duplicated response: 1

}

func TestNewMenu(t *testing.T) {
	assert := assert.New(t)
	for _, c := range newMenuCases {
		menu := NewMenu(c)
		assert.Equal(c, menu.question)
		assert.Nil(menu.defaultFunction)
		assert.Nil(menu.options)
		assert.Equal(" ", menu.multiSeparator)
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
		menu.SetSeparator(c)
		assert.Equal(t, c, menu.multiSeparator)
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

func TestAddColor(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range addColorCases {
		menu.AddColor(c.opt, c.que, c.res, c.err)
		//Nothing to assert on just make sure the function does not fail
	}
}
