package disk

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories/memory"
)

type FileStorage struct {
	memory.MemStorage
	file      *os.File
	fileWrite sync.Mutex
}

func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		file: file,
	}
	st.IDLinkDataDictionary = make(map[repositories.ID]repositories.LinkData)
	st.ExistingURLs = make(map[repositories.URL]repositories.ID)

	reader := bufio.NewReader(file)

	for {
		bytes, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("unable to read bytes: %v", err)
			return nil, err
		}
		line := strings.Trim(string(bytes), "\n")
		splitted := strings.Split(line, ",")

		command := splitted[0]
		id := splitted[1]
		user, err := uuid.Parse(splitted[2])
		if err != nil {
			log.Printf("unable to parse user: %v", err)
			continue
		}

		switch command {
		case "NEW":
			url := splitted[3]
			st.IDLinkDataDictionary[id] = repositories.LinkData{
				URL:  url,
				User: user,
			}
			st.ExistingURLs[url] = id
		case "DELETE":
			link, ok := st.IDLinkDataDictionary[id]
			if !ok {
				continue
			}
			if link.User != user {
				continue
			}
			link.Deleted = true
			st.IDLinkDataDictionary[id] = link
		}
	}

	return st, nil
}

func (st *FileStorage) write(data string) error {
	_, err := st.file.Write([]byte(data + "\n"))
	return err
}

func (st *FileStorage) Close() error {
	return st.file.Close()
}

func (st *FileStorage) Add(ctx context.Context, url repositories.URL, user repositories.User) (id repositories.ID, err error) {
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

	err = st.write(fmt.Sprintf("NEW,%s,%s,%s", id, user.String(), url))
	return id, err
}

func (st *FileStorage) DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error {
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
