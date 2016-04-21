# WMenu [![Build Status](https://travis-ci.org/dixonwille/wmenu.svg?branch=master)](https://travis-ci.org/dixonwille/wmenu) [![Go Report Card](https://goreportcard.com/badge/github.com/dixonwille/wmenu)](https://goreportcard.com/report/github.com/dixonwille/wmenu) [![GoDoc](https://godoc.org/github.com/dixonwille/wmenu?status.svg)](https://godoc.org/github.com/dixonwille/wmenu)

Package wmenu creates menus for cli programs. It uses wlog for it's interface
with the command line. It uses os.Stdin, os.Stdout, and os.Stderr with
concurrency by default. wmenu allows you to change the color of the different
parts of the menu. This package also creates it's own error structure so you can
type assert if you need to.

## Import
    import "github.com/dixonwille/wmenu"

## Examples

### Error Duplicate

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
We caught the err: Duplicated response: 1
```

### Error Invalid

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
We caught the err: Invalid response: 3
```

### Error No Response

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
We caught the err: No response
```

### Error Too Many

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
We caught the err: Too many responses
```

### Multiple

Code:

``` go
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
menu.SetSeperator(",")
menu.Option("Option 0", true, optFunc)
menu.Option("Option 1", false, nil)
menu.Option("Option 2", true, nil)
err := menu.Run()
if err != nil {
    log.Fatal(err)
}
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
Option 1 has an id of 1.
Option 2 has an id of 2.
```

### Multiple With Defaults

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
Option 0 has an id of 0.
Option 2 has an id of 2.
```

### Simple

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
Option 1 has an id of 1.
```

### Simple With Default

Code:

``` go
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
```

Output:

```
0) Option 0
1) Option 1
2) Option 2
Choose an option.
Option 0 was chosen.
```

## Usage

```go
var (
	//ErrInvalid is returned if a response from user was an invalid option
	ErrInvalid = errors.New("Invalid response")

	//ErrTooMany is returned if multiSelect is false and a user tries to select multiple options
	ErrTooMany = errors.New("Too many responses")

	//ErrNoResponse is returned if there were no responses and no action to call
	ErrNoResponse = errors.New("No response")

	//ErrDuplicate is returned is a user selects an option twice
	ErrDuplicate = errors.New("Duplicated response")
)
```

#### func  Clear

```go
func Clear()
```
Clear simply clears the command line interface.

#### func  IsDuplicateErr

```go
func IsDuplicateErr(err error) bool
```
IsDuplicateErr checks to see if err is of type duplicate returned by menu.

#### func  IsInvalidErr

```go
func IsInvalidErr(err error) bool
```
IsInvalidErr checks to see if err is of type invalid error returned by menu.

#### func  IsNoResponseErr

```go
func IsNoResponseErr(err error) bool
```
IsNoResponseErr checks to see if err is of type no response returned by menu.

#### func  IsTooManyErr

```go
func IsTooManyErr(err error) bool
```
IsTooManyErr checks to see if err is of type too many returned by menu.

#### type Menu

```go
type Menu struct {
	// contains filtered or unexported fields
}
```

Menu is used to display options to a user. A user can then select options and
Menu will validate the response and perform the correct action.

#### func  NewMenu

```go
func NewMenu(question string) *Menu
```
NewMenu creates a menu with a wlog.UI as the writer.

#### func (*Menu) Action

```go
func (m *Menu) Action(function func(Option))
```
Action adds a default action to use in certain scenarios. If the selected option
(by default or user selected) does not have a function applied to it this will
be called. If there are no default options and no option was selected this will
be called with an option that has an ID of -1.

#### func (*Menu) AddColor

```go
func (m *Menu) AddColor(optionColor, questionColor, responseColor, errorColor wlog.Color)
```
AddColor will change the color of the menu items. optionColor changes the color
of the options. questionColor changes the color of the questions. errorColor
changes the color of the question. Use wlog.None if you do not want to change
the color.

#### func (*Menu) ChangeReaderWriter

```go
func (m *Menu) ChangeReaderWriter(reader io.Reader, writer, errorWriter io.Writer)
```
ChangeReaderWriter changes where the menu listens and writes to. reader is where
user input is collected. writer and errorWriter is where the menu should write
to.

#### func (*Menu) ClearOnMenuRun

```go
func (m *Menu) ClearOnMenuRun()
```
ClearOnMenuRun will clear the screen when a menu is ran. This is checked when
LoopOnInvalid is activated. Meaning if an error occurred then it will clear the
screen before asking again.

#### func (*Menu) LoopOnInvalid

```go
func (m *Menu) LoopOnInvalid()
```
LoopOnInvalid is used if an invalid option was given then it will prompt the
user again.

#### func (*Menu) MultipleAction

```go
func (m *Menu) MultipleAction(function func([]Option))
```
MultipleAction is called when multiple options are selected (by default or user
selected). If this is set then it uses the separator string specified by
SetSeparator (Default is a space) to separate the responses. If this is not set
then it is implied that the menu only allows for one option to be selected.

#### func (*Menu) Option

```go
func (m *Menu) Option(title string, isDefault bool, function func())
```
Option adds an option to the menu for the user to select from. title is the
string the user will select isDefault is whether this option is a default option
(IE when no options are selected). function is what is called when only this
option is selected. If function is nil then it will default to the menu's
Action.

#### func (*Menu) Run

```go
func (m *Menu) Run() error
```
Run is used to execute the menu. It will print to options and question to the
screen. It will only clear the screen if ClearOnMenuRun is activated. This will
validate all responses. Errors are of type MenuError.

#### func (*Menu) SetSeparator

```go
func (m *Menu) SetSeparator(sep string)
```
SetSeparator sets the separator to use when multiple options are valid
responses. Default value is a space.

#### type MenuError

```go
type MenuError struct {
	Err error
	Res string
}
```

MenuError records menu errors

#### func (*MenuError) Error

```go
func (e *MenuError) Error() string
```
Error prints the error in an easy to read string.

#### type Option

```go
type Option struct {
	ID   int
	Text string
	// contains filtered or unexported fields
}
```

Option is what Menu uses to display options to screen. Also holds information on
what should run and if it is a default option
