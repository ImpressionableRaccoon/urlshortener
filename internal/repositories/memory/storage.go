// Package memory содержит хранилище интерфейса Storager для хранения данных во временной памяти.
package memory

import (
	"context"
	"log"
	"sync"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

// MemStorage - структура для хранилища во временной памяти.
type MemStorage struct {
	sync.RWMutex
	IDLinkDataDictionary map[repositories.ID]repositories.LinkData
	ExistingURLs         map[repositories.URL]repositories.ID
}

// NewMemoryStorage - конструктор для MemStorage.
func NewMemoryStorage() (*MemStorage, error) {
	st := &MemStorage{
		IDLinkDataDictionary: make(map[repositories.ID]repositories.LinkData),
		ExistingURLs:         make(map[repositories.URL]repositories.ID),
	}

	return st, nil
}

// Add - адаптер для AddLink.
func (st *MemStorage) Add(
	ctx context.Context,
	url repositories.URL,
	user repositories.User,
) (id repositories.ID, err error) {
	return st.AddLink(url, user)
}

// AddLink - сократить ссылку.
func (st *MemStorage) AddLink(url repositories.URL, user repositories.User) (id repositories.ID, err error) {
	st.Lock()
	defer st.Unlock()

	value, ok := st.ExistingURLs[url]
	if ok {
		return value, repositories.ErrURLAlreadyExists
	}

	for exists := true; exists; _, exists = st.IDLinkDataDictionary[id] {
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

// Get - получить оригинальную ссылку по ID.
func (st *MemStorage) Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error) {
	st.RLock()
	defer st.RUnlock()

	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, data.Deleted, nil
	}

	return "", false, repositories.ErrURLNotFound
}

// GetUserLinks - получить все ссылки пользователя.
func (st *MemStorage) GetUserLinks(
	ctx context.Context,
	user repositories.User,
) (data []repositories.LinkData, err error) {
	st.RLock()
	defer st.RUnlock()

	data = make([]repositories.LinkData, 0)

	for id, value := range st.IDLinkDataDictionary {
		if value.User != user {
			continue
		}

		if value.Deleted {
			continue
		}

		data = append(data, repositories.LinkData{
			ID:      id,
			URL:     value.URL,
			User:    user,
			Deleted: false,
		})
	}

	return data, nil
}

// DeleteUserLinks - удалить ссылки пользователя.
func (st *MemStorage) DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error {
	for _, id := range ids {
		_ = st.DeleteUserLink(id, user)
	}

	return nil
}

// DeleteUserLink - удалить ссылку пользователя.
func (st *MemStorage) DeleteUserLink(id repositories.ID, user repositories.User) (ok bool) {
	st.Lock()
	defer st.Unlock()

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

// Pool - проверить соединение с базой данных.
func (st *MemStorage) Pool(ctx context.Context) (ok bool) {
	return true
}
