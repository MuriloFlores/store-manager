package vo

import (
	"errors"
	"net/mail"
)

var (
	ErrEmptyEmail   = errors.New("empty email")
	ErrInvalidEmail = errors.New("invalid email")
)

type Email string

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", ErrEmptyEmail
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", ErrInvalidEmail
	}

	return Email(email), nil
}

func (e Email) Equals(o Email) bool {
	return e.String() == o.String()
}

func (e Email) String() string {
	return string(e)
}
