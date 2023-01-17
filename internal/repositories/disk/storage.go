package disk

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/repositories"
	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type FileStorage struct {
	IDLinkDataDictionary map[repositories.ID]repositories.LinkData
	existingURLs         map[repositories.URL]repositories.ID
	file                 *os.File
	writer               *bufio.Writer
}

func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		IDLinkDataDictionary: make(map[repositories.ID]repositories.LinkData),
		existingURLs:         make(map[repositories.URL]repositories.ID),
		file:                 file,
		writer:               bufio.NewWriter(file),
	}

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

		id := splitted[0]
		url := splitted[1]

		userID, err := uuid.Parse(splitted[2])
		if err != nil {
			log.Printf("unable to parse user: %v", err)
			return nil, err
		}

		st.IDLinkDataDictionary[id] = repositories.LinkData{
			URL:  url,
			User: userID,
		}
		st.existingURLs[url] = id
	}

	return st, nil
}

func (st *FileStorage) Add(ctx context.Context, url repositories.URL, userID repositories.User) (id repositories.ID, err error) {
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

	data := []byte(id + "," + url + "," + userID.String() + "\n")
	if _, err = st.writer.Write(data); err != nil {
		log.Printf("write failed: %v", err)
		return "", err
	}
	err = st.writer.Flush()
	if err != nil {
		log.Printf("flush failed: %v", err)
		return "", err
	}

	return id, nil
}

func (st *FileStorage) Get(ctx context.Context, id repositories.ID) (url repositories.URL, deleted bool, err error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, data.Deleted, nil
	}
	return "", false, repositories.ErrURLNotFound
}

func (st *FileStorage) GetUserLinks(ctx context.Context, user repositories.User) (data []repositories.UserLink, err error) {
	data = make([]repositories.UserLink, 0)

	for id, value := range st.IDLinkDataDictionary {
		if value.User != user {
			continue
		}

		if value.Deleted == true {
			continue
		}

		data = append(data, repositories.UserLink{
			ID:  id,
			URL: value.URL,
		})
	}

	return data, err
}

func (st *FileStorage) Pool(ctx context.Context) bool {
	return true
}

func (st *FileStorage) Close() error {
	return st.file.Close()
}

func (st *FileStorage) DeleteUserLinks(ctx context.Context, ids []repositories.ID, user repositories.User) error {
	for _, id := range ids {
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
	return nil
}
