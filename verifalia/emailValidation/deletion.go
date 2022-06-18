package emailValidation

import (
	"errors"
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"net/http"
)

// Delete removes an emailValidation validation job from the Verifalia servers.
func (client *Client) Delete(id string) error {
	response, err := client.RestClient.Invoke(rest.InvocationOptions{
		Method:   http.MethodDelete,
		Resource: fmt.Sprintf("email-validations/%v", id),
	})

	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		{
			// The job has been correctly deleted
			return nil
		}
	}

	// TODO: Return the response body along with the error message

	return errors.New(fmt.Sprintf("Unexpected HTTP response status code: %v", response.StatusCode))
}
