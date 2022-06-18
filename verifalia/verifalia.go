package verifalia

import (
	"crypto/tls"
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia/auth"
	"github.com/verifalia/verifalia-go-sdk/verifalia/credit"
	"github.com/verifalia/verifalia-go-sdk/verifalia/emailValidation"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"runtime"
)

// Client represents a REST client for Verifalia. To start verifying email addresses, use one of the functions available
// through the EmailValidation field, for example:
//  validation, err := client.EmailValidation.Run("batman@gmail.com")
type Client struct {
	authenticationProvider auth.Provider
	restClient             rest.Client

	// Allows to manage the credits for the Verifalia account.
	Credit credit.Client

	// Allows to submit and manage email validations using the Verifalia service, for example:
	//  validation, err := client.EmailValidation.Run("batman@gmail.com")
	EmailValidation emailValidation.Client
}

// NewClient initializes a new REST client for Verifalia with the specified username and password.
// While authenticating with your Verifalia main account credentials is possible, it is strongly advised
// to create one or more users with just the required permissions, for improved
// security. To create a new user or manage existing ones, please visit https://verifalia.com/client-area#/users
func NewClient(username string, password string) *Client {
	return newClientImpl(auth.NewBasicAuthProvider(username, password), rest.BaseUrls)
}

// NewClientWithCertificateAuth initializes a new REST client for Verifalia with the specified client certificate
// (for enterprise-grade mutual TLS authentication). TLS client certificate authentication is available to premium plans only.
// It is strongly advised to create one or more users with just the required permissions,
// for improved security. To create a new user or manage existing ones, please visit https://verifalia.com/client-area#/users
func NewClientWithCertificateAuth(certificate *tls.Certificate) *Client {
	return newClientImpl(auth.NewCertificateAuthProvider(certificate), rest.BaseCcaUrls)
}

func newClientImpl(authenticationProvider auth.Provider, baseUrls []string) *Client {
	client := rest.NewMultiplexedRestClient(authenticationProvider,
		// TODO: Add the git hash of the current SDK version to the user agent string
		fmt.Sprintf("verifalia-rest-client/go/%s/%s", runtime.Version(), runtime.GOOS),
		baseUrls)

	return &Client{
		authenticationProvider: authenticationProvider,
		restClient:             client,
		Credit: credit.Client{
			RestClient: client,
		},
		EmailValidation: emailValidation.Client{
			RestClient: client,
		},
	}
}
