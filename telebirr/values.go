package telebirr

var Endpoints = struct {
	GenerateAppToken string
	MerchantPreOrder string
	GetAuthToken     string
}{
	GenerateAppToken: "/payment/v1/token",
	MerchantPreOrder: "/payment/v1/merchant/preOrder",
	GetAuthToken:     "/payment/v1/auth/authToken",
}
