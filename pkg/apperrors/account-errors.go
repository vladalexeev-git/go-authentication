package apperrors

import "errors"

// account errors
var (
	ErrorAccountAlreadyExists        = errors.New("account already exists")
	ErrorAccountNotFound             = errors.New("account not found")
	ErrorAccountPasswordNotGenerated = errors.New("password hash generation error")
	ErrorAccountWrongPassword        = errors.New("wrong password")
	ErrorValidate                    = errors.New("some fields are incorrect")
)

// session errors
var (
	ErrorSessionNotCreated = errors.New("error occurred while creating session")
)

// jwt errors
var (
	ErrNoSigningKey         = errors.New("empty signing key")
	ErrNoClaims             = errors.New("error getting claims from token")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
)
