package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateAPIKey() (plain, prefix, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", "", err
	}
	plain = "sk_live_" + hex.EncodeToString(b)
	prefix = plain[:8]
	h := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(h[:])
	return
}

func HashAPIKey(plain string) string {
	h := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(h[:])
}

func HasScope(scopesCSV, want string) bool {
	want = strings.TrimSpace(want)
	for _, s := range strings.Split(scopesCSV, ",") {
		if strings.TrimSpace(s) == want {
			return true
		}
	}
	return false
}
