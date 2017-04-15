package wmenu

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wlog "gopkg.in/dixonwille/wlog.v2"
)

var newMenuCases = []string{"Testing this menu.", "", "!@#$%^&*()"}
var setSeperatorCases = []string{"", ".", ",", "~"}
var defaultIconCases = []string{"", ",", "*", "~", "!", "@"}
var setTriesCases = []int{0, -4, 5}
var optionCases = []struct {
	name     string
	value    interface{}
	def      bool
	function func(Opt) error
}{
	{"Options", "Options", true, func(Opt) error { fmt.Println("testing option"); return nil }},
	{"", "Options", true, func(Opt) error { fmt.Println("testing option"); return nil }},
	{"Options", "", false, func(Opt) error { fmt.Println("testing option"); return nil }},
	{"", "", false, nil},
	{"", nil, false, nil},
}
var actionCases = []func([]Opt) error{
	func(opts []Opt) error { fmt.Println(opts); return nil },
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
var errorCases = []struct {
	input       string
	optFunction func(option Opt) error
	defFunction func([]Opt) error
	expected    string
	singDef     bool
	multiDef    bool
}{
	{"0\r\n", func(option Opt) error { return errors.New("Oops") }, nil, "Oops", false, false},
	{"0\r\n", nil, func(opts []Opt) error { return errors.New("Oops") }, "Oops", false, false},
	{"0 1\r\n", nil, func(opts []Opt) error { return errors.New("Oops") }, "Too many responses", false, false},
	{"\r\n", func(option Opt) error { return errors.New("Oops") }, nil, "Oops", true, false},
	{"\r\n", nil, func(opts []Opt) error { return errors.New("Oops") }, "Oops", true, false},
	{"\r\n", nil, func(opts []Opt) error { return errors.New("Oops") }, "Oops", false, false},
	{"\r\n", nil, func(opts []Opt) error { return errors.New("Oops") }, "Oops", true, true},
}
var ynCases = []struct {
	input    string
	defFunc  func([]Opt) error
	expected string
	def      int
}{
	{"\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "y", 0},
	{"\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "n", 1},
	{"YES\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "y", 0},
	{"n\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "n", 0},
	{" Yes \r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "y", 0},
	{"ahh\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "Invalid response: ahh", 0},
	{"boo ahh\r\n", func(opts []Opt) error { return errors.New(opts[0].Text) }, "Too many responses", 0},
}
var trimCases = []struct {
	input    string
	del      string
	expected string
}{
	{"0 \r\n", " ", "0\r\n"},
	{" 0\r\n", " ", "0\r\n"},
	{" 0 \r\n", " ", "0\r\n"},
	{"0 1 \r\n", " ", "0 1\r\n"},
	{" 0 1\r\n", " ", "0 1\r\n"},
	{" 0 1 \r\n", " ", "0 1\r\n"},
	{" 0, 1 \r\n", ",", "0 1\r\n"},
	{"0, 1,\r\n", ",", "0 1\r\n"},
	{"0, 1, \r\n", ",", "0 1\r\n"},
}

func Example_optionValue() {
	type NameEntity struct {
		FirstName string
		LastName  string
	}

	reader := strings.NewReader("1\r\n") //Simulates the user typing "1" and hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	actFunc := func(opts []Opt) error {
		name, ok := opts[0].Value.(NameEntity)
		if !ok {
			log.Fatal("Could not cast option's value to NameEntity")
		}
		fmt.Printf("%s has an id of %d.\n", opts[0].Text, opts[0].ID)
		fmt.Printf("Hello, %s %s.\n", name.FirstName, name.LastName)
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.Option("Option 0", NameEntity{"Bill", "Bob"}, true, optFunc)
	menu.Option("Option 1", NameEntity{"John", "Doe"}, false, nil)
	menu.Option("Option 2", NameEntity{"Jane", "Doe"}, false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) *Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 1 has an id of 1.
	//Hello, John Doe.
}

func Example_yesNo() {
	reader := strings.NewReader("y\r\n") //Simulates the user typing "y" and hitting the [enter] key
	actFunc := func(opts []Opt) error {
		fmt.Printf("%s has an id of %d.\n", opts[0].Text, opts[0].ID)
		fmt.Printf("But has a value of %s.\n", opts[0].Value.(string))
		return nil
	}
	menu := NewMenu("Would you like to start?")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.IsYesNo(0)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//Would you like to start? (Y/n)
	//y has an id of 0.
	//But has a value of yes.
}

func Example_simpleDefault() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Fprint(os.Stdout, "Option 0 was chosen.")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", "0", true, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) *Option 0
	//1) Option 1
	//2) Option 2
	//Choose an option.
	//Option 0 was chosen.
}

func Example_multiple() {
	reader := strings.NewReader("1,2\r\n") //Simulates the user typing "1,2" and hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	actFunc := func(opts []Opt) error {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.AllowMultiple()
	menu.SetSeparator(",")
	menu.Option("Option 0", "0", true, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", true, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) *Option 0
	//1) Option 1
	//2) *Option 2
	//Choose an option.
	//Option 1 has an id of 1.
	//Option 2 has an id of 2.
}

func Example_multipleDefault() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	actFunc := func(opts []Opt) error {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.AllowMultiple()
	menu.Option("Option 0", "0", true, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", true, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
	//Output:
	//0) *Option 0
	//1) Option 1
	//2) *Option 2
	//Choose an option.
	//Option 0 has an id of 0.
	//Option 2 has an id of 2.
}

func Example_errorNoResponse() {
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
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
	reader := strings.NewReader("3\r\n") //Simulates the user typing "3" and hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
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
	reader := strings.NewReader("1 2\r\n") //Simulates the user typing "1 2" and hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
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
	reader := strings.NewReader("1 1\r\n") //Simulates the user typing "1 1" and hitting the [enter] key
	optFunc := func(option Opt) error {
		fmt.Println("Option 0 was chosen.")
		return nil
	}
	actFunc := func(opts []Opt) error {
		for _, opt := range opts {
			fmt.Printf("%s has an id of %d.\n", opt.Text, opt.ID)
		}
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, os.Stdout, os.Stderr)
	menu.Action(actFunc)
	menu.AllowMultiple()
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
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
		assert.Nil(menu.function)
		assert.Nil(menu.options)
		assert.Equal(" ", menu.multiSeparator)
		assert.False(menu.allowMultiple)
		assert.False(menu.loopOnInvalid)
		assert.False(menu.clear)
		assert.NotNil(menu.ui)
		assert.Equal(3, menu.tries)
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

func TestSetTries(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range setTriesCases {
		menu.SetTries(c)
		assert.Equal(t, c, menu.tries)
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
		menu.Option(c.name, c.value, c.def, c.function)
		require.Equal(t, i+1, len(menu.options))
		assert.Equal(i, menu.options[i].ID)
		assert.Equal(c.name, menu.options[i].Text)
		assert.Equal(c.def, menu.options[i].isDefault)
		assert.Equal(c.value, menu.options[i].Value)
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
			assert.NotNil(t, menu.function)
		} else {
			assert.Nil(t, menu.function)
		}
	}
}

func TestAllowMultiple(t *testing.T) {
	menu := NewMenu("Testing")
	menu.AllowMultiple()
	assert.True(t, menu.allowMultiple)
}

func TestAddColor(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range addColorCases {
		menu.AddColor(c.opt, c.que, c.res, c.err)
		//Nothing to assert on just make sure the function does not fail
	}
}

func TestSetDefaultIcon(t *testing.T) {
	menu := NewMenu("Testing")
	for _, c := range defaultIconCases {
		menu.SetDefaultIcon(c)
		assert.Equal(t, c, menu.defIcon)
	}

}

func TestClearInAsk(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("1\r\n") //Simulates the user typing "1" and hitting the [enter] key
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called Option 0's function")
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Len(t, opts, 1)
		assert.Equal(t, 1, opts[0].ID)
		assert.Equal(t, "Option 1", opts[0].Text)
		assert.Nil(t, opts[0].function)
		assert.False(t, opts[0].isDefault)
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.ClearOnMenuRun()
	menu.Option("Option 0", "0", true, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", true, nil)
	err := menu.Run()
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestDefaultAction(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called option 0's function")
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Len(t, opts, 1)
		assert.Equal(t, -1, opts[0].ID)
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", false, nil)
	menu.Option("Option 2", "2", false, nil)
	// spew.Dump(menu)
	err := menu.Run()
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestDefaultActionWithDefaultOption(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("\r\n") //Simulates the user hitting the [enter] key
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called option 0's function")
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Len(t, opts, 1)
		assert.Equal(t, 1, opts[0].ID)
		assert.Equal(t, "Option 1", opts[0].Text)
		assert.Nil(t, opts[0].function)
		assert.True(t, opts[0].isDefault)
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", true, nil)
	menu.Option("Option 2", "2", false, nil)
	err := menu.Run()
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestOptionsFunction(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("0\r\n") //Simulates the user typing "0" and hitting the [enter] key
	optFunc := func(option Opt) error {
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Fail(t, "Should not have called the menu's default function")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", true, nil)
	menu.Option("Option 2", "2", false, nil)
	err := menu.Run()
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestWlogAskErr(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("1") //Simulates the user typing "1" without hitting [enter]. Can't happen when reader is os.Stdin
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called option 0's function")
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Fail(t, "Should not have called the menu's default function")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", true, nil)
	menu.Option("Option 2", "2", false, nil)
	err := menu.Run()
	if err != nil {
		assert.Equal(t, "EOF", err.Error())
	} else {
		assert.Fail(t, "Expected to get an error of EOF but instead got no error")
	}
}

func TestLetterForResponse(t *testing.T) {
	stdOut := initTest()
	reader := strings.NewReader("a\r\n") //Simulates the user typing "a" and hitting [enter].
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called option 0's function")
		return nil
	}
	actFunc := func(opts []Opt) error {
		assert.Fail(t, "Should not have called the menu's default function")
		return nil
	}
	menu := NewMenu("Choose an option.")
	menu.ChangeReaderWriter(reader, stdOut, stdOut)
	menu.Action(actFunc)
	menu.Option("Option 0", "0", false, optFunc)
	menu.Option("Option 1", "1", true, nil)
	menu.Option("Option 2", "2", false, nil)
	err := menu.Run()
	if err != nil {
		require.True(t, IsInvalidErr(err))
		e := err.(*MenuError)
		assert.Equal(t, "a", e.Res)
	} else {
		assert.Fail(t, "Expected to get an Invalid Response error but instead got no error")
	}
}

func TestActionError(t *testing.T) {
	stdout := initTest()
	for _, c := range errorCases {
		reader := strings.NewReader(c.input)
		menu := NewMenu("Choose an option.")
		menu.ChangeReaderWriter(reader, stdout, stdout)
		menu.Action(c.defFunction)
		if c.multiDef {
			menu.AllowMultiple()
		}
		menu.Option("Option 0", "0", c.singDef, c.optFunction)
		menu.Option("Option 1", "1", false, nil)
		menu.Option("Option 2", "2", c.multiDef, nil)
		err := menu.Run()
		if err != nil {
			assert.Equal(t, c.expected, err.Error())
		} else {
			assert.Fail(t, "Expected an error but did not get one")
		}
	}

}

func TestLoopAndTries(t *testing.T) {
	stdOut := initTest()
	optFunc := func(option Opt) error {
		assert.Fail(t, "Should not have called option 0's function")
		return nil
	}
	for _, c := range setTriesCases {
		reader := strings.NewReader("a") //Simulates the user typing "a" and not hitting [enter].
		menu := NewMenu("Choose an option.")
		menu.ChangeReaderWriter(reader, stdOut, stdOut)
		menu.SetTries(c)
		menu.LoopOnInvalid()
		menu.ClearOnMenuRun()
		menu.Option("Option 0", "0", false, optFunc)
		menu.Option("Option 1", "1", false, nil)
		menu.Option("Option 2", "2", false, nil)
		err := menu.Run()
		if err != nil {
			require.True(t, IsMenuErr(err))
			e := err.(*MenuError)
			assert.Equal(t, 0, e.TriesLeft)
		} else {
			assert.Fail(t, "Expected to get a general Menu error but instead got no error")
		}
	}
}

func TestYesNo(t *testing.T) {
	for _, c := range ynCases {
		stdOut := initTest()
		reader := strings.NewReader(c.input)
		menu := NewMenu("Yes or No")
		menu.ChangeReaderWriter(reader, stdOut, stdOut)
		menu.IsYesNo(c.def)
		menu.Action(c.defFunc)
		err := menu.Run()
		if err != nil {
			assert.Equal(t, c.expected, err.Error())
		}
	}
}

func TestTrim(t *testing.T) {
	for _, c := range trimCases {
		stdOut := initTest()
		actOut := initTest()
		reader := strings.NewReader(c.input)
		menu := NewMenu("Testing")
		menu.ChangeReaderWriter(reader, stdOut, stdOut)
		menu.SetSeparator(c.del)
		menu.Action(func(opts []Opt) error {
			var actual []string
			for _, opt := range opts {
				actual = append(actual, opt.Text)
			}
			fmt.Fprintf(actOut, "%s\r\n", strings.Join(actual, " "))
			return nil
		})
		menu.AllowMultiple()
		menu.Option("0", nil, false, nil)
		menu.Option("1", nil, false, nil)
		err := menu.Run()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		act, err := actOut.ReadString((byte)('\n'))
		if err != nil {
			assert.Fail(t, err.Error())
		}
		assert.Equal(t, c.expected, act)
	}
}

func initTest() *bytes.Buffer {
	var b []byte
	return bytes.NewBuffer(b)
}
