package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateToken_InvalidID(t *testing.T) {
	a := NewOTPAuthenticatior(
		1*time.Second,
		func(string) bool { return false },
	)
	token, err := a.GenerateToken("69")

	assert.Empty(t, token)
	assert.EqualError(t, err, "ID: 69 is not valid")
}

func Test_GenerateToken_Success(t *testing.T) {
	a := NewOTPAuthenticatior(
		1*time.Second,
		func(string) bool { return true },
	)
	token, err := a.GenerateToken("69")

	assert.NotEmpty(t, token)
	assert.NoError(t, err)
}

func Test_AuthenticateToken_TokenNotFound(t *testing.T) {
	a := NewOTPAuthenticatior(
		1*time.Second,
		func(string) bool { return true },
	)

	id, err := a.AuthenticateToken("token")

	assert.Empty(t, id)
	assert.EqualError(t, err, "Token: token not found in list")
}

func Test_AuthenticateToken_TokenExpired(t *testing.T) {
	a := NewOTPAuthenticatior(
		100*time.Millisecond,
		func(string) bool { return true },
	)

	token, _ := a.GenerateToken("69")
	time.Sleep(200 * time.Millisecond)
	id, err := a.AuthenticateToken(token)

	assert.Empty(t, id)
	assert.EqualError(t, err, fmt.Sprintf("Token: %s not found in list", token))
}

func Test_AuthenticateToken_Success(t *testing.T) {
	a := NewOTPAuthenticatior(
		1*time.Second,
		func(string) bool { return true },
	)

	token, _ := a.GenerateToken("69")
	id, err := a.AuthenticateToken(token)

	assert.Equal(t, id, "69")
	assert.NoError(t, err)
}
