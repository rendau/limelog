package errs

import (
	"github.com/rendau/dop/dopErrs"
)

const (
	WrongPassword = dopErrs.Err("wrong_password")
	BadDuration   = dopErrs.Err("bad_duration")
)
