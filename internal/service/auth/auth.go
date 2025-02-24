package auth

import "errors"

var ErrUserNotFound = errors.New("user with this creds not found")
var ErrRestricted = errors.New("this user cannot access this app")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")

type UserAuth struct {
	TokenTTL int
	Secret   string
}

func NewUserAuth(tokenTTL int, secret string) *UserAuth {
	return &UserAuth{
		TokenTTL: tokenTTL,
		Secret:   secret,
	}
}

func (ua *UserAuth) Validate(token string) (bool, error) {
	return true, nil
}

func (ua *UserAuth) Login(username, password string, appId int) (string, string, error) {
	return "", "", nil
}

func (ua *UserAuth) Refresh(refreshToken string) (string, error) {
	return "", nil
}
