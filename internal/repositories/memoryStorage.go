package repositories

import (
	"errors"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDURLsDictionary map[ID]URL
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDURLsDictionary: make(map[string]string),
	}

	return st, nil
}

func (st *MemStorage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st.IDURLsDictionary[id] {
		id, err = utils.GetRandomID()
		if err != nil {
			return "", err
		}
	}

	st.IDURLsDictionary[id] = url

	return id, nil
}

func (st *MemStorage) Get(id string) (string, error) {
	url, ok := st.IDURLsDictionary[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}
