package client

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dullkingsman/go-telebirr/internal/httpclient"
	"github.com/dullkingsman/go-telebirr/internal/model"
	"github.com/dullkingsman/go-telebirr/internal/utils"
	"log"
	"net/url"
)

type MerchantPreOrderRequestBody struct {
	NonceStr   string                            `json:"nonce_str"`
	BizContent MerchantPreOrderRequestBizContent `json:"biz_content"`
	Method     string                            `json:"method"`
	Sign       string                            `json:"sign"`
	Version    string                            `json:"version"`
	SignType   string                            `json:"sign_type"`
	Timestamp  string                            `json:"timestamp"`
}

type MerchantPreOrderRequestBizContent struct {
	TransCurrency       string `json:"trans_currency"`
	TotalAmount         string `json:"total_amount"`
	MerchOrderID        string `json:"merch_order_id"`
	AppID               string `json:"appid"`
	MerchCode           string `json:"merch_code"`
	TimeoutExpress      string `json:"timeout_express"`
	TradeType           string `json:"trade_type"`
	NotifyURL           string `json:"notify_url"`
	Title               string `json:"title"`
	BusinessType        string `json:"business_type"`
	PayeeIdentifier     string `json:"payee_identifier"`
	PayeeIdentifierType string `json:"payee_identifier_type"`
	PayeeType           string `json:"payee_type"`
	RedirectURL         string `json:"redirect_url"`
	CallbackInfo        string `json:"callback_info"`
}

func (cb *MerchantPreOrderRequestBody) AttachSignature(key *rsa.PrivateKey) error {
	if cb == nil {
		return fmt.Errorf("object is nil")
	}

	var signString = model.NewSignatureData().Add(
		*cb,
		model.SignatureDataExclusions{
			"sign":        true,
			"sign_type":   true,
			"biz_content": true, // handled below
		},
	).Add(
		cb.BizContent,
		model.SignatureDataExclusions{},
	).Construct()

	sign, err := signString.Sign(key)

	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	(*cb).Sign = sign

	return nil
}

type MerchantPreOrderResponseBody struct {
	NonceStr   string                             `json:"nonce_str"`
	BizContent MerchantPreOrderResponseBizContent `json:"biz_content"`
	Code       string                             `json:"code"`
	Msg        string                             `json:"msg"`
	Result     string                             `json:"result"`
	Sign       string                             `json:"sign"`
	SignType   string                             `json:"sign_type"`
}

type MerchantPreOrderResponseBizContent struct {
	MerchOrderID string `json:"merch_order_id"`
	PrepayID     string `json:"prepay_id"`
}

func (cb *MerchantPreOrderResponseBody) VerifySignature(key *rsa.PublicKey) error {
	if cb == nil {
		return fmt.Errorf("object is nil")
	}

	var signatureString = model.NewSignatureData().Add(
		*cb,
		model.SignatureDataExclusions{
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

func (c *Client) MerchantPreOrder(token string, body MerchantPreOrderRequestBody, config ...httpclient.ClientConfig[MerchantPreOrderResponseBody]) (*httpclient.Response[MerchantPreOrderResponseBody], error) {
	var err = body.AttachSignature(c.config.ParsedPrivateKey)

	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := httpclient.NewHTTPClient[MerchantPreOrderResponseBody](config...).DoRequest(&httpclient.Request{
		Method: "POST",
		Url:    c.config.BaseURL + Endpoints.MerchantPreOrder,
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

	if c.config.VerifyResponseSignature {
		err = resp.Body.VerifySignature(c.config.ParsedPublicKey)

		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

type RawRequest struct {
	AppId     string `json:"appid"`
	MerchCode string `json:"merch_code"`
	NonceStr  string `json:"nonce_str"`
	PrepayId  string `json:"prepay_id"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	SignType  string `json:"sign_type"`
}

func (c *Client) NewRawRequest(prepayId string) (RawRequestString, error) {
	var req = RawRequest{
		AppId:     c.config.MerchantAppID,
		MerchCode: c.config.MerchantCode,
		NonceStr:  utils.CreateNonceStr(32),
		PrepayId:  prepayId,
		Timestamp: utils.GetCurrentUnixTimestampString(),
		SignType:  "SHA256WithRSA",
	}

	var signString = model.NewSignatureData().Add(
		req,
		map[string]bool{
			"sign":      true,
			"sign_type": true,
		},
	).Construct()

	var sign, err = signString.Sign(c.config.ParsedPrivateKey)

	if err != nil {
		return "", errors.New("failed to sign request: " + err.Error())
	}

	req.Sign = sign

	var request = c.config.WebBaseURL

	request += "appid=" + url.QueryEscape(req.AppId) + "&"
	request += "merch_code=" + url.QueryEscape(req.MerchCode) + "&"
	request += "nonce_str=" + url.QueryEscape(req.NonceStr) + "&"
	request += "prepay_id=" + url.QueryEscape(req.PrepayId) + "&"
	request += "timestamp=" + url.QueryEscape(req.Timestamp) + "&"
	request += "sign=" + url.QueryEscape(req.Sign) + "&"
	request += "sign_type=" + url.QueryEscape(req.SignType)

	return RawRequestString(request), nil
}

type RawRequestString string

func (r RawRequestString) String() string {
	return string(r)
}

func (r RawRequestString) StringPrt() *string {
	var tmp = string(r)
	return &tmp
}

func (r RawRequestString) Extend(extension ...map[string]string) RawRequestString {
	var _extension = ""

	if len(extension) > 0 && extension[0] != nil {
		for k, v := range extension[0] {
			if _extension != "" {
				_extension += "&"
			}

			_extension += k + "=" + url.QueryEscape(v)
		}
	}

	return r + "&" + RawRequestString(_extension)
}
