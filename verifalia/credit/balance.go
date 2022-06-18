package credit

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
