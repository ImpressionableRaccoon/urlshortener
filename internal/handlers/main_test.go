package handlers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	// создаем хранилище для наших тестов
	st := storage.NewStorage()
	st["test"] = "https://google.com"

	// пробуем получить тестовый URL из хранилища
	t.Run("get test URL", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/test", nil)

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)

		// сравниваем ссылку, на которую нас редиректит
		assert.Equal(t, st["test"], res.Header.Get("Location"))
	})

	// пробуем получить несуществующий URL
	t.Run("get URL by wrong ID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/test123", nil)

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	// пробуем получить короткую ссылку для пустого URL
	t.Run("try to get short link for empty URL", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	// пробуем отправить PUT-запрос
	t.Run("wrong PUT request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/test123", nil)

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	originalLink := "https://impressionablracoon.com"
	var shortLink string

	// получаем короткую ссылку для URL
	t.Run("get short link for URL", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(originalLink))

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		// проверяем body на пустоту
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		shortLink = string(resBody)
	})

	// получаем обратно URL из короткой ссылки
	t.Run("get URL from short link", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, shortLink, nil)

		w := httptest.NewRecorder()
		h := RootHandler(st)
		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		// сравниваем код ответа
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)

		// сравниваем ссылку, на которую нас редиректит
		assert.Equal(t, originalLink, res.Header.Get("Location"))
	})
}
