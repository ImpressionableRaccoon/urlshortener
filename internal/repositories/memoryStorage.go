package repositories

import (
	"errors"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDLinkDataDictionary map[ID]LinkData
	UserIDs              []UserID
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[ID]LinkData),
	}

	return st, nil
}

func (st *MemStorage) Add(url URL, userID UserID) (id ID, err error) {
	for ok := true; ok; _, ok = st.IDLinkDataDictionary[id] {
		id, err = utils.GetRandomID()
		if err != nil {
			return "", err
		}
	}

	st.IDLinkDataDictionary[id] = LinkData{
		URL:    url,
		UserID: userID,
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
