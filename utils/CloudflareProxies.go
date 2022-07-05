package utils

import (
	"net/http"
	"strings"
)

func GetCloudflareProxyIPs() ([]string, error) {

	// Get Proxy IPs in IPv4 Format
	res, err := http.Get("https://www.cloudflare.com/ips-v4")
	if err != nil {
		return nil, err
	}

	// Ready Request Body
	var responseBody []byte
	if _, err := res.Body.Read(responseBody); err != nil {
		return nil, err
	}

	// Split String by \n
	proxyIPs := strings.Split(string(responseBody), "\n")
	return proxyIPs, nil
}
