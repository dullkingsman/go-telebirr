package telebirr

import "time"

type Client struct {
	config                    ClientConfig
	fabricToken               *string
	fabricTokenExpirationDate *time.Time
	fabricTokenEffectiveDate  *time.Time
	tokenClearTimer           *time.Timer
}

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
		ef = *c.fabricTokenExpirationDate
		ex = *c.fabricTokenExpirationDate
	}

	if t != "" {
		return &t, &ef, &ex
	}

	return nil, nil, nil
}

func (c *Client) SetToken(token string, effectiveDate time.Time, expirationDate time.Time) {
	c.clearToken()
	c.clearTimer()

	c.fabricToken = &token
	c.fabricTokenEffectiveDate = &effectiveDate
	c.fabricTokenExpirationDate = &expirationDate

	var startTimer = func() {
		var d = time.Until(expirationDate)

		if d <= 0 {
			c.clearToken()
			return
		}

		c.tokenClearTimer = time.AfterFunc(d, func() { c.clearToken() })
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
	}
}

func (c *Client) clearToken() {
	c.fabricToken = nil
	c.fabricTokenEffectiveDate = nil
	c.fabricTokenExpirationDate = nil
}
