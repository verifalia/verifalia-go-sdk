package auth

import (
	"crypto/tls"
	"net/http"
	"time"
)

// X.509 mutual TLS client certificate authentication

type certificateAuthProvider struct {
	Certificate *tls.Certificate
	Client      *http.Client
}

func (provider certificateAuthProvider) Authenticate(request *http.Request) error {
	return nil
}

func (provider certificateAuthProvider) HandleUnauthorizedRequest() error {
	return nil
}

func (provider certificateAuthProvider) BuildClient() *http.Client {
	return provider.Client
}

func NewCertificateAuthProvider(certificate *tls.Certificate) Provider {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*certificate},
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &certificateAuthProvider{
		Certificate: certificate,
		Client:      client,
	}
}
