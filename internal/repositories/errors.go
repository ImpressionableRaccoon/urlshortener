package repositories

import (
	"errors"
)

var (
	ErrURLNotFound      error = errors.New("URL not found")
	ErrURLAlreadyExists error = errors.New("URL already exists")
)
