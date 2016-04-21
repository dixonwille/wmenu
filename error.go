package wmenu

import "errors"

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

//MenuError records menu errors
type MenuError struct {
	Err error
	Res string
}

//Error prints the error in an easy to read string.
func (e *MenuError) Error() string {
	if e.Res != "" {
		return e.Err.Error() + ": " + e.Res
	}
	return e.Err.Error()
}

func newMenuError(err error, res string) *MenuError {
	return &MenuError{
		Err: err,
		Res: res,
	}
}
