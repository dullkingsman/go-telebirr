package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Config struct {
	BaseURL                 string          `json:"baseUrl"`
	WebBaseURL              string          `json:"webBaseUrl"`
	FabricAppID             string          `json:"fabricAppId"`
	AppSecret               string          `json:"appSecret"`
	MerchantAppID           string          `json:"merchantAppId"`
	MerchantCode            string          `json:"merchantCode"`
	PrivateKey              string          `json:"privateKey"`
	ParsedPrivateKey        *rsa.PrivateKey `json:"-"`
	VerifyResponseSignature bool            `json:"verifyResponseSignature"`
	PublicKey               string          `json:"publicKey"`
	ParsedPublicKey         *rsa.PublicKey  `json:"-"`
}

func (c *Config) ParsePrivateKey() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.PrivateKey == "" {
		return fmt.Errorf("private key is empty")
	}

	block, _ := pem.Decode([]byte(c.PrivateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse key failed: %w", err)
	}

	c.ParsedPrivateKey = key.(*rsa.PrivateKey)

	return nil
}

func (c *Config) ParsePublicKey() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.PublicKey == "" {
		return fmt.Errorf("public key is empty")
	}

	block, _ := pem.Decode([]byte(c.PublicKey))

	var pubKey *rsa.PublicKey
	var err error

	switch block.Type {
	case "PUBLIC KEY":
		var parsed interface{}
		parsed, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err == nil {
			var ok bool
			pubKey, ok = parsed.(*rsa.PublicKey)
			if !ok {
				return fmt.Errorf("not an RSA public key")
			}
		}
	case "RSA PUBLIC KEY":
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return fmt.Errorf("unsupported key type: %s", block.Type)
	}

	if err != nil {
		return fmt.Errorf("parse key failed: %w", err)
	}

	c.ParsedPublicKey = pubKey

	return nil
}

func (c *Config) GetParsedPrivateKey() *rsa.PrivateKey {
	if c == nil {
		return nil
	}

	if c.ParsedPrivateKey == nil && c.PrivateKey != "" {
		_ = c.ParsePrivateKey()
	} else if c.ParsedPrivateKey == nil {
		return nil
	}

	return c.ParsedPrivateKey
}

func (c *Config) GetParsedPublicKey() *rsa.PublicKey {
	if c == nil {
		return nil
	}

	if c.ParsedPublicKey == nil && c.PublicKey != "" {
		_ = c.ParsePublicKey()
	} else if c.ParsedPublicKey == nil {
		return nil
	}

	return c.ParsedPublicKey
}
