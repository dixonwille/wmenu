package wmenu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dixonwille/wlog"
)

//Menu is the structure for all menus created
type Menu struct {
	question        string
	defaultFunction func(Option)
	options         []Option
	ui              wlog.UI
	multiSeperator  string
	multiFunction   func([]Option)
	loopOnInvalid   bool
	clear           bool
}

//NewMenu creates a menu to use
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
		multiSeperator:  " ",
		multiFunction:   nil,
		loopOnInvalid:   false,
		clear:           false,
	}
}

//AddColor will change the color of the menu items.
//optionColor changes the color of the options.
//questionColor changes the color of the questions.
//errorColor changes the color of the question.
//Use wlog.None if you do not want to change the color.
func (m *Menu) AddColor(optionColor, questionColor, errorColor wlog.Color) {
	m.ui = wlog.AddColor(wlog.None, optionColor, wlog.None, questionColor, errorColor, wlog.None, wlog.None, m.ui)
}

//ClearOnMenuRun will clear the screen when a menu is ran.
func (m *Menu) ClearOnMenuRun() {
	m.clear = true
}

//SetSeperator sets the seperator to use for multi select.
//Default value is a space.
func (m *Menu) SetSeperator(sep string) {
	m.multiSeperator = sep
}

//LoopOnInvalid is used if an invalid option was given then clear the screen and ask again
func (m *Menu) LoopOnInvalid() {
	m.loopOnInvalid = true
}

//Option adds options to the menu
func (m *Menu) Option(title string, isDefault bool, function func()) {
	option := newOption(len(m.options), title, isDefault, function)
	m.options = append(m.options, *option)
}

//Action adds a default action if an option does not have an action or no option set as default.
func (m *Menu) Action(function func(Option)) {
	m.defaultFunction = function
}

//MultipleAction is called when multiple options are selected.
func (m *Menu) MultipleAction(function func([]Option)) {
	m.multiFunction = function
}

//Run is used to execute the menu (print to screen and ask quesiton)
func (m *Menu) Run() error {
	if m.clear {
		Clear()
	}
	valid := false
	var options []Option
	//Loop and on error check if loopOnInvalid is enabled.
	//If it is Clear the screen and write error.
	//Then ask again
	for !valid {
		//step 1 print things to screen
		m.print()
		//step 2 get and validate response
		opt, err := m.ask()
		if err != nil {
			if m.loopOnInvalid {
				if m.clear {
					Clear()
				}
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
	switch len(options) {
	//if no options go through options and look for default options
	case 0:
		opt := m.getDefault()
		switch len(opt) {
		//if there are no default options call the defaultFunction of the menu
		case 0:
			m.defaultFunction(Option{id: -1})
			//if there is one default option call it's function if it exist
			//if it does not, call the menu's defaultFunction
		case 1:
			if opt[0].function == nil {
				m.defaultFunction(opt[0])
			} else {
				opt[0].function()
			}
			//if there is more than one default option call the menu's multiFunction
		default:
			m.multiFunction(opt)
		}
		//if there is one option call it's funciton if it exist
		//if it does not, call the menu's defaultFunction
	case 1:
		if options[0].function == nil {
			m.defaultFunction(options[0])
		} else {
			options[0].function()
		}
		//if there is more than one option call the menu's multiFunction
	default:
		m.multiFunction(options)
	}
	return nil
}

func (m *Menu) print() {
	for _, opt := range m.options {
		m.ui.Output(fmt.Sprintf("%d) %s", opt.id, opt.text))
	}
	m.ui.Info(m.question)
}

func (m *Menu) ask() ([]Option, error) {
	reader := bufio.NewReader(os.Stdin)
	res, _ := reader.ReadString('\n')
	res = strings.Replace(res, "\r", "", -1) //this will only be useful under windows
	res = strings.Replace(res, "\n", "", -1)

	//Validate responses
	//Check if no responses are returned and no action to call
	if res == "" {
		//get default options
		opt := m.getDefault()
		if m.checkOptAndFunc(opt) {
			return nil, newMenuError(ErrNoResponse, "")
		}
		return nil, nil
	}

	resStrings := strings.Split(res, m.multiSeperator) //split responses by spaces
	//Check if we don't want multiple responses
	if m.multiFunction == nil && len(resStrings) > 1 {
		return nil, newMenuError(ErrTooMany, "")
	}

	//Convert responses to intigers
	var responses []int
	for _, response := range resStrings {
		//Check if it is an intiger
		r, err := strconv.Atoi(response)
		if err != nil {
			return nil, newMenuError(ErrInvalid, response)
		}
		responses = append(responses, r)
	}

	//Check if response is in the range of options
	//If it is make sure it is not duplicated
	var tmp []int
	for _, response := range responses {
		if response < 0 || len(m.options)-1 < response {
			return nil, newMenuError(ErrInvalid, strconv.Itoa(response))
		}

		if exist(tmp, response) {
			return nil, newMenuError(ErrDuplicate, strconv.Itoa(response))
		}

		tmp = append(tmp, response)
	}

	//Parse responses and return them as options
	var finalOptions []Option
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

func (m *Menu) getDefault() []Option {
	var opt []Option
	for _, o := range m.options {
		if o.isDefault {
			opt = append(opt, o)
		}
	}
	return opt
}

//make sure that there is an action available to be called in certain cases
func (m *Menu) checkOptAndFunc(opt []Option) bool {
	return ((len(opt) == 0 && m.defaultFunction == nil) || (len(opt) == 1 && opt[0].function == nil && m.defaultFunction == nil) || (len(opt) > 0 && m.multiFunction == nil))
}
