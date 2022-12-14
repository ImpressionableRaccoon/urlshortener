package file

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ImpressionableRaccoon/urlshortener/internal/utils"

	"github.com/ImpressionableRaccoon/urlshortener/internal/storage"
)

type Storage struct {
	IDURLsDictionary map[storage.ID]storage.URL
	file             *os.File
	writer           *bufio.Writer
}

func NewStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	st := &Storage{
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

func (st *Storage) Add(url string) (id string, err error) {
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

func (st *Storage) Get(id string) (string, error) {
	url, ok := st.IDURLsDictionary[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}

func (st *Storage) Close() error {
	return st.file.Close()
}
