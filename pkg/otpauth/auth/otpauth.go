package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type OTPAuthenticator struct {
	tokens             *cache.Cache
	authenticationFunc func(string) bool
}

func NewOTPAuthenticatior(duration time.Duration, authenticationFunc func(string) bool) *OTPAuthenticator {
	return &OTPAuthenticator{
		tokens:             cache.New(duration, duration),
		authenticationFunc: authenticationFunc,
	}
}

func (a *OTPAuthenticator) AuthenticateToken(token string) (ret interface{}, err error) {
	ret, ok := a.tokens.Get(token)
	if !ok {
		err = fmt.Errorf("Token: %s not found in list", token)
		return
	}
	return
}

func (a *OTPAuthenticator) GenerateToken(ID string) (token string, err error) {
	valid := a.authenticationFunc(ID)
	if !valid {
		err = fmt.Errorf("ID: %s is not valid", ID)
		return
	}

	token = createToken()

	err = a.tokens.Add(token, ID, cache.NoExpiration)
	if err != nil {
		token = ""
		err = fmt.Errorf("Error while inserting token: %s with ID: %s, error: %s", token, ID, err.Error())
		return
	}

	log.Infof("Succesfully added token: %s with ID: %s to the tokens map", token, ID)

	return
}

func createToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
