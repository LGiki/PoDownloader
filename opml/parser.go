package opml

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// parseOPMLFromBytes parses OPML from bytes slice
// and returns an OPML instance
func parseOPMLFromBytes(bytes []byte) (*OPML, error) {
	var opml OPML
	err := xml.Unmarshal(bytes, &opml)
	if err != nil {
		return nil, err
	}
	return &opml, nil
}

// ParseOPMLFromText parses OPML from text
// and returns an OPML instance
func ParseOPMLFromText(text string) (*OPML, error) {
	return parseOPMLFromBytes([]byte(text))
}

// ParseOPMLFromFile parses OPML from specified file path
// and returns an OPML instance
func ParseOPMLFromFile(filePath string) (*OPML, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parseOPMLFromBytes(bytes)
}

// ParseOPMLFromURL parses OPML from specified URL using specified http client
// and returns an OPML instance
func ParseOPMLFromURL(httpClient *http.Client, url string) (*OPML, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseOPMLFromBytes(bytes)
}
