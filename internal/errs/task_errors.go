package errs

import "errors"

var ErrTaskNotFound = errors.New("not found")
var ErrInvalidTask = errors.New("invalid task")
var ErrNotValidFields = errors.New("not valid fields")
var ErrBadConvertation = errors.New("problem to convertation")
