package auth

import (
	"fmt"
	"net/http"
	"time"
)

// Basic authentication

type basicAuthProvider struct {
	Username string
	Password string
}

func (provider basicAuthProvider) Authenticate(request *http.Request) error {
	if provider.Username == "" {
		return fmt.Errorf("empty username, please specify a valid value before authenticating")
	}

	request.SetBasicAuth(provider.Username, provider.Password)
	return nil
}

func (provider basicAuthProvider) HandleUnauthorizedRequest() error {
	return nil
}

func (provider basicAuthProvider) BuildClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func NewBasicAuthProvider(username string, password string) Provider {
	return &basicAuthProvider{
		Username: username,
		Password: password,
	}
}
