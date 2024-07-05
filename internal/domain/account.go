package domain

import (
	"fmt"
	"go-authentication/pkg/apperrors"
	"go-authentication/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (a *Account) GenPasswordHash() error {

	// the argument 11 represents the "cost" factor,
	//indicating the number of iterations performed during password hashing.
	//Typically, recommended values for cost range from 10 to 14,
	//depending on security and performance requirements
	b, err := bcrypt.GenerateFromPassword([]byte(a.Password), 11)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", apperrors.ErrorAccountPasswordNotGenerated)
	}

	a.PasswordHash = string(b)
	return nil
}

func (a *Account) CompareHashAndPassword() error {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(a.Password))
	if err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", apperrors.ErrorAccountWrongPassword)
	}
	return nil
}

func (a *Account) RandomPassword() {
	a.Password = utils.RandomSpecialString(16)
}
