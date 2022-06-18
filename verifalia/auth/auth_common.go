package auth

import (
	"net/http"
)

type Provider interface {
	Authenticate(request *http.Request) error
	HandleUnauthorizedRequest() error
	BuildClient() *http.Client
}
