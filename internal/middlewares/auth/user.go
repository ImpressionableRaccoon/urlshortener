package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/ImpressionableRaccoon/urlshortener/configs"

	"github.com/google/uuid"
)

func UserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER")
		if err == http.ErrNoCookie || len(cookie.Value) < 16 {
			user, err := setNewUser(w)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		payload, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if len(payload) < 16 {
			user, err := setNewUser(w)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		h := hmac.New(sha256.New, configs.CookieKey)
		h.Write(payload[:16])
		sign := h.Sum(nil)

		if !hmac.Equal(sign, payload[16:]) {
			user, err := setNewUser(w)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := uuid.FromBytes(payload[:16])
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setNewUser(w http.ResponseWriter) (string, error) {
	user := uuid.New()

	b, err := user.MarshalBinary()
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, configs.CookieKey)
	h.Write(b)
	sign := h.Sum(nil)

	content := append(b, sign...)

	encoded := base64.StdEncoding.EncodeToString(content)

	cookie := http.Cookie{
		Name:    "USER",
		Value:   encoded,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	http.SetCookie(w, &cookie)

	return user.String(), nil
}
