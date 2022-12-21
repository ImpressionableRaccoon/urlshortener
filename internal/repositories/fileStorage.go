package repositories

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"
)

type FileStorage struct {
	IDURLsDictionary map[ID]URL
	file             *os.File
	writer           *bufio.Writer
}

func NewFileStorage(file *os.File) (*FileStorage, error) {
	st := &FileStorage{
		IDURLsDictionary: make(map[string]string),
		file:             file,
		writer:           bufio.NewWriter(file),
	}

	reader := bufio.NewReader(file)

	for {
		id, err := reader.ReadBytes(',')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		url, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		sID := strings.Trim(string(id), ",")
		sURL := strings.Trim(string(url), "\n")

		st.IDURLsDictionary[sID] = sURL
	}

	return st, nil
}

func (st *FileStorage) Add(url string) (id string, err error) {
	for ok := true; ok; _, ok = st.IDURLsDictionary[id] {
		id, err = utils.GetRandomID()
		if err != nil {
			return "", err
		}
	}

	st.IDURLsDictionary[id] = url

	data := []byte(id + "," + url + "\n")
	if _, err := st.writer.Write(data); err != nil {
		return "", err
	}
	err = st.writer.Flush()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (st *FileStorage) Get(id string) (string, error) {
	url, ok := st.IDURLsDictionary[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func (st *FileStorage) Close() error {
	return st.file.Close()
}
