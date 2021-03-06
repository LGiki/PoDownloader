package podcast

import (
	"PoDownloader/util"
	"encoding/json"
	"net/http"
)

// Enclosure is the enclosure item of Podcast
type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

// GetEnclosureFileExtensionName determine file extension name by enclosure type first
// When can not determine file extension name by enclosure type,
// will call util.GetRemoteFileExtensionName to get file extension name
// by sending a HTTP HEAD request
func (e *Enclosure) GetEnclosureFileExtensionName(httpClient *http.Client) (string, error) {
	if extensionName, ok := util.GetExtensionNameByMimeType(e.Type); ok {
		return extensionName, nil
	}
	extensionName, err := util.GetRemoteFileExtensionName(httpClient, e.URL)
	if err != nil {
		return "", err
	}
	return extensionName, nil
}

// GetJSON returns an Enclosure instance in JSON format
func (e *Enclosure) GetJSON() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
