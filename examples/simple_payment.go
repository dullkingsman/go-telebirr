package main

import (
	"errors"
	"fmt"
	"net/http"
	"telebirr-go/client"
	"telebirr-go/internal/config"
	"telebirr-go/internal/utils"
)

var TelebirrClient = client.NewClient(config.Config{
	BaseURL:                 "https://196.188.120.3:38443/apiaccess/payment/gateway",
	WebBaseURL:              "https://developerportal.ethiotelebirr.et:38443/payment/web/paygate?",
	FabricAppID:             "<fabric-id>",
	AppSecret:               "<app-secret>",
	MerchantAppID:           "<merchant-app-id>",
	MerchantCode:            "<merchant-code>",
	VerifyResponseSignature: false,
	PrivateKey: `-----BEGIN PRIVATE KEY-----
<private-key>
-----END PRIVATE KEY-----`,
	PublicKey: `-----BEGIN PUBLIC KEY-----
<public-key>
-----END PUBLIC KEY-----`,
})

func main() {
	_ = (&http.Server{
		Addr: ":8083",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var resp, err = TelebirrClient.GenerateAppToken()

			if resp != nil {
				fmt.Printf("resp: %+v\n", *resp)
			}

			if err != nil {
				fmt.Printf("err: %+v\n", err)
			}

			var _config = TelebirrClient.Config()

			var body = client.MerchantPreOrderRequestBody{
				Timestamp: utils.GetCurrentUnixTimestampString(),
				NonceStr:  utils.CreateNonceStr(32),
				Method:    "payment.preorder",
				Version:   "1.0",
				SignType:  "SHA256WithRSA",
				BizContent: client.MerchantPreOrderRequestBizContent{
					TransCurrency:       "ETB",
					TotalAmount:         "1.5",
					MerchOrderID:        utils.GetCurrentUnixMilliTimestampString(),
					AppID:               _config.MerchantAppID,
					MerchCode:           _config.MerchantCode,
					TimeoutExpress:      "120m",
					TradeType:           "Checkout",
					NotifyURL:           "<your-callback-url>",
					Title:               "diamond_1.5",
					BusinessType:        "BuyGoods",
					PayeeIdentifier:     _config.MerchantCode,
					PayeeIdentifierType: "04",
					PayeeType:           "5000",
					RedirectURL:         "<your-redirect-url>",
					CallbackInfo:        "From web",
				},
			}

			if resp == nil {
				fmt.Printf("err: %+v\n", errors.New("token is nil"))
				return
			}

			orderResp, err := TelebirrClient.MerchantPreOrder(resp.Body.Token, body)

			if orderResp != nil {
				fmt.Printf("resp: %+v\n", *orderResp)
			}

			if err != nil {
				fmt.Printf("err: %+v\n", err)
			}

			if orderResp != nil && orderResp.Body.BizContent.PrepayID != "" {
				request, err := TelebirrClient.NewRawRequest(orderResp.Body.BizContent.PrepayID)

				if err != nil {
					fmt.Printf("err: %+v\n", err)
				}

				var requestStr = request.Extend(map[string]string{
					"version":    "1.0",
					"trade_type": "Checkout",
				}).String()

				fmt.Printf("redirect: %+v\n", requestStr)

				http.Redirect(w, r, requestStr, http.StatusFound)
			}
		}),
	}).ListenAndServe()
}
