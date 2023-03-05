package repositories

import (
	"errors"
)

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrUnableParseUser  = errors.New("unable parse user")
	ErrUnableDecodeURL  = errors.New("unable decode URL")
	ErrLinkNotExists    = errors.New("link not exists")
	ErrUserNotMatch     = errors.New("user not match")
)
