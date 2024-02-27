package apperrors

import "errors"

var (
	ErrorAccountAlreadyExists        = errors.New("account already exists")
	ErrAccountNotFound               = errors.New("account not found")
	ErrorAccountPasswordNotGenerated = errors.New("password hash generation error")
	ErrorAccountWrongPassword        = errors.New("wrong password")
)
