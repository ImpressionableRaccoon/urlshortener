package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

var (
	ErrValueIsNotString = errors.New("value is not string")
)

func getUser(r *http.Request) (user uuid.UUID, err error) {
	userID, ok := r.Context().Value(utils.ContextKey("userID")).(string)
	if !ok {
		return uuid.Nil, ErrValueIsNotString
	}
	return uuid.Parse(userID)
}
