package emailValidation

/*
* Verifalia - Email list cleaning and real-time email verification service
* https://verifalia.com/
* support@verifalia.com
*
* Copyright (c) 2005-2024 Cobisi Research
*
* Cobisi Research
* Via Della Costituzione, 31
* 35010 Vigonza
* Italy - European Union
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in
* all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
* THE SOFTWARE.
 */

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
	"time"
)

// RetrievalOptions allows to define retrieval options for an e-mail verification job.
type RetrievalOptions struct {
	// Defines how much time to ask the Verifalia API to wait for the completion of the job on the server side, during the
	// job retrieval request.
	RetrievalWaitTime time.Duration
}

// Get fetches an email validation job previously submitted for processing.
func (client *Client) Get(id string) (*Job, error) {
	return client.GetWithOptions(id, nil)
}

// GetWithOptions fetches an email validation job previously submitted for processing.
func (client *Client) GetWithOptions(id string, options *RetrievalOptions) (*Job, error) {
	var queryParams map[string][]string

	if options != nil {
		queryParams = make(map[string][]string)
		queryParams["waitTime"] = []string{fmt.Sprintf("%v", options.RetrievalWaitTime.Seconds())}
	}

	response, err := client.RestClient.Invoke(rest.InvocationOptions{
		Method:      http.MethodGet,
		Resource:    fmt.Sprintf("email-validations/%v", id),
		QueryParams: queryParams,
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
		Overview: buildOverview(partial.Overview),
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

func buildOverview(rawOverview overview) Overview {
	return Overview{
		Id:            rawOverview.Id,
		SubmittedOn:   rawOverview.SubmittedOn,
		CompletedOn:   rawOverview.CompletedOn,
		Priority:      rawOverview.Priority,
		Name:          rawOverview.Name,
		Owner:         rawOverview.Owner,
		ClientIP:      net.ParseIP(rawOverview.ClientIP),
		CreatedOn:     rawOverview.CreatedOn,
		Quality:       rawOverview.Quality,
		Retention:     common.TimeSpanStringToDuration(rawOverview.Retention),
		Deduplication: rawOverview.Deduplication,
		Status:        rawOverview.Status,
		NoOfEntries:   rawOverview.NoOfEntries,
	}
}

// GetOverview fetches an overview of an email validation job previously submitted for processing.
func (client *Client) GetOverview(id string) (*Overview, error) {
	return client.GetOverviewWithOptions(id, nil)
}

// GetOverviewWithOptions fetches an overview of an email validation job previously submitted for processing.
func (client *Client) GetOverviewWithOptions(id string, options *RetrievalOptions) (*Overview, error) {
	var queryParams map[string][]string

	if options != nil {
		queryParams = make(map[string][]string)
		queryParams["waitTime"] = []string{fmt.Sprintf("%v", options.RetrievalWaitTime.Seconds())}
	}

	response, err := client.RestClient.Invoke(rest.InvocationOptions{
		Method:      http.MethodGet,
		Resource:    fmt.Sprintf("email-validations/%v/overview", id),
		QueryParams: queryParams,
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

			// Unmarshal the job overview

			var rawOverview = overview{}

			if err := json.Unmarshal(responseData, &rawOverview); err != nil {
				return nil, err
			}

			var overview = buildOverview(rawOverview)

			return &overview, err
		}
	case http.StatusNotFound, http.StatusGone:
		return nil, nil
	}

	// TODO: Return the response body along with the error message

	return nil, errors.New(fmt.Sprintf("Unexpected HTTP response status code: %v", response.StatusCode))
}
