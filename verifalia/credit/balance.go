package credit

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
	"fmt"
	"github.com/ericlagergren/decimal"
	"github.com/verifalia/verifalia-go-sdk/verifalia/rest"
	"io/ioutil"
	"net/http"
)

type Client struct {
	RestClient rest.Client
}

type Balance struct {
	CreditPacks        decimal.Big  `json:"creditPacks"`
	FreeCredits        *decimal.Big `json:"freeCredits"`
	FreeCreditsResetIn *string      `json:"freeCreditsResetIn"`
}

func (client *Client) GetBalance() (*Balance, error) {
	return client.getBalance(nil)
}

func (client *Client) GetBalanceWithContext(ctx context.Context) (*Balance, error) {
	return client.getBalance(ctx)
}

func (client *Client) getBalance(ctx context.Context) (*Balance, error) {
	response, err := client.RestClient.Invoke(rest.InvocationOptions{
		Method:   http.MethodGet,
		Resource: "credits/balance",
		Context:  ctx,
	})

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		var balance Balance

		if err := json.Unmarshal(responseData, &balance); err != nil {
			return nil, err
		}

		return &balance, nil
	}

	return nil, fmt.Errorf("unexpected HTTP response: %d", response.StatusCode)
}
