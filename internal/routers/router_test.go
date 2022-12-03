package routers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string, http.Header) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	// создаем клиент, который не будет переходить по редиректам
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody), resp.Header
}

func TestRouter(t *testing.T) {
	// создаем хранилище
	st, err := storage.GetStorage()
	if err != nil {
		panic(err)
	}

	st.Values["test"] = "https://google.com"

	// создаем роутер и сервер для тестов
	r, err := NewRouter()
	if err != nil {
		panic(err)
	}

	ts := httptest.NewServer(r)
	defer ts.Close()

	// пробуем получить тестовый URL из хранилища
	t.Run("get test URL", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, http.MethodGet, "/test", nil)
		// сравниваем код ответа
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		// сравниваем ссылку, на которую редиректит
		assert.Equal(t, st.Values["test"], header.Get("Location"))
	})

	// пробуем получить несуществующий URL
	t.Run("get URL by wrong ID", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodGet, "/test123", nil)
		// сравниваем код ответа
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	// пробуем получить короткую ссылку для пустого URL
	t.Run("try to get short link for empty URL", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(""))
		// сравниваем код ответа
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	// пробуем отправить PUT-запрос
	t.Run("wrong PUT request", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodPut, "/test123", nil)
		// сравниваем код ответа
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	originalLink := "https://impressionablracoon.com"
	var shortLinkID string

	// получаем короткую ссылку для URL
	t.Run("get short link for URL", func(t *testing.T) {
		statusCode, body, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(originalLink))
		// сравниваем код ответа
		assert.Equal(t, http.StatusCreated, statusCode)
		// вытаскиваем ID и сохраняем его
		splitted := strings.Split(body, "/")
		shortLinkID = splitted[len(splitted)-1]
	})

	// получаем обратно URL из короткой ссылки
	t.Run("get URL from short link", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, http.MethodGet, "/"+shortLinkID, nil)
		// сравниваем код ответа
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		// сравниваем ссылку, на которую редиректит
		assert.Equal(t, originalLink, header.Get("Location"))
	})
}
