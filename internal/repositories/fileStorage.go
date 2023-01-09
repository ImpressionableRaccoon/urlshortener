package repositories

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type FileStorage struct {
	IDLinkDataDictionary map[ID]LinkData
	existingURLs         map[URL]ID
	file                 *os.File
	writer               *bufio.Writer
}

func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		IDLinkDataDictionary: make(map[ID]LinkData),
		existingURLs:         make(map[URL]ID),
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
			return nil, err
		}
		line := strings.Trim(string(bytes), "\n")
		splitted := strings.Split(line, ",")

		id := splitted[0]
		url := splitted[1]

		userID, err := uuid.Parse(splitted[2])
		if err != nil {
			return nil, err
		}

		st.IDLinkDataDictionary[id] = LinkData{
			URL:  url,
			User: userID,
		}
		st.existingURLs[url] = id
	}

	return st, nil
}

func (st *FileStorage) Add(ctx context.Context, url URL, userID User) (id ID, err error) {
	value, ok := st.existingURLs[url]
	if ok {
		return value, URLAlreadyExists
	}

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
	st.existingURLs[url] = id

	data := []byte(id + "," + url + "," + userID.String() + "\n")
	if _, err = st.writer.Write(data); err != nil {
		return "", err
	}
	err = st.writer.Flush()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (st *FileStorage) Get(ctx context.Context, id ID) (URL, error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, nil
	}
	return "", URLNotFound
}

func (st *FileStorage) GetUserLinks(ctx context.Context, user User) (data []UserLink, err error) {
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

func (st *FileStorage) Pool(ctx context.Context) bool {
	return true
}

func (st *FileStorage) Close() error {
	return st.file.Close()
}
