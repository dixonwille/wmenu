package wmenu

//TODO:0 Fix this so that helpers know which option I am talking about issue:1

//Opt is what Menu uses to display options to screen.
//Also holds information on what should run and if it is a default option
type Opt struct {
	ID        int
	Text      string
	function  func()
	isDefault bool
}

func newOption(id int, text string, def bool, function func()) *Opt {
	return &Opt{
		ID:        id,
		Text:      text,
		isDefault: def,
		function:  function,
	}
}
