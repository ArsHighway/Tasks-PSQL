package errs

import "errors"

var ErrUserNotFound = errors.New("not found")
var ErrInvalidUser = errors.New("invalid user")
var ErrNotValidFieldsUser = errors.New("not valid fields")
var ErrUserBadConvertation = errors.New("problem to convertation")
