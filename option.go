package wmenu

//Option are the options that are used inside the menus
type Option struct {
	id        int
	text      string
	function  func()
	isDefault bool
}

func newOption(id int, text string, def bool, function func()) *Option {
	return &Option{
		id:        id,
		text:      text,
		isDefault: def,
		function:  function,
	}
}
