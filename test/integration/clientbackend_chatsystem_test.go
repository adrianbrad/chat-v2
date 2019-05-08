package integration

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/server"
	"github.com/adrianbrad/chat-v2/internal/user"
	"github.com/adrianbrad/chat-v2/pkg/hashauth"
	"github.com/stretchr/testify/assert"

	testutils "github.com/adrianbrad/chat-v2/test/utils"
)

func TestAddUserSucces(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	initDependencies(db)

	body, err := json.Marshal(map[string]interface{}{
		"id":       "1",
		"nickname": "chat_user",
		"permissions": map[string]struct{}{
			user.SendMessage.String(): struct{}{},
		},
	})

	r := testutils.NewTestRequest(t, http.MethodPost, baseAddress+"/users", bytes.NewReader(body))

	hash := hashauth.GenerateHash(hasher, body)
	r.Header.Set("Authenticate", hash)

	createUserHandler := hashAuthenticator.Auth(userService)

	stopServer := startServer(server.PathHandler{
		Path:    "/users",
		Handler: createUserHandler})
	defer stopServer()

	resp, err := httpClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	u, _ := userRepository.GetOne("1")
	assert.NotNil(t, u)
	assert.Equal(t, u.ID, "1")
	assert.Equal(t, u.Nickname, "chat_user")
	assert.Equal(t, u.Permissions, map[string]struct{}{
		user.SendMessage.String(): struct{}{},
	})
}

func TestRequestTokenSuccess(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	initDependencies(db)

	body := []byte("user_a")
	r := testutils.NewTestRequest(t, http.MethodPost, baseAddress+"/auth", bytes.NewReader(body))

	h := hmac.New(sha256.New, []byte(secret))
	hash := hashauth.GenerateHash(h, body)

	r.Header.Set("Authenticate", hash)

	tokenGeneratorHandler := httpOTPAuthenticator.Auth(nil)

	stopServer := startServer(server.PathHandler{
		Path:    "/auth",
		Handler: tokenGeneratorHandler})
	defer stopServer()

	resp, err := httpClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Authorization"))
}
