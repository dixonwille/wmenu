package wmenu

//Opt is what Menu uses to display options to screen.
//Also holds information on what should run and if it is a default option
type Opt struct {
	ID        int
	Text      string
	Value     string
	function  func() error
	isDefault bool
}

func newOption(id int, text string, value string, def bool, function func() error) *Opt {
	return &Opt{
		ID:        id,
		Text:      text,
		Value:     value,
		isDefault: def,
		function:  function,
	}
}
