package opml

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func parseOPMLFromBytes(bytes []byte) (*OPML, error) {
	var opml OPML
	err := xml.Unmarshal(bytes, &opml)
	if err != nil {
		return nil, err
	}
	return &opml, nil
}

func ParseOPMLFromText(text string) (*OPML, error) {
	return parseOPMLFromBytes([]byte(text))
}

func ParseOPMLFromFile(filePath string) (*OPML, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parseOPMLFromBytes(bytes)
}

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
