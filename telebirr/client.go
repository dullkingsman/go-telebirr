package telebirr

import (
	"fmt"
	"time"
)

var logDateLayout = "2006-01-02 15:04:05"

type Client struct {
	config           ClientConfig
	fabricTokenCache FabricTokenCache
}

func NewClient(config ClientConfig, fabricTokenCache ...FabricTokenCache) *Client {
	if err := config.ParsePrivateKey(); err != nil {
		panic(err)
	}

	if config.VerifyResponseSignature {
		if err := config.ParsePublicKey(); err != nil {
			panic(err)
		}
	}

	var tmp = &Client{
		config: config,
	}

	if len(fabricTokenCache) > 0 {
		fmt.Printf("using custom fabric token cache...\n")
		tmp.fabricTokenCache = fabricTokenCache[0]
	} else {
		fmt.Printf("using default fabric token cache...\n")
		tmp.fabricTokenCache = &DefaultFabricTokenCache{}
	}

	return tmp
}

func (c *Client) Config() ClientConfig {
	return c.config
}

func (c *Client) GetToken() (*string, *time.Time, *time.Time) {
	return c.fabricTokenCache.GetToken()
}

func (c *Client) SetToken(token string, effectiveDate time.Time, expirationDate time.Time) {
	c.fabricTokenCache.SetToken(token, effectiveDate, expirationDate)
}
