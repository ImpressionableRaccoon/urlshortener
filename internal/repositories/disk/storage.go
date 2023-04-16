// Package disk содержит хранилище интерфейса Storager для взаимодействия с текстовым файлом.
package disk

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/memory"
)

// FileStorage - структура для хранилища в файле.
type FileStorage struct {
	file *os.File
	memory.MemStorage
	fileMutex sync.Mutex
}

// NewFileStorage - конструктор для FileStorage.
func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		file: file,
	}
	st.IDLinkDataDictionary = make(map[repositories.ID]repositories.LinkData)
	st.ExistingURLs = make(map[repositories.URL]repositories.ID)

	err := st.load()
	if err != nil {
		return nil, err
	}

	return st, nil
}

// Add - адаптер для AddLink.
func (st *FileStorage) Add(
	_ context.Context,
	url repositories.URL,
	user repositories.User,
) (id repositories.ID, err error) {
	id, err = st.AddLink(url, user)
	if err != nil {
		return
	}

	err = st.write(fmt.Sprintf("NEW,%s,%s,%s", id, user.String(), base64.StdEncoding.EncodeToString([]byte(url))))
	return
}

// DeleteUserLinks - удалить ссылки пользователя.
func (st *FileStorage) DeleteUserLinks(_ context.Context, ids []repositories.ID, user repositories.User) error {
	for _, id := range ids {
		ok := st.DeleteUserLink(id, user)
		if ok {
			err := st.write(fmt.Sprintf("DELETE,%s,%s", id, user.String()))
			if err != nil {
				log.Printf("unable to write delete: %v", err)
			}
		}
	}
	return nil
}

// Close - мягко завершить работу хранилища.
func (st *FileStorage) Close(_ context.Context) error {
	return st.file.Close()
}

func (st *FileStorage) load() error {
	st.Lock()
	defer st.Unlock()

	reader := bufio.NewReader(st.file)

	i := 0
	for {
		bytes, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("unable to read bytes: %v", err)
			return err
		}
		line := strings.Trim(string(bytes), "\n")
		splitted := strings.Split(line, ",")

		switch splitted[0] {
		case "NEW":
			err = st.loadNew(splitted)
		case "DELETE":
			err = st.loadDelete(splitted)
		}
		if err != nil {
			log.Printf("unable to parse line %d: %v", i, err)
		}

		i++
	}

	return nil
}

func (st *FileStorage) loadNew(splitted []string) error {
	id := splitted[1]
	user, err := uuid.Parse(splitted[2])
	if err != nil {
		return repositories.ErrUnableParseUser
	}

	var data []byte
	data, err = base64.StdEncoding.DecodeString(splitted[3])
	if err != nil {
		return repositories.ErrUnableDecodeURL
	}
	url := repositories.URL(data)

	st.IDLinkDataDictionary[id] = repositories.LinkData{
		URL:  url,
		User: user,
	}
	st.ExistingURLs[url] = id

	return nil
}

func (st *FileStorage) loadDelete(splitted []string) error {
	id := splitted[1]
	user, err := uuid.Parse(splitted[2])
	if err != nil {
		return repositories.ErrUnableParseUser
	}

	link, ok := st.IDLinkDataDictionary[id]
	if !ok {
		return repositories.ErrLinkNotExists
	}
	if link.User != user {
		return repositories.ErrUserNotMatch
	}

	link.Deleted = true
	st.IDLinkDataDictionary[id] = link

	return nil
}

func (st *FileStorage) write(data string) error {
	st.fileMutex.Lock()
	defer st.fileMutex.Unlock()

	_, err := st.file.Write([]byte(data + "\n"))
	if err != nil {
		return err
	}

	return nil
}
