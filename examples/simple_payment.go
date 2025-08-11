package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dullkingsman/go-telebirr/pkg/client"
	"github.com/dullkingsman/go-telebirr/telebirr"
)

var TelebirrClient = telebirr.NewClient(telebirr.ClientConfig{
	BaseURL:                 "https://196.188.120.3:38443/apiaccess/payment/gateway",
	WebBaseURL:              "https://developerportal.ethiotelebirr.et:38443/payment/web/paygate?",
	FabricAppID:             "c4182ef8-9249-458a-985e-06d191f4d505",
	AppSecret:               "fad0f06383c6297f545876694b974599",
	MerchantAppID:           "1464688410880003",
	MerchantCode:            "729748",
	VerifyResponseSignature: false,
	PrivateKey: `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC/ZcoOng1sJZ4CegopQVCw3HYqqVRLEudgT+dDpS8fRVy7zBgqZunju2VRCQuHeWs7yWgc9QGd4/8kRSLY+jlvKNeZ60yWcqEY+eKyQMmcjOz2Sn41fcVNgF+HV3DGiV4b23B6BCMjnpEFIb9d99/TsjsFSc7gCPgfl2yWDxE/Y1B2tVE6op2qd63YsMVFQGdre/CQYvFJENpQaBLMq4hHyBDgluUXlF0uA1X7UM0ZjbFC6ZIB/Hn1+pl5Ua8dKYrkVaecolmJT/s7c/+/1JeN+ja8luBoONsoODt2mTeVJHLF9Y3oh5rI+IY8HukIZJ1U6O7/JcjH3aRJTZagXUS9AgMBAAECggEBALBIBx8JcWFfEDZFwuAWeUQ7+VX3mVx/770kOuNx24HYt718D/HV0avfKETHqOfA7AQnz42EF1Yd7Rux1ZO0e3unSVRJhMO4linT1XjJ9ScMISAColWQHk3wY4va/FLPqG7N4L1w3BBtdjIc0A2zRGLNcFDBlxl/CVDHfcqD3CXdLukm/friX6TvnrbTyfAFicYgu0+UtDvfxTL3pRL3u3WTkDvnFK5YXhoazLctNOFrNiiIpCW6dJ7WRYRXuXhz7C0rENHyBtJ0zura1WD5oDbRZ8ON4v1KV4QofWiTFXJpbDgZdEeJJmFmt5HIi+Ny3P5n31WwZpRMHGeHrV23//0CgYEA+2/gYjYWOW3JgMDLX7r8fGPTo1ljkOUHuH98H/a/lE3wnnKKx+2ngRNZX4RfvNG4LLeWTz9plxR2RAqqOTbX8fj/NA/sS4mru9zvzMY1925FcX3WsWKBgKlLryl0vPScq4ejMLSCmypGz4VgLMYZqT4NYIkU2Lo1G1MiDoLy0CcCgYEAwt77exynUhM7AlyjhAA2wSINXLKsdFFF1u976x9kVhOfmbAutfMJPEQWb2WXaOJQMvMpgg2rU5aVsyEcuHsRH/2zatrxrGqLqgxaiqPz4ELINIh1iYK/hdRpr1vATHoebOv1wt8/9qxITNKtQTgQbqYci3KV1lPsOrBAB5S57nsCgYAvw+cagS/jpQmcngOEoh8I+mXgKEET64517DIGWHe4kr3dO+FFbc5eZPCbhqgxVJ3qUM4LK/7BJq/46RXBXLvVSfohR80Z5INtYuFjQ1xJLveeQcuhUxdK+95W3kdBBi8lHtVPkVsmYvekwK+ukcuaLSGZbzE4otcn47kajKHYDQKBgDbQyIbJ+ZsRw8CXVHu2H7DWJlIUBIS3s+CQ/xeVfgDkhjmSIKGX2to0AOeW+S9MseiTE/L8a1wY+MUppE2UeK26DLUbH24zjlPoI7PqCJjl0DFOzVlACSXZKV1lfsNEeriC61/EstZtgezyOkAlSCIH4fGr6tAeTU349Bnt0RtvAoGBAObgxjeH6JGpdLz1BbMj8xUHuYQkbxNeIPhH29CySn0vfhwg9VxAtIoOhvZeCfnsCRTj9OZjepCeUqDiDSoFznglrKhfeKUndHjvg+9kiae92iI6qJudPCHMNwP8wMSphkxUqnXFR3lr9A765GA980818UWZdrhrjLKtIIZdh+X1
-----END PRIVATE KEY-----`,
	PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAv2XKDp4NbCWeAnoKKUFQ
sNx2KqlUSxLnYE/nQ6UvH0Vcu8wYKmbp47tlUQkLh3lrO8loHPUBneP/JEUi2Po5
byjXmetMlnKhGPniskDJnIzs9kp+NX3FTYBfh1dwxoleG9twegQjI56RBSG/Xfff
07I7BUnO4Aj4H5dslg8RP2NQdrVROqKdqnet2LDFRUBna3vwkGLxSRDaUGgSzKuI
R8gQ4JblF5RdLgNV+1DNGY2xQumSAfx59fqZeVGvHSmK5FWnnKJZiU/7O3P/v9SX
jfo2vJbgaDjbKDg7dpk3lSRyxfWN6IeayPiGPB7pCGSdVOju/yXIx92kSU2WoF1E
vQIDAQAB
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

			var body = telebirr.MerchantPreOrderRequestBody{
				Timestamp: client.GetCurrentUnixTimestampString(),
				NonceStr:  client.CreateNonceStr(32),
				Method:    "payment.preorder",
				Version:   "1.0",
				SignType:  "SHA256WithRSA",
				BizContent: telebirr.MerchantPreOrderRequestBizContent{
					TransCurrency:       "ETB",
					TotalAmount:         "1.5",
					MerchOrderID:        client.GetCurrentUnixMilliTimestampString(),
					AppID:               _config.MerchantAppID,
					MerchCode:           _config.MerchantCode,
					TimeoutExpress:      "120m",
					TradeType:           "Checkout",
					NotifyURL:           "https://api.vps.gebeta.app/v1/external/payments/verify/post_box_subscription/sale-id/telebirr",
					Title:               "diamond_1.5",
					BusinessType:        "BuyGoods",
					PayeeIdentifier:     _config.MerchantCode,
					PayeeIdentifierType: "04",
					PayeeType:           "5000",
					RedirectURL:         "https://api.vps.gebeta.app/v1/external/payments/status/post_box_subscription/sale-id/",
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
