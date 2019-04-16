package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	s := UserService{}

	req, err := http.NewRequest("POST", "http://wtf.com/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	s.ServeHTTP(rr, req)
}
