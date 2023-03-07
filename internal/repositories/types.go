package repositories

import (
	"github.com/google/uuid"
)

// Типы, которые используем для хранения информации о ссылках.
type (
	ID   = string    // Тип для хранения ID сокращенной ссылки.
	URL  = string    // Тип для хранения исходного URL.
	User = uuid.UUID // Тип для хранения ID пользователя.
)

// LinkData - структура для хранения данных о ссылке по ключу ID.
type LinkData struct {
	URL     URL  // Исходный URL.
	User    User // Пользователь, которому принадлежит ссылка.
	Deleted bool // Удалена ли ссылка.
}

// UserLink - структура для хранения ссылки пользователя.
type UserLink struct {
	ID  ID  // ID сокращенной ссылки.
	URL URL // Исходный URL.
}

// LinkPendingDeletion - структура с информацией о ссылке, ожидающей удаления.
type LinkPendingDeletion struct {
	ID   ID   // ID сокращенной ссылки.
	User User // Пользователь, которому принадлежит ссылка.
}
