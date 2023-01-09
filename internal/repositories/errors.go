package repositories

import "errors"

var (
	URLNotFound      error = errors.New("URL not found")
	URLAlreadyExists error = errors.New("URL already exists")
)
