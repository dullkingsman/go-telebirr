package telebirr

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dullkingsman/go-telebirr/internal/httpclient"
)

// GetAuthTokenRequestBody represents the request payload for /payment/v1/auth/authToken
type GetAuthTokenRequestBody struct {
	Timestamp  string                 `json:"timestamp"`   // UTC timestamp, seconds
	Method     string                 `json:"method"`      // e.g. "payment.authtoken"
	NonceStr   string                 `json:"nonce_str"`   // Random string ≤ 32 chars
	SignType   string                 `json:"sign_type"`   // e.g. "SHA256WithRSA"
	Sign       string                 `json:"sign"`        // Signature ≤ 512 chars
	Version    string                 `json:"version"`     // e.g. "1.0"
	BizContent GetAuthTokenBizContent `json:"biz_content"` // Business content
}

// GetAuthTokenBizContent contains merchant and access token details
type GetAuthTokenBizContent struct {
	AppID        string `json:"appid"`         // Merchant App ID
	AccessToken  string `json:"access_token"`  // Provided from SuperApp
	TradeType    string `json:"trade_type"`    // e.g. "InApp"
	ResourceType string `json:"resource_type"` // e.g. "OpenId"
}

func (cb *GetAuthTokenRequestBody) AttachSignature(key *rsa.PrivateKey) error {
	if cb == nil {
		return fmt.Errorf("object is nil")
	}

	var signString = NewSignatureData().Add(
		*cb,
		SignatureDataExclusions{
			"sign":        true,
			"sign_type":   true,
			"biz_content": true, // handled below
		},
	).Add(
		cb.BizContent,
		SignatureDataExclusions{},
	).Construct()

	sign, err := signString.Sign(key)

	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	(*cb).Sign = sign

	return nil
}

// GetAuthTokenResponseBody represents the full response from the API
type GetAuthTokenResponseBody struct {
	Result     string                         `json:"result"` // "SUCCESS" or "FAIL"
	Code       string                         `json:"code"`   // 0 or business error code
	Msg        string                         `json:"msg"`    // Error description
	Sign       string                         `json:"sign"`   // Response signature
	NonceStr   string                         `json:"nonce_str"`
	SignType   string                         `json:"sign_type"` // "SHA256WithRSA"
	BizContent GetAuthTokenResponseBizContent `json:"biz_content"`
}

// GetAuthTokenResponseBizContent contains user identity details
type GetAuthTokenResponseBizContent struct {
	OpenID            string `json:"open_id"`
	IdentityID        string `json:"identityId"`
	IdentityType      string `json:"identityType"`
	WalletIdentityID  string `json:"walletIdentityId"`
	Identifier        string `json:"identifier,omitempty"`        // Optional
	NickName          string `json:"nickName,omitempty"`          // Optional
	Status            string `json:"status,omitempty"`            // Optional
	ShortCode         string `json:"shortcode,omitempty"`         // Optional
	WalletOrgOperator string `json:"walletOrgOperator,omitempty"` // Optional
}

func (cb *GetAuthTokenResponseBody) VerifySignature(key *rsa.PublicKey) error {
	if cb == nil {
		return fmt.Errorf("object is nil")
	}

	var signatureString = NewSignatureData().Add(
		*cb,
		SignatureDataExclusions{
			"sign":        true,
			"sign_type":   true,
			"biz_content": true,
		},
	).Add(
		cb.BizContent,
		nil,
	).Construct()

	var err = signatureString.VerifySignature(cb.Sign, key)

	if err != nil {
		return fmt.Errorf("failed to verify response signature: %w", err)
	}

	return nil
}

func (c *Client) GetAuthToken(token string, body GetAuthTokenRequestBody, config ...httpclient.ClientConfig[GetAuthTokenResponseBody]) (*httpclient.Response[GetAuthTokenResponseBody], error) {
	var err = body.AttachSignature(c.config.ParsedPrivateKey)

	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := httpclient.NewHTTPClient[GetAuthTokenResponseBody](config...).DoRequest(&httpclient.Request{
		Method: "POST",
		Url:    c.config.BaseURL + Endpoints.GetAuthToken,
		Body:   reqBody,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"X-APP-Key":     c.config.FabricAppID,
			"Authorization": token,
		},
	})

	if err != nil {
		return resp, err
	}

	if resp.Status < 200 || resp.Status >= 300 {
		log.Printf("Error response body: %+v\n", resp.Body)
		return nil, fmt.Errorf("HTTP error! status: %d", resp.Status)
	}

	return resp, nil
}
