package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/go-resty/resty/v2"
)

func ExampleHandler_CreateShortURL() {
	url := "https://google.com"
	link := "http://localhost:8080/"

	client := resty.New()

	r, err := client.R().SetBody(url).Post(link)
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 201 {
		log.Println("success")
	} else {
		log.Println("unable to short URL")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())
}

func ExampleHandler_GetURL() {
	link := "http://localhost:8080/yourIDHere"

	client := resty.New().SetRedirectPolicy(resty.RedirectPolicyFunc(
		func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}))

	r, err := client.R().Get(link)
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 307 {
		log.Println("success")
	} else {
		log.Println("unable to get URL")
	}

	log.Printf("Response body: %s\n", strings.Trim(string(r.Body()), "\n\r"))
	log.Printf("Response status code: %d", r.StatusCode())
	log.Printf("Response location header: %s", r.Header().Get("Location"))
}

func ExampleHandler_PingDB() {
	link := "http://localhost:8080/ping"

	client := resty.New()

	r, err := client.R().Get(link)
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 200 {
		log.Println("success")
	} else {
		log.Println("unable to ping storage")
	}

	log.Printf("Response body: %s\n", strings.Trim(string(r.Body()), "\n\r"))
	log.Printf("Response status code: %d", r.StatusCode())
}

func ExampleHandler_ShortenURL() {
	url := "https://google.com"
	link := "http://localhost:8080/api/shorten"

	request := struct {
		URL string `json:"url"`
	}{
		URL: url,
	}

	client := resty.New()

	r, err := client.R().SetBody(request).Post(link)
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 201 {
		log.Println("success")
	} else {
		log.Println("unable to short URL")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())
}

func ExampleHandler_ShortenBatch() {
	link := "http://localhost:8080/api/shorten/batch"

	request := []struct {
		CorrelationID string `json:"correlation_id"`
		URL           string `json:"original_url"`
	}{
		{CorrelationID: "google", URL: "https://google.com/"},
		{CorrelationID: "yandex", URL: "https://yandex.ru"},
	}

	client := resty.New()

	r, err := client.R().SetBody(request).Post(link)
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 201 {
		log.Println("success")
	} else {
		log.Println("unable to short URLs")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())
}

func ExampleHandler_GetUserURLs() {
	// сокращаем ссылки
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		panic(err)
	}

	client := resty.New().SetCookieJar(jar)

	request := []struct {
		CorrelationID string `json:"correlation_id"`
		URL           string `json:"original_url"`
	}{
		{CorrelationID: "google", URL: "https://google.com/"},
		{CorrelationID: "yandex", URL: "https://yandex.ru"},
	}

	r, err := client.R().SetBody(request).Post("http://localhost:8080/api/shorten/batch")
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 201 {
		log.Println("success")
	} else {
		log.Println("unable to short URLs")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())

	// получаем сокращенные ссылки
	r, err = client.R().Get("http://localhost:8080/api/user/urls")
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 200 {
		log.Println("success")
	} else {
		log.Println("unable to get user URLs")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())
}

func ExampleHandler_DeleteUserURLs() {
	// сокращаем ссылки
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		panic(err)
	}

	client := resty.New().SetCookieJar(jar)

	request := []struct {
		CorrelationID string `json:"correlation_id"`
		URL           string `json:"original_url"`
	}{
		{CorrelationID: "google", URL: "https://google.com/"},
		{CorrelationID: "yandex", URL: "https://yandex.ru"},
	}

	r, err := client.R().SetBody(request).Post("http://localhost:8080/api/shorten/batch")
	if err != nil {
		panic(err)
	}

	response := make([]struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}, 0)

	if r.StatusCode() == 201 {
		log.Println("success")
		err = json.Unmarshal(r.Body(), &response)
		if err != nil {
			log.Printf("Error while unmarshal JSON: %v", err)
		}
	} else {
		log.Println("unable to short URLs")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())

	// удаляем сокращенные ссылки
	links := make([]string, 0, len(response))
	for _, link := range response {
		splitted := strings.Split(link.ShortURL, "/")
		links = append(links, splitted[len(splitted)-1])
	}

	r, err = client.R().SetBody(links).Delete("http://localhost:8080/api/user/urls")
	if err != nil {
		panic(err)
	}

	if r.StatusCode() == 202 {
		log.Println("success")
	} else {
		log.Println("unable to delete user URLs")
	}

	log.Printf("Response body: %s\n", string(r.Body()))
	log.Printf("Response status code: %d", r.StatusCode())
}
