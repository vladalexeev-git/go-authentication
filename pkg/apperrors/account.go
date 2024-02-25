package apperrors

import "errors"

var (
	ErrorAccountAlreadyExists = errors.New("account.go already exists")
	ErrAccountNotFound        = errors.New("account.go not found")
)
