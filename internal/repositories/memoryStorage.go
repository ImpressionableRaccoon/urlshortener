package repositories

import (
	"errors"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDLinkDataDictionary map[ID]LinkData
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[ID]LinkData),
	}

	return st, nil
}

func (st *MemStorage) Add(url URL, userID User) (id ID, err error) {
	for ok := true; ok; _, ok = st.IDLinkDataDictionary[id] {
		id, err = utils.GetRandomID()
		if err != nil {
			return "", err
		}
	}

	st.IDLinkDataDictionary[id] = LinkData{
		URL:  url,
		User: userID,
	}

	return id, nil
}

func (st *MemStorage) Get(id ID) (string, error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, nil
	}

	return "", errors.New("URL not found")
}

func (st *MemStorage) GetUserLinks(user User) (data []UserLink, err error) {
	data = make([]UserLink, 0)

	for id, value := range st.IDLinkDataDictionary {
		if value.User != user {
			continue
		}

		data = append(data, UserLink{
			ID:  id,
			URL: value.URL,
		})
	}

	return
}

func (st *MemStorage) Pool() bool {
	return true
}
