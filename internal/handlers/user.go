package handlers

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

var (
	ErrValueIsNotString error = errors.New("value is not string")
)

func getUser(r *http.Request) (user uuid.UUID, err error) {
	value := r.Context().Value(utils.ContextKey("userID"))
	if reflect.ValueOf(value).Kind() != reflect.String {
		return uuid.Nil, ErrValueIsNotString
	}

	return uuid.Parse(r.Context().Value(utils.ContextKey("userID")).(string))
}
