package hashauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HashAuthenticator struct {
	h hash.Hash
}

func NewHashAuthenticator(secret string) *HashAuthenticator {
	return &HashAuthenticator{
		h: hmac.New(sha256.New, []byte(secret)),
	}
}

func (a *HashAuthenticator) Authenticate(givenHash string, data string) (valid bool) {
	a.h.Write([]byte(data))
	valid = base64.StdEncoding.EncodeToString(a.h.Sum(nil)) == givenHash
	a.h.Reset()
	return
}

func generateHash(h hash.Hash, data []byte) (hash string) {
	h.Write(data)
	hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
	h.Reset()
	return
}
