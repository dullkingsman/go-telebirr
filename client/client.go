package client

import (
	"telebirr-go/internal/config"
)

type Client struct {
	config config.Config
}

func NewClient(config config.Config) *Client {
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

func (c *Client) Config() config.Config {
	return c.config
}
