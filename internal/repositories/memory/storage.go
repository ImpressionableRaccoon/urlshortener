package memory

import (
	"context"
	"log"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDLinkDataDictionary map[repositories.ID]repositories.LinkData
	existingURLs         map[repositories.URL]repositories.ID
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[repositories.ID]repositories.LinkData),
		existingURLs:         make(map[repositories.URL]repositories.ID),
	}

	return st, nil
}

func (st *MemStorage) Add(ctx context.Context, url repositories.URL, userID repositories.User) (id repositories.ID, err error) {
	value, ok := st.existingURLs[url]
	if ok {
		return value, repositories.ErrURLAlreadyExists
	}

	for ok := true; ok; _, ok = st.IDLinkDataDictionary[id] {
		id, err = utils.GenRandomID()
		if err != nil {
			log.Printf("generate id failed: %v", err)
			return "", err
		}
	}

	st.IDLinkDataDictionary[id] = repositories.LinkData{
		URL:  url,
		User: userID,
	}
	st.existingURLs[url] = id

	return id, nil
}

func (st *MemStorage) Get(ctx context.Context, id repositories.ID) (repositories.URL, error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, nil
	}

	return "", repositories.ErrURLNotFound
}

func (st *MemStorage) GetUserLinks(ctx context.Context, user repositories.User) (data []repositories.UserLink, err error) {
	data = make([]repositories.UserLink, 0)

	for id, value := range st.IDLinkDataDictionary {
		if value.User != user {
			continue
		}

		data = append(data, repositories.UserLink{
			ID:  id,
			URL: value.URL,
		})
	}

	return
}

func (st *MemStorage) Pool(ctx context.Context) bool {
	return true
}
