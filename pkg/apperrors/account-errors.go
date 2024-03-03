package apperrors

import "errors"

var (
	ErrorAccountAlreadyExists        = errors.New("account already exists")
	ErrorAccountNotFound             = errors.New("account not found")
	ErrorAccountPasswordNotGenerated = errors.New("password hash generation error")
	ErrorAccountWrongPassword        = errors.New("wrong password")
	ErrorValidate                    = errors.New("some fields are incorrect")
)
