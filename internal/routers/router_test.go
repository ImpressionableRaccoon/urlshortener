package routers

import (
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
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
	st := storage.NewStorage()
	st.Values["test"] = "https://google.com"

	handler := handlers.NewHandler(st)

	r := NewRouter(handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("get test URL", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, http.MethodGet, "/test", nil)
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, st.Values["test"], header.Get("Location"))
	})

	t.Run("get URL by wrong ID", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodGet, "/test123", nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("try to get short link for empty URL", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(""))
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("wrong PUT request", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, http.MethodPut, "/test123", nil)
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	originalLink := "https://impressionablracoon.com"
	var shortLinkID string

	t.Run("get short link for URL", func(t *testing.T) {
		statusCode, body, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(originalLink))
		assert.Equal(t, http.StatusCreated, statusCode)
		splitted := strings.Split(body, "/")
		shortLinkID = splitted[len(splitted)-1]
	})

	t.Run("get URL from short link", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, http.MethodGet, "/"+shortLinkID, nil)
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, originalLink, header.Get("Location"))
	})
}
