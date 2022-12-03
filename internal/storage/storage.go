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

var storage *Storage

type Storage struct {
	Values map[string]string // map[id]url
}

func GetStorage() (*Storage, error) {
	if storage != nil {
		return storage, nil
	}

	storage = &Storage{
		Values: make(map[string]string),
	}

	return storage, nil
}

func (st *Storage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st.Values[id] {
		id, err = getRandomID()
		if err != nil {
			return "", err
		}
	}

	st.Values[id] = url

	return id, nil
}

func (st *Storage) Get(id string) (url string, e error) {
	url, ok := st.Values[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func getRandomID() (string, error) {
	rand.Seed(time.Now().UnixNano())
	allowedCharactersLength := int32(len(allowedCharacters))
	var b strings.Builder
	for i := 0; i < idLength; i++ {
		_, err := fmt.Fprint(&b, string(allowedCharacters[rand.Int31n(allowedCharactersLength)]))
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}