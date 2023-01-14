package repositories

import (
	"context"
	"log"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDLinkDataDictionary map[ID]LinkData
	existingURLs         map[URL]ID
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[ID]LinkData),
		existingURLs:         make(map[URL]ID),
	}

	return st, nil
}

func (st *MemStorage) Add(ctx context.Context, url URL, userID User) (id ID, err error) {
	value, ok := st.existingURLs[url]
	if ok {
		return value, ErrURLAlreadyExists
	}

	for ok := true; ok; _, ok = st.IDLinkDataDictionary[id] {
		id, err = utils.GenRandomID()
		if err != nil {
			log.Printf("generate id failed: %v", err)
			return "", err
		}
	}

	st.IDLinkDataDictionary[id] = LinkData{
		URL:  url,
		User: userID,
	}
	st.existingURLs[url] = id

	return id, nil
}

func (st *MemStorage) Get(ctx context.Context, id ID) (URL, error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, nil
	}

	return "", ErrURLNotFound
}

func (st *MemStorage) GetUserLinks(ctx context.Context, user User) (data []UserLink, err error) {
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

func (st *MemStorage) Pool(ctx context.Context) bool {
	return true
}
