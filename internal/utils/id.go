package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength          = 5
)

func GetRandomID() (string, error) {
	rand.Seed(time.Now().UnixNano())
	allowedCharactersLength := int32(len(allowedCharacters))
	var b strings.Builder
	for i := 0; i < idLength; i++ {
		_, err := fmt.Fprint(&b, string(allowedCharacters[rand.Int31n(allowedCharactersLength)]))
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}
