package integration_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adrianbrad/chat-v2/internal/repository/userrepository"
	"github.com/adrianbrad/chat-v2/internal/user"
	"github.com/adrianbrad/chat-v2/pkg/hashauth"
	"github.com/adrianbrad/chat-v2/pkg/otpauth/httpotpauth"
	"github.com/stretchr/testify/assert"

	"github.com/adrianbrad/chat-v2/test"
	testutils "github.com/adrianbrad/chat-v2/test/utils"
)

func TestAddUserSucces(t *testing.T) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	ur := userrepository.NewUserRepositoryDB(db)

	us := user.NewUserService(ur)

	secret := "chat"
	hAuth := hashauth.NewHTTPHashAuthenticator(
		secret,

		us,

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			data = string(bodyBytes)
			hash = r.Header.Get("Authenticate")

			r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return hash, data, nil, false
		})

	body, err := json.Marshal(map[string]interface{}{
		"id":       "1",
		"nickname": "chat_user",
		"permissions": map[string]struct{}{
			user.SendMessage.String(): struct{}{},
		},
	})

	r := testutils.NewTestRequest(t, http.MethodPost, "", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h := hmac.New(sha256.New, []byte(secret))
	hash := hashauth.GenerateHash(h, body)

	r.Header.Set("Authenticate", hash)
	hAuth.ServeHTTP(rr, r)
	assert.Equal(t, http.StatusCreated, rr.Code)

	u, _ := ur.GetOne("1")
	assert.NotNil(t, u)
	assert.Equal(t, u.ID, "1")
	assert.Equal(t, u.Nickname, "chat_user")
	assert.Equal(t, u.Permissions, map[string]struct{}{
		user.SendMessage.String(): struct{}{},
	})
}

func TestRequestTokenSuccess(t *testing.T) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	ur := userrepository.NewUserRepositoryDB(db)

	authFunc := func(ID string) bool {
		_, err := ur.GetOne(ID)
		if err != nil {
			return false
		}
		return true
	}

	tokenAuth := httpotpauth.NewHTTPOTPAuthenticator(10, authFunc, nil)

	secret := "chat"
	hAuth := hashauth.NewHTTPHashAuthenticator(
		secret,

		tokenAuth,

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			data = string(bodyBytes)
			hash = r.Header.Get("Authenticate")

			r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return hash, data, nil, false
		})

	body := []byte("user_a")
	r := testutils.NewTestRequest(t, http.MethodPost, "", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h := hmac.New(sha256.New, []byte(secret))
	hash := hashauth.GenerateHash(h, body)

	r.Header.Set("Authenticate", hash)

	hAuth.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("Authorization"))
}

func TestAuthenticateTokenSuccess(t *testing.T) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	ur := userrepository.NewUserRepositoryDB(db)

	authFunc := func(ID string) bool {
		_, err := ur.GetOne(ID)
		if err != nil {
			return false
		}
		return true
	}

	websocketsLevel := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("reached"))
	})

	tokenAuth := httpotpauth.NewHTTPOTPAuthenticator(10*time.Second, authFunc, websocketsLevel)

	secret := "chat"
	hAuth := hashauth.NewHTTPHashAuthenticator(
		secret,

		tokenAuth,

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			data = string(bodyBytes)
			hash = r.Header.Get("Authenticate")

			r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return hash, data, nil, false
		})

	body := []byte("user_a")
	r := testutils.NewTestRequest(t, http.MethodPost, "", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h := hmac.New(sha256.New, []byte(secret))
	hash := hashauth.GenerateHash(h, body)

	r.Header.Set("Authenticate", hash)

	hAuth.ServeHTTP(rr, r)

	authToken := rr.Header().Get("Authorization")

	r = testutils.NewTestRequest(t, http.MethodGet, "?key="+authToken, bytes.NewReader(body))
	rr = httptest.NewRecorder()

	time.Sleep(1 * time.Second)

	tokenAuth.ServeHTTP(rr, r)

	fmt.Println(rr)
}
