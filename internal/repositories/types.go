package repositories

import (
	"github.com/google/uuid"
)

// Типы, которые используем для хранения информации о ссылках.
type (
	ID      = string    // Тип для хранения ID сокращенной ссылки.
	URL     = string    // Тип для хранения исходного URL.
	User    = uuid.UUID // Тип для хранения ID пользователя.
	Deleted = bool      // Тип для хранения, удалена ли ссылка.
)

// LinkData - структура для хранения данных о ссылке.
type LinkData struct {
	ID      ID      // ID сокращенной ссылки.
	URL     URL     // Исходный URL.
	User    User    // Пользователь, которому принадлежит ссылка.
	Deleted Deleted // Удалена ли ссылка.
}

// ServiceStats - структура для хранения статистики сервиса.
type ServiceStats struct {
	URLs  uint64 `json:"urls"`  // Количество сокращённых URL в сервисе.
	Users uint64 `json:"users"` // Количество пользователей в сервисе.
}
