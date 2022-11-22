package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var ids map[string]string = make(map[string]string)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength          = 5
)

func genID() (string, error) {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	for i := 0; i < idLength; i++ {
		_, err := fmt.Fprint(&b, string(allowedCharacters[rand.Int31n(int32(len(allowedCharacters)))]))
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func Make(url string) (id string, err error) {
	for ok := true; ok; _, ok = ids[id] {
		id, err = genID()
		if err != nil {
			return
		}
	}

	ids[id] = url

	return id, nil
}

func Get(id string) (string, error) {
	url, ok := ids[id]
	if ok {
		return url, nil
	}
	return "", errors.New("URL not found")
}
