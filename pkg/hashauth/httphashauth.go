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
	retrieveData func(r *http.Request) (hash, data string, err error, skipAuth bool)
}

func NewHTTPHashAuthenticator(
	secret string,
	retrieveData func(r *http.Request) (hash, data string, err error, skipAuth bool)) *HTTPHashAuthenticator {
	return &HTTPHashAuthenticator{
		auther:       NewHashAuthenticator(secret),
		retrieveData: retrieveData,
	}
}

func (h *HTTPHashAuthenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash, data, err, skipAuth := h.retrieveData(r)
		if !skipAuth {
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
				log.Errorf("Invalid hash: %s", hash)
				return
			}
		}

		next.ServeHTTP(w, r)
	}))
}
