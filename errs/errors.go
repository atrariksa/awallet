package errs

import "errors"

var (
	ErrUserAlreadyExists error = errors.New("User Already Exist")
	ErrUnauthorized      error = errors.New("Unauthorized")
	ErrInternalServer    error = errors.New("Internal Server Error")
)
