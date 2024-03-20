package apperrors

import "errors"

// account errors
var (
	ErrorAccountAlreadyExists        = errors.New("account already exists")
	ErrorAccountNotFound             = errors.New("account not found")
	ErrorAccountPasswordNotGenerated = errors.New("password hash generation error")
	ErrorAccountWrongPassword        = errors.New("wrong password")
	ErrorValidate                    = errors.New("some fields are incorrect")
	ErrorContextAccountIdNotFount    = errors.New("account id in context not found")
)

// auth errors
var (
	ErrorLoginOrPasswordIncorrect = errors.New("wrong login or password")
)

// session errors
var (
	ErrorSessionNotCreated      = errors.New("error occurred while creating session")
	ErrorSessionDeviceMismatch  = errors.New("device doesn't match with device of current session")
	ErrorContextSessionNotFound = errors.New("session id in context not found")
)

// jwt errors
var (
	ErrNoSigningKey         = errors.New("empty signing key")
	ErrNoClaims             = errors.New("error getting claims from token")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
)
