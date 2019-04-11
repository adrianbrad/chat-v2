package testutils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func ReadRequestBody(t *testing.T, body *bytes.Buffer) []byte {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	return bodyBytes
}

func NewTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
