package util

import (
	url2 "net/url"
)

// IsValidHTTPLink returns true if specified url is a valid http link, otherwise it returns false
func IsValidHTTPLink(url string) bool {
	parsedURL, err := url2.Parse(url)
	if err != nil {
		return false
	}
	return (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != ""
}

// StripQueryParam returns a URL that does not contain any query parameters
func StripQueryParam(inURL string) string {
	u, err := url2.Parse(inURL)
	if err != nil {
		return inURL
	}
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}
