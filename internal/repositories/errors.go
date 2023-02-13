package repositories

import (
	"errors"
)

var (
	ErrURLNotFound      error = errors.New("URL not found")
	ErrURLAlreadyExists error = errors.New("URL already exists")
	ErrUnableParseUser  error = errors.New("unable parse user")
	ErrUnableDecodeURL  error = errors.New("unable decode URL")
	ErrLinkNotExists    error = errors.New("link not exists")
	ErrUserNotMatch     error = errors.New("user not match")
)
