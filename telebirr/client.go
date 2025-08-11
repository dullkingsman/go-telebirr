package telebirr

type Client struct {
	config ClientConfig
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
