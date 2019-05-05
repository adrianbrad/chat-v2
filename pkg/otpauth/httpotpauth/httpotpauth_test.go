package httpotpauth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ServeHTTP_InvalidMethod(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		nil,
	)

	r, err := http.NewRequest(http.MethodDelete, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, string(bodyBytes), "Invalid method\n")
}

func Test_HandleGenerate_NilBody(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		func(string) bool { return false },
	)

	r, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Body not present\n", string(bodyBytes))
}

func Test_HandleGenerate_InvalidID(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		func(string) bool { return false },
	)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader("69"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Invalid ID\n", string(bodyBytes))
}

func Test_HandleGenerate_Success(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		func(string) bool { return true },
	)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader("69"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("Authorization"))
}

func Test_HandleAuthenticate_NoTokenProvided(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		func(string) bool { return true },
	)

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Token:  not found in list\n", string(bodyBytes))
}

func Test_HandleAuthenticate_InvalidToken(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1,
		func(string) bool { return false },
	)

	r, err := http.NewRequest(http.MethodGet, "?key=token", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(nil).ServeHTTP(rr, r)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Token: token not found in list\n", string(bodyBytes))
}

func Test_HandleAuthenticate_Success(t *testing.T) {
	a := NewHTTPOTPAuthenticator(
		1*time.Second,
		func(string) bool { return true },
	)

	token, err := a.GenerateToken("id")
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodGet, "?key="+token, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	a.Auth(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, r)

	assert.Equal(t, "id", r.Header.Get("X-OTPAuth-ID"))
}
