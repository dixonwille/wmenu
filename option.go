package wmenu

//Option are the options that are used inside the menus
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
