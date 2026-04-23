package telebirr

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dullkingsman/go-telebirr/core/httpclient"
)

// QueryOrderRequestBody represents the request payload for /payment/v1/auth/authToken
type QueryOrderRequestBody struct {
	Timestamp  string               `json:"timestamp"`   // UTC timestamp, seconds
	Method     string               `json:"method"`      //"payment.queryorder"
	NonceStr   string               `json:"nonce_str"`   // Random string ≤ 32 chars
	SignType   string               `json:"sign_type"`   // e.g. "SHA256WithRSA"
	Sign       string               `json:"sign"`        // Signature ≤ 512 chars
	AppCode    string               `json:"app_code"`    // The merchant's unique code
	Version    string               `json:"version"`     // e.g. "1.0"
	BizContent QueryOrderBizContent `json:"biz_content"` // Business content
}

// QueryOrderBizContent contains merchant and access token details
type QueryOrderBizContent struct {
	AppID        string `json:"appid"`          // Merchant App ID
	MerchCode    string `json:"merch_code"`     // Short code
	MerchOrderID string `json:"merch_order_id"` // trx_ref
}

func (cb *QueryOrderRequestBody) AttachSignature(key *rsa.PrivateKey) error {
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

// QueryOrderResponseBody represents the full response from the API
type QueryOrderResponseBody struct {
	Result     string                       `json:"result"` // "SUCCESS" or "FAIL"
	Code       string                       `json:"code"`   // 0 or business error code
	Msg        string                       `json:"msg"`    // Error description
	Sign       string                       `json:"sign"`   // Response signature
	NonceStr   string                       `json:"nonce_str"`
	SignType   string                       `json:"sign_type"` // "SHA256WithRSA"
	BizContent QueryOrderResponseBizContent `json:"biz_content"`
}

// QueryOrderResponseBizContent contains user identity details
type QueryOrderResponseBizContent struct {
	MerchOrderId   string `json:"merch_order_id"`
	OrderStatus    string `json:"order_status"`
	TradeStatus    string `json:"trade_status"`
	PaymentOrderId string `json:"payment_order_id"`
	TransTime      string `json:"trans_time"`
	TransCurrency  string `json:"trans_currency"`
	TotalAmount    string `json:"total_amount"`
}

func (cb *QueryOrderResponseBody) VerifySignature(key *rsa.PublicKey) error {
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

func (c *Client) QueryOrder(token string, body QueryOrderRequestBody, config ...httpclient.ClientConfig[QueryOrderResponseBody]) (*httpclient.Response[QueryOrderResponseBody], error) {
	if body.Method == "" {
		body.Method = "payment.queryorder"
	}

	var err = body.AttachSignature(c.config.ParsedPrivateKey)

	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := httpclient.NewHTTPClient[QueryOrderResponseBody](config...).DoRequest(&httpclient.Request{
		Method: "POST",
		Url:    c.config.BaseURL + Endpoints.QueryOrder,
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
