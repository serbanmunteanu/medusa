package auth

import (
	"net/http"
	"os"
)

type Auth struct {
	authTrustedProxy string
}

func NewAuth() *Auth {
	return &Auth{
		authTrustedProxy: os.Getenv("AUTH_TRUSTED_PROXY"),
	}
}

func (a *Auth) Authenticate(request *http.Request) error {
	return nil
}

func (a *Auth) Authorize(request *http.Request) error {
	return nil
}
