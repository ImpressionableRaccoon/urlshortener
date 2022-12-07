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

type id = string
type url = string

type Storage struct {
	IDURLsDictionary map[id]url
}

func NewStorage() *Storage {
	storage := &Storage{
		IDURLsDictionary: make(map[string]string),
	}

	return storage
}

func (st *Storage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st.IDURLsDictionary[id] {
		id, err = st.getRandomID()
		if err != nil {
			return "", err
		}
	}

	st.IDURLsDictionary[id] = url

	return id, nil
}

func (st *Storage) Get(id string) (string, error) {
	url, ok := st.IDURLsDictionary[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func (st *Storage) getRandomID() (string, error) {
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
