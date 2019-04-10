package hashauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPHash_ErrorWhileRetrievingData(t *testing.T) {
	h := NewHTTPHashAuthenticator(
		"",
		nil,
		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			err = fmt.Errorf("error")
			return
		},
	)

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, nil)
	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Error while retrieving data\n", string(bodyBytes))
}

func Test_HTTPHash_Success(t *testing.T) {
	secret := "secret"

	h := NewHTTPHashAuthenticator(
		secret,

		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { return }), //returns status code 200

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			data = "some data"
			h := hmac.New(sha256.New, []byte(secret))
			hash = generateHash(h, []byte(data))

			return hash, data, nil, false
		},
	)

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, nil)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_HTTPHash_SkipAuth(t *testing.T) {
	h := NewHTTPHashAuthenticator(
		"",

		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { return }), //returns status code 200

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			return "", "", nil, true
		},
	)

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, nil)

	assert.Equal(t, http.StatusOK, rr.Code)
}
