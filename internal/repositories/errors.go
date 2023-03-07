package repositories

import (
	"errors"
)

var (
	ErrURLNotFound      = errors.New("URL not found")      // Ссылки с таким ID не существует.
	ErrURLAlreadyExists = errors.New("URL already exists") // Ссылка с таким исходным URL уже есть.
	ErrUnableParseUser  = errors.New("unable parse user")  // Не получается распарсить пользователя из файла.
	ErrUnableDecodeURL  = errors.New("unable decode URL")  // Не получается загрузить ссылку из файла.
	ErrLinkNotExists    = errors.New("link not exists")    // Ссылки с таким ID не существует.
	ErrUserNotMatch     = errors.New("user not match")     // Пользователь не может удалить чужую ссылку.
)
