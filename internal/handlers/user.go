package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

// ErrValueIsNotUUID - значение не может быть преобразовано к типу uuid.UUID.
var ErrValueIsNotUUID = errors.New("value is not uuid.UUID")

func getUser(r *http.Request) (user uuid.UUID, err error) {
	user, ok := r.Context().Value(utils.ContextKey("userID")).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrValueIsNotUUID
	}
	return
}
