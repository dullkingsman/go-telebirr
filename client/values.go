package client

var Endpoints = struct {
	GenerateAppToken string
	MerchantPreOrder string
}{
	GenerateAppToken: "/payment/v1/token",
	MerchantPreOrder: "/payment/v1/merchant/preOrder",
}
