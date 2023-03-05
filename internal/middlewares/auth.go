package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

func (m *Middlewares) UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER")
		if errors.Is(err, http.ErrNoCookie) || len(cookie.Value) < 16 {
			m.setNewUser(next, w, r, m.createNewUser(w))
			return
		}

		payload, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil || len(payload) < 16 {
			m.setNewUser(next, w, r, m.createNewUser(w))
			return
		}

		h := hmac.New(sha256.New, m.cfg.CookieKey)
		h.Write(payload[:16])
		sign := h.Sum(nil)

		if !hmac.Equal(sign, payload[16:]) {
			m.setNewUser(next, w, r, m.createNewUser(w))
			return
		}

		user, err := uuid.FromBytes(payload[:16])
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		m.setNewUser(next, w, r, user.String())
	})
}

func (m *Middlewares) setNewUser(next http.Handler, w http.ResponseWriter, r *http.Request, user string) {
	ctx := context.WithValue(r.Context(), utils.ContextKey("userID"), user)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func (m *Middlewares) createNewUser(w http.ResponseWriter) string {
	user := uuid.New()

	b, _ := user.MarshalBinary()

	h := hmac.New(sha256.New, m.cfg.CookieKey)
	h.Write(b)
	sign := h.Sum(nil)

	cookie := http.Cookie{
		Name:    "USER",
		Value:   base64.StdEncoding.EncodeToString(append(b, sign...)),
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)

	return user.String()
}
