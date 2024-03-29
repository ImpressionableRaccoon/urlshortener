package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Символы для ID.
	idLength          = 5                                                                // Длина генерируемого ID.
)

// GenRandomID - функция дли генерации случайного ID для ссылки.
func GenRandomID() (string, error) {
	allowedCharactersLength := big.NewInt(int64(len(allowedCharacters)))

	var b strings.Builder

	for i := 0; i < idLength; i++ {
		n, err := rand.Int(rand.Reader, allowedCharactersLength)
		if err != nil {
			return "", err
		}

		_, err = fmt.Fprint(&b, string(allowedCharacters[n.Int64()]))
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}
