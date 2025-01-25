package domain

import "github.com/waliqueiroz/mystery-gifter-api/pkg/validator"

type Credentials struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

func NewCredentials(email, password string) (*Credentials, error) {
	credentials := Credentials{
		Email:    email,
		Password: password,
	}

	if err := credentials.Validate(); err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (c *Credentials) Validate() error {
	if errs := validator.Validate(c); len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}

type AuthSession struct {
	User        User   `validate:"required"`
	AccessToken string `validate:"required"`
	TokenType   string `validate:"required"`
	ExpiresIn   int64  `validate:"required"`
}

func NewAuthSession(user User, accessToken, tokenType string, expiresIn int64) (*AuthSession, error) {
	authSession := AuthSession{
		User:        user,
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiresIn,
	}

	if err := authSession.Validate(); err != nil {
		return nil, err
	}

	return &authSession, nil
}

func (a *AuthSession) Validate() error {
	if errs := validator.Validate(a); len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}
