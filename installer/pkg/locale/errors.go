package locale

import "fmt"

type LocaleGenError struct {
	Err error
}

func (e LocaleGenError) Error() string {
	return fmt.Sprintf("Error while configuring locales: error=%s", e.Err.Error())
}

func (e LocaleGenError) Unwrap() error {
	return e.Err
}
