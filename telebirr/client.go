package telebirr

import (
	"fmt"
	"time"
)

type Client struct {
	config                    ClientConfig
	fabricToken               *string
	fabricTokenExpirationDate *time.Time
	fabricTokenEffectiveDate  *time.Time
	tokenClearTimer           *time.Timer
}

var logDateLayout = "2006-01-02 15:04:05"

func NewClient(config ClientConfig) *Client {
	if err := config.ParsePrivateKey(); err != nil {
		panic(err)
	}

	if config.VerifyResponseSignature {
		if err := config.ParsePublicKey(); err != nil {
			panic(err)
		}
	}

	return &Client{
		config: config,
	}
}

func (c *Client) Config() ClientConfig {
	return c.config
}

func (c *Client) GetToken() (*string, *time.Time, *time.Time) {
	var (
		t  string
		ef time.Time
		ex time.Time
	)

	if c.fabricToken != nil {
		t = *c.fabricToken
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

func (c *Client) SetToken(token string, effectiveDate time.Time, expirationDate time.Time) {
	fmt.Printf(`set fabric token: %s
effective-date: %s
expiration-date: %s
`, token, effectiveDate.Format(logDateLayout), expirationDate.Format(logDateLayout))

	c.clearTimer()

	c.fabricToken = &token
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

func (c *Client) clearTimer() {
	if c.tokenClearTimer != nil {
		c.tokenClearTimer.Stop()
		c.tokenClearTimer = nil

		fmt.Printf("stopped fabric token clear timer for %s\n", *c.fabricToken)
	}
}

func (c *Client) clearToken(token string) {
	if c.fabricToken != nil && *c.fabricToken != token {
		return
	}

	fmt.Printf("cleared fabric token %s\n", token)

	c.fabricToken = nil
	c.fabricTokenEffectiveDate = nil
	c.fabricTokenExpirationDate = nil
}
