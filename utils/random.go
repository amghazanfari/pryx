package utils

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
)

func CreateRandomInt(startRange, endRange int) int {
	ranRange := endRange - startRange
	if ranRange <= 0 {
		panic("the range for random is non positive, clearly something is wrong")
	}
	return rand.IntN(ranRange)

}

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := cryptoRand.Read(b)

	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}

	if nRead < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}

	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
