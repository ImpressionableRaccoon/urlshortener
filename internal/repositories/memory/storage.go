package memory

import (
	"context"
	"log"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type MemStorage struct {
	IDLinkDataDictionary map[repositories.ID]repositories.LinkData
	ExistingURLs         map[repositories.URL]repositories.ID
}

func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[repositories.ID]repositories.LinkData),
		ExistingURLs:         make(map[repositories.URL]repositories.ID),
	}

	return st, nil
}

func (st *MemStorage) add() {

}

func (st *MemStorage) Add(ctx context.Context, url repositories.URL, user repositories.User) (id repositories.ID, err error) {
	value, ok := st.ExistingURLs[url]
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
		User: user,
	}
	st.ExistingURLs[url] = id

	return id, nil
}

func (st *MemStorage) Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, data.Deleted, nil
	}
	return "", false, repositories.ErrURLNotFound
}

func (st *MemStorage) GetUserLinks(ctx context.Context, user repositories.User) (data []repositories.UserLink, err error) {
	data = make([]repositories.UserLink, 0)

	for id, value := range st.IDLinkDataDictionary {
		if value.User != user {
			continue
		}

		if value.Deleted {
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

func (st *MemStorage) DeleteUserLink(id repositories.ID, user repositories.User) (ok bool) {
	link, ok := st.IDLinkDataDictionary[id]
	if !ok {
		return false
	}
	if link.User != user {
		return false
	}
	link.Deleted = true
	st.IDLinkDataDictionary[id] = link
	return true
}

func (st *MemStorage) DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error {
	for _, id := range ids {
		_ = st.DeleteUserLink(id, user)
	}
	return nil
}
