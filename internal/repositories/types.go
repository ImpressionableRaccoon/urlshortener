package repositories

import "github.com/google/uuid"

type ID = string
type URL = string
type User = uuid.UUID

type LinkData struct {
	URL  URL
	User User
}

type UserLink struct {
	ID  ID
	URL URL
}
