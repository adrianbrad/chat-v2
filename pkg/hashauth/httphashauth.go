package hashauth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type auther interface {
	Authenticate(givenHash string, data string) (valid bool)
}

type HTTPHashAuthenticator struct {
	auther
	next         http.Handler
	retrieveData func(r *http.Request) (hash, data string, err error)
}

func NewHTTPHashAuthenticator(
	secret string,
	next http.Handler,
	retrieveData func(r *http.Request) (hash, data string, err error)) *HTTPHashAuthenticator {
	return &HTTPHashAuthenticator{
		auther:       NewHashAuthenticator(secret),
		next:         next,
		retrieveData: retrieveData,
	}
}

func (h *HTTPHashAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash, data, err := h.retrieveData(r)
	if err != nil {
		http.Error(
			w,
			"Error while retrieving data",
			http.StatusBadRequest,
		)
		log.Errorf("Error: %s while retrieving data from request: %+v", err.Error(), r)
		return
	}

	valid := h.Authenticate(hash, data)
	if !valid {
		http.Error(
			w,
			"Invalid hash",
			http.StatusUnauthorized,
		)
		return
	}

	h.next.ServeHTTP(w, r)
}
