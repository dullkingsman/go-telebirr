package telebirr

import (
	"context"
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

func (c *Client) GetToken(ctx context.Context, key ...string) (*string, *time.Time, *time.Time, error) {
	return c.fabricTokenCache.GetToken(ctx, key...)
}

func (c *Client) SetToken(ctx context.Context, token string, effectiveDate time.Time, expirationDate time.Time) {
	_ = c.fabricTokenCache.SetToken(ctx, token, effectiveDate, expirationDate)
}
