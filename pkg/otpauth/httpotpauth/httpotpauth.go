package httpotpauth

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/adrianbrad/chat-v2/pkg/otpauth/auth"

	log "github.com/sirupsen/logrus"
)

type auther interface {
	GenerateToken(ID string) (token string, err error)
	AuthenticateToken(token string) (ret interface{}, err error)
}

type HTTPOTPAuthenticator struct {
	auther
	next http.Handler
}

func NewHTTPOTPAuthenticator(duration time.Duration, authenticationFunc func(string) bool) *HTTPOTPAuthenticator {
	return &HTTPOTPAuthenticator{
		auther: auth.NewOTPAuthenticatior(duration, authenticationFunc),
	}
}

func (ha *HTTPOTPAuthenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost: // * generate token
			ha.handleGenerate(w, r)
		case http.MethodGet: // * authenticate token
			ha.handleAuthenticate(w, r, next)
		default:
			http.Error(
				w,
				"Invalid method",
				http.StatusMethodNotAllowed,
			)
			return
		}
	}))
}

func (ha *HTTPOTPAuthenticator) handleAuthenticate(w http.ResponseWriter, r *http.Request, next http.Handler) {
	token := r.FormValue("key")

	ID, err := ha.AuthenticateToken(token)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusUnauthorized,
		)
		log.Infof("Invalid token: %s from ip: %s, error: %s", token, r.RemoteAddr, err.Error())
		return
	}

	IDstr, ok := ID.(string)
	if !ok {
		http.Error(
			w,
			"ID retrieved is not string",
			http.StatusInternalServerError,
		)
		log.Infof("Retrieved id is not string %s", ID)
		return
	}

	r.Header.Add("X-OTPAuth-ID", IDstr)
	next.ServeHTTP(w, r)
}

func (ha *HTTPOTPAuthenticator) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(
			w,
			"Body not present",
			http.StatusBadRequest,
		)
		log.Errorf("Received nil body from request: %+v", r)
		return
	}

	IDBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(
			w,
			"Error while reading body",
			http.StatusInternalServerError,
		)
		log.Errorf("Error while reading body from request: %+v", r)
		return
	}

	ID := string(IDBytes)

	token, err := ha.GenerateToken(ID)
	if err != nil {
		http.Error(
			w,
			"Invalid ID",
			http.StatusBadRequest,
		)
		log.Errorf("Request: %+v provided invalid ID: %s, error: %s", r, ID, err.Error())
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusCreated)
	return
}
