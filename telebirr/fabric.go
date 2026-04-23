package telebirr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dullkingsman/go-telebirr/core/httpclient"
)

type GenerateAppTokenRequestBody struct {
	AppSecret string `json:"appSecret"`
}

type GenerateAppTokenResponseBody struct {
	Token          string `json:"token"`
	EffectiveDate  string `json:"effectiveDate"`
	ExpirationDate string `json:"expirationDate"`
}

func (c *Client) GenerateAppToken(config ...httpclient.ClientConfig[GenerateAppTokenResponseBody]) (*httpclient.Response[GenerateAppTokenResponseBody], error) {
	var dateLayout = "20060102150405"

	var token, tokenEffectiveDate, tokenExpirationDate = c.GetToken()

	if token != nil {
		return &httpclient.Response[GenerateAppTokenResponseBody]{
			Status: http.StatusOK,
			Body: GenerateAppTokenResponseBody{
				Token:          *token,
				EffectiveDate:  tokenEffectiveDate.Format(dateLayout),
				ExpirationDate: tokenExpirationDate.Format(dateLayout),
			},
		}, nil
	}

	var reqBody, err = json.Marshal(GenerateAppTokenRequestBody{AppSecret: c.config.AppSecret})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := httpclient.NewHTTPClient[GenerateAppTokenResponseBody](config...).DoRequest(&httpclient.Request{
		Method: "POST",
		Url:    c.config.BaseURL + Endpoints.GenerateAppToken,
		Body:   reqBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-APP-Key":    c.config.FabricAppID,
		},
	})

	if err != nil {
		return resp, err
	}

	if resp.Status < 200 || resp.Status >= 300 {
		log.Printf("Error response body: %+v\n", resp.Body)
		return nil, fmt.Errorf("HTTP error! status: %d", resp.Status)
	}

	var (
		expirationDate, _ = time.Parse(dateLayout, resp.Body.ExpirationDate)
		effectiveDate, _  = time.Parse(dateLayout, resp.Body.EffectiveDate)
	)

	c.SetToken(resp.Body.Token, effectiveDate, expirationDate)

	return resp, nil
}
