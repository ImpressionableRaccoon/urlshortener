package repositories

import "github.com/google/uuid"

type (
	ID   = string
	URL  = string
	User = uuid.UUID
)

type LinkData struct {
	URL  URL
	User User
}

type UserLink struct {
	ID  ID
	URL URL
}
