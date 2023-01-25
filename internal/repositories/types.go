package repositories

import (
	"github.com/google/uuid"
)

type (
	ID   = string
	URL  = string
	User = uuid.UUID
)

type LinkData struct {
	URL     URL
	User    User
	Deleted bool
}

type UserLink struct {
	ID  ID
	URL URL
}

type LinkPendingDeletion struct {
	ID   ID
	User User
}
