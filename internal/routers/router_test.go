package routers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, jar http.CookieJar, method, path string, body io.Reader) (int, []byte, http.Header) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return resp.StatusCode, respBody, resp.Header
}

func TestRouter(t *testing.T) {
	s, err := storage.NewStorager()
	if err != nil {
		panic(err)
	}

	testURL := "https://google.com"
	testID, err := s.Add(context.Background(), testURL, uuid.New())
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler(s)
	r := NewRouter(h)

	ts := httptest.NewServer(r)
	defer ts.Close()

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatal(err)
	}

	t.Run("get test URL", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", testID), nil)
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, testURL, header.Get("Location"))
	})

	t.Run("get URL by wrong ID", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodGet, "/test123", nil)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("try to get short link for empty URL", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodPost, "/", strings.NewReader(""))
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("wrong PUT request", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodPut, "/test123", nil)
		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	originalLink := "https://example.com"
	var shortLink, shortLinkID string

	t.Run("get short link for URL", func(t *testing.T) {
		statusCode, body, _ := testRequest(t, ts, jar, http.MethodPost, "/", strings.NewReader(originalLink))
		assert.Equal(t, http.StatusCreated, statusCode)
		shortLink = string(body)
		splitted := strings.Split(shortLink, "/")
		shortLinkID = splitted[len(splitted)-1]
		assert.Equal(t, fmt.Sprintf("%s/%s", configs.ServerBaseURL, shortLinkID), string(body))
	})

	t.Run("get URL from short link", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, jar, http.MethodGet, "/"+shortLinkID, nil)
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, originalLink, header.Get("Location"))
	})

	t.Run("API: get short link for URL", func(t *testing.T) {
		request := handlers.ShortenURLRequest{
			URL: originalLink,
		}

		requestJSON, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}

		reader := strings.NewReader(string(requestJSON))

		statusCode, body, _ := testRequest(t, ts, jar, http.MethodPost, "/api/shorten", reader)
		assert.Equal(t, http.StatusConflict, statusCode)

		var response handlers.ShortenURLResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			panic(err)
		}

		url := response.Result
		splitted := strings.Split(url, "/")
		shortLinkID = splitted[len(splitted)-1]
	})

	t.Run("API: get URL from short link", func(t *testing.T) {
		statusCode, _, header := testRequest(t, ts, jar, http.MethodGet, "/"+shortLinkID, nil)
		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, originalLink, header.Get("Location"))
	})

	t.Run("API: get URL from user URLs", func(t *testing.T) {
		statusCode, body, _ := testRequest(t, ts, jar, http.MethodGet, "/api/user/urls", nil)
		assert.Equal(t, http.StatusOK, statusCode)
		var links []handlers.UserLink
		json.Unmarshal(body, &links)
		assert.Contains(t, links, handlers.UserLink{
			ShortURL:    shortLink,
			OriginalURL: originalLink,
		})
	})
}
