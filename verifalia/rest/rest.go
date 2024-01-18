package rest

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
	"errors"
	"fmt"
	"github.com/verifalia/verifalia-go-sdk/verifalia/auth"
	"io"
	"net/http"
	"net/url"
)

// BaseUrls contain the standard base URLs for the Verifalia API.
var BaseUrls = []string{
	"https://api-1.verifalia.com/v2.5",
	"https://api-2.verifalia.com/v2.5",
	"https://api-3.verifalia.com/v2.5",
}

// BaseCcaUrls contain the client-certificate authentication base URLs for the Verifalia API.
var BaseCcaUrls = []string{
	"https://api-cca-1.verifalia.com/v2.5",
	"https://api-cca-2.verifalia.com/v2.5",
	"https://api-cca-3.verifalia.com/v2.5",
}

// ContentType is an enum-like struct which contains the MIME content types understood by the Verifalia API.
var ContentType = struct {
	// application/json MIME content type.
	ApplicationJson string
	// Plain-text files (.txt), with one emailValidation address per line.
	TextPlain string
	// Comma-separated values (.csv).
	TextCsv string
	// Tab-separated values (usually coming with the .tsv extension).
	TextTsv string
	// Microsoft Excel 97-2003 Worksheet (.xls).
	ExcelXls string
	// Microsoft Excel workbook (.xslx).
	ExcelXlsx string
}{
	ApplicationJson: "application/json",
	TextPlain:       "text/plain",
	TextCsv:         "text/csv",
	TextTsv:         "text/tab-separated-values",
	ExcelXls:        "application/vnd.ms-excel",
	ExcelXlsx:       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

type multiplexedRestClient struct {
	userAgent              string
	underlyingClient       *http.Client
	authenticationProvider auth.Provider
	currentBaseUrlIdx      int
	baseUrls               []string
}

type InvocationOptions struct {
	Method      string
	Resource    string
	QueryParams url.Values
	Body        io.Reader
	Context     context.Context
	Headers     map[string]string
}

type Client interface {
	Invoke(options InvocationOptions) (*http.Response, error)
}

func NewMultiplexedRestClient(authenticationProvider auth.Provider, userAgent string, baseUrls []string) *multiplexedRestClient {
	httpClient := authenticationProvider.BuildClient()

	return &multiplexedRestClient{
		userAgent:              userAgent,
		underlyingClient:       httpClient,
		authenticationProvider: authenticationProvider,
		currentBaseUrlIdx:      0,
		baseUrls:               baseUrls,
	}
}

type invocationError struct {
	url   string
	error error
}

func (client *multiplexedRestClient) Invoke(options InvocationOptions) (*http.Response, error) {
	errs := make([]invocationError, 0)

	// Performs a maximum of as many attempts as the number of configured base API endpoints, keeping track
	// of the last used endpoint after each call, in order to try to distribute the load evenly across the
	// available endpoints.

	for idxAttempt := 0; idxAttempt < len(client.baseUrls); idxAttempt++ {
		// Retrieve the API base URL

		baseUrl := client.baseUrls[client.currentBaseUrlIdx%len(client.baseUrls)]

		// The next request will be performed on a subsequent API endpoint

		client.currentBaseUrlIdx++

		// Set up the query string

		queryString := ""

		if options.QueryParams != nil && len(options.QueryParams) > 0 {
			queryString = options.QueryParams.Encode()
		}

		// Build the final URL

		finalUrl := fmt.Sprintf("%s/%s?%s", baseUrl, options.Resource, queryString)

		// log.Printf("%v %v (attempt %d)...\n", options.Method, finalUrl, idxAttempt)

		// Init the HTTP request

		var request *http.Request
		var err error

		if options.Context == nil {
			request, err = http.NewRequest(options.Method, finalUrl, options.Body)
		} else {
			request, err = http.NewRequestWithContext(options.Context, options.Method, finalUrl, options.Body)
		}

		if err != nil {
			errs = append(errs, invocationError{
				url:   finalUrl,
				error: err,
			})

			continue
		}

		// Default headers

		request.Header.Set("User-Agent", client.userAgent)
		request.Header.Set("Content-Type", ContentType.ApplicationJson)
		request.Header.Set("Accept", ContentType.ApplicationJson)

		// Custom headers

		if options.Headers != nil {
			for k, v := range options.Headers {
				request.Header.Set(k, v)
			}
		}

		// Authenticate the underlying client, if needed

		err = client.authenticationProvider.Authenticate(request)

		if err != nil {
			errs = append(errs, invocationError{
				url:   finalUrl,
				error: err,
			})

			continue
		}

		// Send the request to the Verifalia servers

		response, err := client.underlyingClient.Do(request)

		if err != nil {
			errs = append(errs, invocationError{
				url:   finalUrl,
				error: err,
			})

			continue
		}

		// Fails on the first occurrence of an HTTP { 401, 403 } status codes

		if response.StatusCode == 401 || response.StatusCode == 403 {
			return nil, fmt.Errorf("can't authenticate to Verifalia using the provided credential (HTTP status code: %d)", response.StatusCode)
		}

		return response, nil
	}

	// Generate an error out of the potentially multiple invocation errors

	finalErrorMessage := "All the base URIs are unreachable."

	for _, err := range errs {
		finalErrorMessage = fmt.Sprintf("%v\n%v => %v", finalErrorMessage, err.url, err.error)
	}

	return nil, errors.New(finalErrorMessage)
}
