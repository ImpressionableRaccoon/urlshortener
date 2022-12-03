package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength          = 5
)

type Storage map[string]string

func (st Storage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st[id] {
		rand.Seed(time.Now().UnixNano())
		var b strings.Builder
		for i := 0; i < idLength; i++ {
			_, err := fmt.Fprint(&b, string(allowedCharacters[rand.Int31n(int32(len(allowedCharacters)))]))
			if err != nil {
				return "", err
			}
		}
		id = b.String()
	}

	st[id] = url

	return id, nil
}

func (st Storage) Get(id string) (string, error) {
	url, ok := st[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func NewStorage() Storage {
	return make(Storage)
}
