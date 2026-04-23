package values

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

func GetDefaultHttpClient(serverURL ...*string) *http.Client {
	var tmp = &http.Client{
		Timeout: 30 * time.Second,
	}

	var transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if len(serverURL) > 0 && serverURL[0] != nil {
		var parsedURL, _ = url.Parse(*serverURL[0])

		if parsedURL != nil {
			transport.TLSClientConfig.ServerName = parsedURL.Hostname()
		}
	}

	tmp.Transport = transport

	return tmp
}
