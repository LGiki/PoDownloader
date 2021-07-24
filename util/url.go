package util

import (
	url2 "net/url"
)

func IsValidHttpLink(str string) bool {
	url, err := url2.Parse(str)
	if err != nil {
		return false
	}
	return (url.Scheme == "http" || url.Scheme == "https") && url.Host != ""
}
