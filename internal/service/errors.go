package service

import "errors"

var ErrUserNotFound = errors.New("user with this creds not found")
var ErrRestricted = errors.New("this user cannot access this app")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token has expired")
