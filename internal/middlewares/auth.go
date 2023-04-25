package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

// UserCookie - middleware для аутентификации пользователя.
//
// Если пользователь обращается первый раз, то генерируем userID и передаем его в cookie.
// Если у пользователя уже есть ID, то проверяем подпись.
func (m *Middlewares) UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER")
		if errors.Is(err, http.ErrNoCookie) || len(cookie.Value) < 16 {
			m.setUser(next, w, r, m.createNewUser(w))
			return
		}

		user, err := m.a.Load(cookie.Value)
		if errors.Is(err, authenticator.ErrUnauthorized) {
			user = m.createNewUser(w)
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		m.setUser(next, w, r, user)
	})
}

// setUser добавляет userID в контекст и передает запрос следующему обработчику
func (m *Middlewares) setUser(next http.Handler, w http.ResponseWriter, r *http.Request, user uuid.UUID) {
	ctx := context.WithValue(r.Context(), utils.ContextKey("userID"), user)
	next.ServeHTTP(w, r.WithContext(ctx))
}

// createNewUser - генерирует пользователя, подписывает cookie и передает их клиенту.
func (m *Middlewares) createNewUser(w http.ResponseWriter) uuid.UUID {
	user, signed := m.a.Gen()

	cookie := http.Cookie{
		Name:    "USER",
		Value:   signed,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	return user
}
