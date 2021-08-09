package util

import (
	"golang.org/x/net/http/httpproxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

// NewHTTPClient initializes and returns a http client that will send http requests using specified user agent
func NewHTTPClient(userAgent string) *http.Client {
	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (uri *url.URL, err error) {
				r.Header.Set("User-Agent", userAgent)
				return proxyFunc(r.URL)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}
}
