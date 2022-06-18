package emailValidation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia/common"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"io/ioutil"
	"net/http"
)

type ListingResult struct {
	JobOverview Overview
	Error       error
}

type ListingField int

const (
	CreatedOn ListingField = iota
)

type ListingOptions struct {
	Context context.Context

	// The maximum number of items to return with a listing request. The Verifalia API may choose to override the
	// specified limit if it is either too small or too big. A single listing operations
	// may automatically perform different listing requests to the Verifalia API: this value limits the number of items
	// returned by each API request, *not the overall total number of returned items*. To limit the total number of returned
	// items, keep track of the number of processed items
	Limit int

	// The job overview field the results will be sorted by.
	OrderBy ListingField

	// The direction of the listing.
	Direction common.Direction
}

// List returns a list of validation jobs according to the user permissions.
func (client *Client) List() chan ListingResult {
	return client.ListWithOptions(ListingOptions{})
}

// ListWithOptions returns a list of validation jobs according to the specified options and user permissions.
func (client *Client) ListWithOptions(options ListingOptions) chan ListingResult {
	// The channel is unbuffered, so that the caller can choose when / whether to abort the listing

	results := make(chan ListingResult)

	go func() {
		defer close(results)

		// First page

		filterParams := make(map[string][]string)

		if options.Limit > 0 {
			filterParams["limit"] = []string{fmt.Sprintf("%v", options.Limit)}
		}

		switch options.OrderBy {
		case CreatedOn:
			{
				if options.Direction == common.Backward {
					filterParams["sort"] = []string{"-createdOn"}
				} else {
					filterParams["sort"] = []string{"createdOn"}
				}
			}
		}

		invOptions := rest.InvocationOptions{
			Method:      http.MethodGet,
			Resource:    "email-validations",
			QueryParams: filterParams,
			Context:     options.Context,
		}

		// Iterate over the subsequent segments

		for {
			segment, err := client.listSegment(invOptions)

			if err != nil {
				results <- ListingResult{
					Error: err,
				}

				break
			}

			for _, jobOverview := range *segment.Data {
				results <- ListingResult{
					JobOverview: jobOverview,
				}
			}

			if !segment.Meta.IsTruncated {
				break
			}

			// Prepare for next page request

			invOptions = rest.InvocationOptions{
				Method:   http.MethodGet,
				Resource: "email-validations",
				QueryParams: map[string][]string{
					"cursor": {
						segment.Meta.Cursor,
					},
				},
				Context: options.Context,
			}
		}
	}()

	return results
}

func (client *Client) listSegment(invocationOptions rest.InvocationOptions) (*common.ListingSegment[Overview], error) {
	response, err := client.RestClient.Invoke(invocationOptions)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		segment := common.ListingSegment[Overview]{}

		if err := json.Unmarshal(responseData, &segment); err != nil {
			return nil, err
		}

		return &segment, nil
	} else {
		// TODO: Return the response body along with the error message

		return nil, errors.New(fmt.Sprintf("Unexpected HTTP response status code: %v", response.StatusCode))
	}
}
