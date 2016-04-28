//Package wmenu creates menus for cli programs.
//It uses wlog for it's interface with the command line.
//It uses os.Stdin, os.Stdout, and os.Stderr with concurrency by default.
//wmenu allows you to change the color of the different parts of the menu.
//This package also creates it's own error structure so you can type assert if you need to.
package wmenu

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dixonwille/wlog"
)

//TODO:0 Refactor the README to not include godoc as file issue:3
//TODO:0 Add default icon and method issue:4
//TODO:0 Add yes no option to the menu issue:5

//Menu is used to display options to a user.
//A user can then select options and Menu will validate the response and perform the correct action.
type Menu struct {
	question        string
	defaultFunction func(Opt) error
	options         []Opt
	ui              wlog.UI
	multiSeparator  string
	multiFunction   func([]Opt) error
	loopOnInvalid   bool
	clear           bool
	tries           int
	defIcon         string
}

//NewMenu creates a menu with a wlog.UI as the writer.
func NewMenu(question string) *Menu {
	//Create a default ui to use for menu
	var ui wlog.UI
	ui = wlog.New(os.Stdin, os.Stdout, os.Stderr)
	ui = wlog.AddConcurrent(ui)

	return &Menu{
		question:        question,
		defaultFunction: nil,
		options:         nil,
		ui:              ui,
		multiSeparator:  " ",
		multiFunction:   nil,
		loopOnInvalid:   false,
		clear:           false,
		tries:           3,
		defIcon:         "*",
	}
}

//AddColor will change the color of the menu items.
//optionColor changes the color of the options.
//questionColor changes the color of the questions.
//errorColor changes the color of the question.
//Use wlog.None if you do not want to change the color.
func (m *Menu) AddColor(optionColor, questionColor, responseColor, errorColor wlog.Color) {
	m.ui = wlog.AddColor(questionColor, errorColor, wlog.None, wlog.None, optionColor, responseColor, wlog.None, wlog.None, wlog.None, m.ui)
}

//ClearOnMenuRun will clear the screen when a menu is ran.
//This is checked when LoopOnInvalid is activated.
//Meaning if an error occurred then it will clear the screen before asking again.
func (m *Menu) ClearOnMenuRun() {
	m.clear = true
}

//SetSeparator sets the separator to use when multiple options are valid responses.
//Default value is a space.
func (m *Menu) SetSeparator(sep string) {
	m.multiSeparator = sep
}

//SetTries sets the number of tries on the loop before failing out.
//Default is 3.
//Negative values act like 0.
func (m *Menu) SetTries(i int) {
	m.tries = i
}

//LoopOnInvalid is used if an invalid option was given then it will prompt the user again.
func (m *Menu) LoopOnInvalid() {
	m.loopOnInvalid = true
}

//SetDefaultIcon sets the icon used to identify which options will be selected by default
func (m *Menu) SetDefaultIcon(icon string) {
	m.defIcon = icon
}

//Option adds an option to the menu for the user to select from.
//title is the string the user will select
//isDefault is whether this option is a default option (IE when no options are selected).
//function is what is called when only this option is selected.
//If function is nil then it will default to the menu's Action.
func (m *Menu) Option(title string, isDefault bool, function func() error) {
	option := newOption(len(m.options), title, isDefault, function)
	m.options = append(m.options, *option)
}

//Action adds a default action to use in certain scenarios.
//If the selected option (by default or user selected) does not have a function applied to it this will be called.
//If there are no default options and no option was selected this will be called with an option that has an ID of -1.
func (m *Menu) Action(function func(Opt) error) {
	m.defaultFunction = function
}

//MultipleAction is called when multiple options are selected (by default or user selected).
//If this is set then it uses the separator string specified by SetSeparator (Default is a space) to separate the responses.
//If this is not set then it is implied that the menu only allows for one option to be selected.
func (m *Menu) MultipleAction(function func([]Opt) error) {
	m.multiFunction = function
}

//ChangeReaderWriter changes where the menu listens and writes to.
//reader is where user input is collected.
//writer and errorWriter is where the menu should write to.
func (m *Menu) ChangeReaderWriter(reader io.Reader, writer, errorWriter io.Writer) {
	var ui wlog.UI
	ui = wlog.New(reader, writer, errorWriter)
	m.ui = ui
}

//Run is used to execute the menu.
//It will print to options and question to the screen.
//It will only clear the screen if ClearOnMenuRun is activated.
//This will validate all responses.
//Errors are of type MenuError.
func (m *Menu) Run() error {
	if m.clear {
		Clear()
	}
	valid := false
	var options []Opt
	//Loop and on error check if loopOnInvalid is enabled.
	//If it is Clear the screen and write error.
	//Then ask again
	for !valid {
		//step 1 print options to screen
		m.print()
		//step 2 ask question, get and validate response
		opt, err := m.ask()
		if err != nil {
			if !IsMenuErr(err) {
				err = newMenuError(err, "", m.triesLeft())
			}
			if m.loopOnInvalid && m.tries > 0 {
				if m.clear {
					Clear()
				}
				m.tries = m.tries - 1
				m.ui.Error(err.Error())
			} else {
				return err
			}
		} else {
			options = opt
			valid = true
		}
	}
	//step 3 call appropriate action with the responses
	return m.callAppropriate(options)
}

func (m *Menu) callAppropriate(options []Opt) (err error) {
	switch len(options) {
	//if no options go through options and look for default options
	case 0:
		return m.callAppropriateNoOptions()
		//if there is one option call it's funciton if it exist
		//if it does not, call the menu's defaultFunction
	case 1:
		if options[0].function == nil {
			if err := m.defaultFunction(options[0]); err != nil {
				return err
			}
		} else {
			if err := options[0].function(); err != nil {
				return err
			}
		}
		//if there is more than one option call the menu's multiFunction
	default:
		if err := m.multiFunction(options); err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) callAppropriateNoOptions() (err error) {
	opt := m.getDefault()
	switch len(opt) {
	//if there are no default options call the defaultFunction of the menu
	case 0:
		if err := m.defaultFunction(Opt{ID: -1}); err != nil {
			return err
		}
		//if there is one default option call it's function if it exist
		//if it does not, call the menu's defaultFunction
	case 1:
		if opt[0].function == nil {
			if err := m.defaultFunction(opt[0]); err != nil {
				return err
			}
		} else {
			if err := opt[0].function(); err != nil {
				return err
			}
		}
		//if there is more than one default option call the menu's multiFunction
	default:
		if err := m.multiFunction(opt); err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) print() {

	for _, opt := range m.options {
		icon := m.defIcon
		if !opt.isDefault {
			icon = ""
		}
		m.ui.Output(fmt.Sprintf("%d) %s%s", opt.ID, icon, opt.Text))
	}
}

func (m *Menu) ask() ([]Opt, error) {
	res, err := m.ui.Ask(m.question)
	if err != nil {
		return nil, err
	}
	//Validate responses
	//Check if no responses are returned and no action to call
	if res == "" {
		//get default options
		opt := m.getDefault()
		if m.checkOptAndFunc(opt) {
			return nil, newMenuError(ErrNoResponse, "", m.triesLeft())
		}
		return nil, nil
	}

	resStrings := strings.Split(res, m.multiSeparator) //split responses by spaces
	//Check if we don't want multiple responses
	if m.multiFunction == nil && len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "", m.triesLeft())
	}

	//Convert responses to intigers
	var responses []int
	for _, response := range resStrings {
		//Check if it is an intiger
		r, err := strconv.Atoi(response)
		if err != nil {
			return nil, newMenuError(ErrInvalid, response, m.triesLeft())
		}
		responses = append(responses, r)
	}

	//Check if response is in the range of options
	//If it is make sure it is not duplicated
	var tmp []int
	for _, response := range responses {
		if response < 0 || len(m.options)-1 < response {
			return nil, newMenuError(ErrInvalid, strconv.Itoa(response), m.triesLeft())
		}

		if exist(tmp, response) {
			return nil, newMenuError(ErrDuplicate, strconv.Itoa(response), m.triesLeft())
		}

		tmp = append(tmp, response)
	}

	//Parse responses and return them as options
	var finalOptions []Opt
	for _, response := range responses {
		finalOptions = append(finalOptions, m.options[response])
	}

	return finalOptions, nil
}

//Simply checks if number exists in the slice
func exist(slice []int, number int) bool {
	for _, s := range slice {
		if number == s {
			return true
		}
	}
	return false
}

//gets a list of default options
func (m *Menu) getDefault() []Opt {
	var opt []Opt
	for _, o := range m.options {
		if o.isDefault {
			opt = append(opt, o)
		}
	}
	return opt
}

//make sure that there is an action available to be called in certain cases
func (m *Menu) checkOptAndFunc(opt []Opt) bool {
	return ((len(opt) == 0 && m.defaultFunction == nil) || (len(opt) == 1 && opt[0].function == nil && m.defaultFunction == nil) || (len(opt) > 1 && m.multiFunction == nil))
}

func (m *Menu) triesLeft() int {
	if m.loopOnInvalid && m.tries > 0 {
		return m.tries
	}
	return 0
}
