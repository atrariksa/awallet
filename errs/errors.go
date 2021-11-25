package errs

import "errors"

var (
	ErrUserAlreadyExists error = errors.New("User Already Exist")
	ErrInternalServer    error = errors.New("Internal Server Error")
)
