package client

import (
	"crypto/rsa"
	"fmt"
	"telebirr-go/internal/model"
)

type PaymentCallbackRequestBody struct {
	NotifyURL      string `json:"notify_url,omitempty"`       // Example: "http://197.156.68.29:5050/v2/api/order-v2/mini/payment" — Callback URL for payment notification
	AppID          string `json:"appid,omitempty"`            // Example: "853694808089634" — Application ID assigned by Mobile Payment system
	NotifyTime     string `json:"notify_time,omitempty"`      // Example: "1670575472482" — Timestamp in UTC seconds (as string, <= 13 chars)
	MerchCode      string `json:"merch_code,omitempty"`       // Example: "245445" — Merchant short code registered with Mobile Money
	MerchOrderID   string `json:"merch_order_id,omitempty"`   // Example: "1670575560882" — Order ID from merchant system
	PaymentOrderID string `json:"payment_order_id,omitempty"` // Example: "00801104C911443200001002" — Order ID assigned by payment system
	TotalAmount    string `json:"total_amount,omitempty"`     // Example: "10.00" — Transaction amount, up to 2 decimal places
	TransID        string `json:"trans_id,omitempty"`         // Example: "49485948475845" — Transaction ID (appears in sample but not schema description)
	TransCurrency  string `json:"trans_currency,omitempty"`   // Example: "ETB" — Currency code (e.g., ETB for Ethiopian Birr)
	TradeStatus    string `json:"trade_status,omitempty"`     // Example: "Completed" — Transaction status (e.g., Paying, Expired, Pending, Completed, Failure)
	TransEndTime   string `json:"trans_end_time,omitempty"`   // Example: "1670575472000" — End time of transaction in UTC seconds (as string)
	CallbackInfo   string `json:"callback_info,omitempty"`    // Example: "your-custom-callback-info" — Merchant-specific data passed back in callback (optional)
	Sign           string `json:"sign,omitempty"`             // Example: "AOwWQF0QDg0jzzs5otLYOunoR65GGgC3hyr+oYn8mm1Qph6Een7C…" — Response signature string
	SignType       string `json:"sign_type,omitempty"`        // Example: "SHA256WithRSA" — Signature algorithm used
}

func (cb *PaymentCallbackRequestBody) VerifySignature(key *rsa.PublicKey) error {
	if cb == nil {
		return fmt.Errorf("object is nil")
	}

	var signatureString = model.NewSignatureData().Add(
		*cb,
		model.SignatureDataExclusions{
			"sign":      true,
			"sign_type": true,
		},
	).Construct()

	var err = signatureString.VerifySignature(cb.Sign, key)

	if err != nil {
		return fmt.Errorf("failed to verify response signature: %w", err)
	}

	return nil
}
