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

type FabricTokenCache interface {
	GetToken() (*string, *time.Time, *time.Time)
	SetToken(token string, effectiveDate time.Time, expirationDate time.Time)
}

type DefaultFabricTokenCache struct {
	token                     *string
	fabricTokenExpirationDate *time.Time
	fabricTokenEffectiveDate  *time.Time
	tokenClearTimer           *time.Timer
}

func (c *DefaultFabricTokenCache) GetToken() (*string, *time.Time, *time.Time) {
	var (
		t  string
		ef time.Time
		ex time.Time
	)

	if c.token != nil {
		t = *c.token
		ef = *c.fabricTokenEffectiveDate
		ex = *c.fabricTokenExpirationDate
	}

	if t != "" {
		fmt.Printf(
			`using cached fabric token: %s
effective-date: %s 
expiration-date: %s
`,
			t,
			ef.Format(logDateLayout),
			ex.Format(logDateLayout),
		)

		return &t, &ef, &ex
	}

	fmt.Println("No cached fabric token found")

	return nil, nil, nil
}

func (c *DefaultFabricTokenCache) SetToken(token string, effectiveDate time.Time, expirationDate time.Time) {
	fmt.Printf(`set fabric token: %s
effective-date: %s
expiration-date: %s
`, token, effectiveDate.Format(logDateLayout), expirationDate.Format(logDateLayout))

	c.clearTimer()

	c.token = &token
	c.fabricTokenEffectiveDate = &effectiveDate
	c.fabricTokenExpirationDate = &expirationDate

	var startTimer = func() {
		var d = time.Until(expirationDate)

		if d <= 0 {
			c.clearToken(token)
			return
		}

		c.tokenClearTimer = time.AfterFunc(d, func() { c.clearToken(token) })

		fmt.Printf("started fabric token clear timer for %s\n", token)
	}

	var d = time.Until(effectiveDate)

	if d <= 0 {
		startTimer()
		return
	}

	time.AfterFunc(d, func() { startTimer() })
}

func (c *DefaultFabricTokenCache) clearTimer() {
	if c.tokenClearTimer != nil {
		c.tokenClearTimer.Stop()
		c.tokenClearTimer = nil

		fmt.Printf("stopped fabric token clear timer for %s\n", *c.token)
	}
}

func (c *DefaultFabricTokenCache) clearToken(token string) {
	if c.token != nil && *c.token != token {
		return
	}

	fmt.Printf("cleared fabric token %s\n", token)

	c.token = nil
	c.fabricTokenEffectiveDate = nil
	c.fabricTokenExpirationDate = nil
}
