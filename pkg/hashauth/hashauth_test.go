package hashauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Authenticate_Fail(t *testing.T) {
	secret := "hashauth"
	ha := NewHashAuthenticator(secret)

	valid := ha.Authenticate("given hash", "test")

	assert.False(t, valid)
}

func Test_Authenticate_Succes(t *testing.T) {
	secret := "hashauth"
	ha := NewHashAuthenticator(secret)

	testHasher := hmac.New(sha256.New, []byte(secret))
	testHasher.Write([]byte("test"))
	hash := base64.StdEncoding.EncodeToString(testHasher.Sum(nil))

	valid := ha.Authenticate(hash, "test")

	assert.True(t, valid)
}
