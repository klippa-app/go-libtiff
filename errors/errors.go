package errors

import "fmt"

type TiffError struct {
	Module    string
	TiffError error
}

func (e TiffError) Error() string {
	return fmt.Sprintf("%s: %s", e.Module, e.TiffError.Error())
}

// Is allows use via errors.Is
func (e *TiffError) Is(err error) bool {
	if _, ok := err.(*TiffError); ok {
		return true
	}
	return false
}

func (e *TiffError) Unwrap() error {
	return e.TiffError
}
