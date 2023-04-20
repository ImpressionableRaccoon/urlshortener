package routers

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"
	"github.com/ImpressionableRaccoon/urlshortener/internal/handlers"
	"github.com/ImpressionableRaccoon/urlshortener/internal/middlewares"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

// TestLink - структура для хранения тестовой ссылки.
type TestLink struct {
	URL       repositories.URL // Исходный URL.
	ShortLink repositories.URL // Сокращенный URL.
	ID        repositories.ID  // ID сокращенного URL.
	Delete    bool             // Должна ли ссылка быть удаленной.
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	jar http.CookieJar,
	method, path string,
	body io.Reader,
	headers map[string]string,
) (statusCode int, respBody []byte, header http.Header) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	return resp.StatusCode, respBody, resp.Header
}

func genTestLinks(count int) (res []TestLink, err error) {
	res = make([]TestLink, 0, count)

	var link TestLink
	for i := 0; i < count; i++ {
		link, err = genTestLink()
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}

	return res, nil
}

func genTestLink() (link TestLink, err error) {
	sites := []string{
		"google.com",
		"yandex.ru",
		"example.com",
		"github.com",
		"awesome.go",
		"go.dev",
	}
	sitesLength := big.NewInt(int64(len(sites)))

	var n *big.Int
	n, err = rand.Int(rand.Reader, sitesLength)
	if err != nil {
		return
	}
	site := sites[n.Int64()]

	var page string
	page, err = utils.GenRandomID()
	if err != nil {
		return
	}

	var del *big.Int
	del, err = rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return
	}

	link = TestLink{
		URL:    fmt.Sprintf("https://%s/%s", site, page),
		Delete: del.Int64() != 0,
	}

	return
}

// TestRouter - тесты для роутера NewRouter.
func TestRouter(t *testing.T) {
	cfg := configs.Config{
		ServerAddress: ":31222",
		ServerBaseURL: "http://localhost:31222",
		CookieKey:     []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		TrustedSubnet: &net.IPNet{
			IP:   net.IP{127, 0, 0, 1},
			Mask: net.IPMask{255, 255, 255, 255},
		},
	}

	s, err := storage.NewStorager(cfg)
	require.NoError(t, err)

	a := authenticator.New(cfg)

	h := handlers.NewHandler(s, cfg)
	m := middlewares.NewMiddlewares(cfg, a)
	r := NewRouter(h, m)

	ts := httptest.NewServer(r)
	defer ts.Close()

	jar, err := cookiejar.New(&cookiejar.Options{})
	require.NoError(t, err)

	t.Run("GET /ping: ping", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodGet, "/ping", nil, nil)

		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("GET /{id}: get URL by wrong ID", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodGet, "/test123", nil, nil)

		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("POST /: try to get short link for empty URL", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodPost, "/", strings.NewReader(""), nil)

		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("GET /api/user/urls: no user URLs", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodGet, "/api/user/urls", nil, nil)

		assert.Equal(t, http.StatusNoContent, statusCode)
	})

	t.Run("wrong PUT request", func(t *testing.T) {
		statusCode, _, _ := testRequest(t, ts, jar, http.MethodPut, "/test123", nil, nil)

		assert.Equal(t, http.StatusMethodNotAllowed, statusCode)
	})

	links, err := genTestLinks(10)
	require.NoError(t, err)

	t.Run("POST /: short URLs", func(t *testing.T) {
		for i := range links {
			statusCode, body, _ := testRequest(
				t, ts, jar, http.MethodPost, "/",
				strings.NewReader(links[i].URL), nil,
			)

			links[i].ShortLink = string(body)
			splitted := strings.Split(links[i].ShortLink, "/")
			links[i].ID = splitted[len(splitted)-1]

			assert.Equal(t, http.StatusCreated, statusCode)
			assert.Equal(t, fmt.Sprintf("%s/%s", cfg.ServerBaseURL, links[i].ID), links[i].ShortLink)
		}
	})

	t.Run("POST /: short conflict URLs", func(t *testing.T) {
		for i := range links {
			statusCode, body, _ := testRequest(
				t, ts, jar, http.MethodPost, "/",
				strings.NewReader(links[i].URL), nil,
			)

			links[i].ShortLink = string(body)
			splitted := strings.Split(links[i].ShortLink, "/")
			links[i].ID = splitted[len(splitted)-1]

			assert.Equal(t, http.StatusConflict, statusCode)
			assert.Equal(t, fmt.Sprintf("%s/%s", cfg.ServerBaseURL, links[i].ID), links[i].ShortLink)
		}
	})

	t.Run("GET /{id}: get URLs", func(t *testing.T) {
		for _, link := range links {
			statusCode, _, header := testRequest(
				t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", link.ID),
				nil, nil,
			)

			assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
			assert.Equal(t, link.URL, header.Get("Location"))
		}
	})

	t.Run("GET /api/user/urls: get user URLs", func(t *testing.T) {
		statusCode, body, _ := testRequest(t, ts, jar, http.MethodGet, "/api/user/urls", nil, nil)

		assert.Equal(t, http.StatusOK, statusCode)

		data := make([]handlers.UserLink, 0)
		err = json.Unmarshal(body, &data)
		require.NoError(t, err)

		for _, link := range links {
			assert.Contains(t, data, handlers.UserLink{ShortURL: link.ShortLink, OriginalURL: link.URL})
		}
	})

	t.Run("DELETE /api/user/urls: delete user URLs", func(t *testing.T) {
		linksIDs := make([]repositories.ID, 0)

		for _, link := range links {
			if link.Delete {
				linksIDs = append(linksIDs, link.ID)
			}
		}

		var data []byte
		data, err = json.Marshal(linksIDs)
		require.NoError(t, err)

		statusCode, _, _ := testRequest(
			t, ts, jar, http.MethodDelete, "/api/user/urls",
			bytes.NewReader(data), nil,
		)

		assert.Equal(t, http.StatusAccepted, statusCode)
	})

	t.Run("GET /{id}: check if only needed links deleted", func(t *testing.T) {
		for _, link := range links {
			statusCode, _, header := testRequest(
				t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", link.ID),
				nil, nil,
			)

			if link.Delete {
				assert.Equal(t, http.StatusGone, statusCode)
				continue
			}

			assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
			assert.Equal(t, link.URL, header.Get("Location"))
		}
	})

	t.Run("GET /api/internal/stats: get stats", func(t *testing.T) {
		stats := repositories.ServiceStats{
			URLs:  10,
			Users: 1,
		}
		expected, err := json.Marshal(stats)
		require.NoError(t, err)

		statusCode, body, _ := testRequest(
			t, ts, jar, http.MethodGet, "/api/internal/stats",
			nil, nil,
		)

		assert.Equal(t, http.StatusOK, statusCode)
		assert.JSONEq(t, string(expected), string(body))
	})

	shortenLinks, err := genTestLinks(10)
	require.NoError(t, err)

	t.Run("POST /api/shorten: short URLs", func(t *testing.T) {
		for i, link := range shortenLinks {
			var request []byte
			request, err = json.Marshal(handlers.ShortenURLRequest{
				URL: link.URL,
			})
			require.NoError(t, err)

			statusCode, body, _ := testRequest(
				t, ts, jar, http.MethodPost, "/api/shorten",
				bytes.NewReader(request), nil,
			)

			response := handlers.ShortenURLResponse{}
			err = json.Unmarshal(body, &response)
			require.NoError(t, err)

			link.ShortLink = response.Result
			splitted := strings.Split(link.ShortLink, "/")
			link.ID = splitted[len(splitted)-1]
			shortenLinks[i] = link

			assert.Equal(t, http.StatusCreated, statusCode)
			assert.Equal(t, fmt.Sprintf("%s/%s", cfg.ServerBaseURL, link.ID), link.ShortLink)
		}
	})

	t.Run("DELETE /api/user/urls: delete URLs from /api/shorten", func(t *testing.T) {
		linksIDs := make([]repositories.ID, 0)

		for _, link := range shortenLinks {
			if link.Delete {
				linksIDs = append(linksIDs, link.ID)
			}
		}

		var data []byte
		data, err = json.Marshal(linksIDs)
		require.NoError(t, err)

		statusCode, _, _ := testRequest(
			t, ts, jar, http.MethodDelete, "/api/user/urls",
			bytes.NewReader(data), nil,
		)

		assert.Equal(t, http.StatusAccepted, statusCode)
	})

	t.Run("GET /{id}: get URLs from /api/shorten", func(t *testing.T) {
		for _, link := range shortenLinks {
			statusCode, _, header := testRequest(
				t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", link.ID),
				nil, nil,
			)

			if link.Delete {
				assert.Equal(t, http.StatusGone, statusCode)
				continue
			}

			assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
			assert.Equal(t, link.URL, header.Get("Location"))
		}
	})

	shortenBatch, err := genTestLinks(1000)
	require.NoError(t, err)

	t.Run("POST /api/shorten/batch: short URLs", func(t *testing.T) {
		data := make([]handlers.BatchRequest, 0, len(shortenBatch))

		for index, link := range shortenBatch {
			data = append(data, handlers.BatchRequest{CorrelationID: strconv.Itoa(index), OriginalURL: link.URL})
		}

		var request []byte
		request, err = json.Marshal(data)
		require.NoError(t, err)

		statusCode, body, _ := testRequest(
			t, ts, jar, http.MethodPost, "/api/shorten/batch",
			bytes.NewReader(request), nil,
		)

		assert.Equal(t, http.StatusCreated, statusCode)

		response := make([]handlers.BatchResponse, 0, len(data))
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		for _, link := range response {
			var id int
			id, err = strconv.Atoi(link.CorrelationID)
			require.NoError(t, err)

			shortenBatch[id].ShortLink = link.ShortURL
			splitted := strings.Split(shortenBatch[id].ShortLink, "/")
			shortenBatch[id].ID = splitted[len(splitted)-1]
		}
	})

	t.Run("DELETE /api/user/urls: delete URLs from /api/shorten/batch", func(t *testing.T) {
		linksIDs := make([]repositories.ID, 0)

		for _, link := range shortenBatch {
			if link.Delete {
				linksIDs = append(linksIDs, link.ID)
			}
		}

		var data []byte
		data, err = json.Marshal(linksIDs)
		require.NoError(t, err)

		statusCode, _, _ := testRequest(
			t, ts, jar, http.MethodDelete, "/api/user/urls",
			bytes.NewReader(data), nil,
		)

		assert.Equal(t, http.StatusAccepted, statusCode)
	})

	t.Run("GET /{id}: get URLs from /api/shorten/batch", func(t *testing.T) {
		for _, link := range shortenBatch {
			statusCode, _, header := testRequest(
				t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", link.ID),
				nil, nil,
			)

			if link.Delete {
				assert.Equal(t, http.StatusGone, statusCode)
				continue
			}

			assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
			assert.Equal(t, link.URL, header.Get("Location"))
		}
	})

	var link TestLink

	t.Run("POST /: short url with gzip-compressed request", func(t *testing.T) {
		link, err = genTestLink()
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		gzipWriter := gzip.NewWriter(buf)
		_, err = gzipWriter.Write([]byte(link.URL))
		require.NoError(t, err)
		err = gzipWriter.Close()
		require.NoError(t, err)

		headers := map[string]string{"Content-Encoding": "gzip"}

		statusCode, body, _ := testRequest(t, ts, jar, http.MethodPost, "/", buf, headers)

		link.ShortLink = string(body)
		splitted := strings.Split(link.ShortLink, "/")
		link.ID = splitted[len(splitted)-1]

		assert.Equal(t, http.StatusCreated, statusCode)
		assert.Equal(t, fmt.Sprintf("%s/%s", cfg.ServerBaseURL, link.ID), link.ShortLink)
	})

	t.Run("GET /{id}: get original URL from shorted by gzip-compressed request", func(t *testing.T) {
		statusCode, _, header := testRequest(
			t, ts, jar, http.MethodGet, fmt.Sprintf("/%s", link.ID),
			nil, nil,
		)

		assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
		assert.Equal(t, link.URL, header.Get("Location"))
	})
}
