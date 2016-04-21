package wmenu

//Option is what Menu uses to display options to screen.
//Also holds information on what should run and if it is a default option
type Option struct {
	ID        int
	Text      string
	function  func()
	isDefault bool
}

func newOption(id int, text string, def bool, function func()) *Option {
	return &Option{
		ID:        id,
		Text:      text,
		isDefault: def,
		function:  function,
	}
}
