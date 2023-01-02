package repositories

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type FileStorage struct {
	IDLinkDataDictionary map[ID]LinkData
	UserIDs              map[User]bool
	file                 *os.File
	writer               *bufio.Writer
}

func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		IDLinkDataDictionary: make(map[ID]LinkData),
		UserIDs:              make(map[User]bool),
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
			panic(err)
		}

		st.IDLinkDataDictionary[id] = LinkData{
			URL:  url,
			User: userID,
		}
		st.UserIDs[userID] = true
	}

	return st, nil
}

func (st *FileStorage) Add(url URL, userID User) (id ID, err error) {
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
	st.UserIDs[userID] = true

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

func (st *FileStorage) Get(id ID) (URL, error) {
	data, ok := st.IDLinkDataDictionary[id]
	if ok {
		return data.URL, nil
	}
	return "", errors.New("URL not found")
}

func (st *FileStorage) GetUserLinks(user User) (data []UserLink) {
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

func (st *FileStorage) IsUserExists(userID User) bool {
	return st.UserIDs[userID]
}

func (st *FileStorage) Close() error {
	return st.file.Close()
}
