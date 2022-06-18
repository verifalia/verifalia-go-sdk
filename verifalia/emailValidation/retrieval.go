package emailValidation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia/common"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"io/ioutil"
	"net"
	"net/http"
)

// Get fetches an e-mail validation job previously submitted for processing. This function does not wait for the eventual
// completion of the job: use the WaitForCompletion() function to do that.
func (client *Client) Get(id string) (*Job, error) {
	response, err := client.RestClient.Invoke(rest.InvocationOptions{
		Method:   http.MethodGet,
		Resource: fmt.Sprintf("email-validations/%v", id),
	})

	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		{
			var responseData []byte
			if responseData, err = ioutil.ReadAll(response.Body); err != nil {
				return nil, err
			}

			// Unmarshal a simplified projection of the original object, without segmentation/metadata info

			var partial partialJob

			if err := json.Unmarshal(responseData, &partial); err != nil {
				return nil, err
			}

			// Empty context here because this specific API returns the job data as a whole

			result, err := client.buildJob(partial, context.TODO())

			return result, err
		}
	case http.StatusNotFound, http.StatusGone:
		return nil, nil
	}

	// TODO: Return the response body along with the error message

	return nil, errors.New(fmt.Sprintf("Unexpected HTTP response status code: %v", response.StatusCode))
}

func (client *Client) buildJob(partial partialJob, ctx context.Context) (*Job, error) {
	// TODO: Retrieve the other entries from the API (needed when we will support the filtering API)

	var result = &Job{
		Overview: Overview{
			Id:            partial.Overview.Id,
			SubmittedOn:   partial.Overview.SubmittedOn,
			CompletedOn:   partial.Overview.CompletedOn,
			Priority:      partial.Overview.Priority,
			Name:          partial.Overview.Name,
			Owner:         partial.Overview.Owner,
			ClientIP:      net.ParseIP(partial.Overview.ClientIP),
			CreatedOn:     partial.Overview.CreatedOn,
			Quality:       partial.Overview.Quality,
			Retention:     common.TimeSpanStringToDuration(partial.Overview.Retention),
			Deduplication: partial.Overview.Deduplication,
			Status:        partial.Overview.Status,
			NoOfEntries:   partial.Overview.NoOfEntries,
		},
	}

	if partial.Entries != nil && partial.Entries.Data != nil {
		result.Entries = *partial.Entries.Data
	}

	if partial.Overview.Progress != nil {
		result.Overview.Progress = &Progress{
			Percentage: partial.Overview.Progress.Percentage,
		}

		if partial.Overview.Progress.EstimatedTimeRemaining != "" {
			eta := common.TimeSpanStringToDuration(partial.Overview.Progress.EstimatedTimeRemaining)
			result.Overview.Progress.EstimatedTimeRemaining = &eta
		}
	}

	return result, nil
}
