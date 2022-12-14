package memory

import (
	"errors"

	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type Storage struct {
	IDURLsDictionary map[storage.ID]storage.URL
}

func NewStorage() (*Storage, error) {
	st := &Storage{
		IDURLsDictionary: make(map[string]string),
	}

	return st, nil
}

func (st *Storage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st.IDURLsDictionary[id] {
		id, err = utils.GetRandomID()
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
