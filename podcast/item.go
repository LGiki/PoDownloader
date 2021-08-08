package podcast

import (
	"encoding/json"
	"path"
	"time"
)

// Item is the item (episode) of Podcast
type Item struct {
	Title       string               `json:"title,omitempty"`
	SafeTitle   string               `json:"safeTitle,omitempty"`
	Description string               `json:"description,omitempty"`
	PubDate     *time.Time           `json:"pubDate,omitempty"`
	GUID        string               `json:"guid,omitempty"`
	ITunesExt   *ITunesItemExtension `json:"iTunesExt,omitempty"`
	Enclosures  []*Enclosure         `json:"enclosures,omitempty"`
}

// ITunesItemExtension is the extension fields of Podcast items
type ITunesItemExtension struct {
	Author   string `json:"author,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Image    string `json:"image,omitempty"`
	Duration string `json:"duration,omitempty"`
	Order    string `json:"order,omitempty"`
}

// GetItemDownloadDestDir returns item download destination dir
// item download destination dir = Podcast download destination dir + episode title
func (i *Item) GetItemDownloadDestDir(podcastDir string) string {
	return path.Join(podcastDir, i.SafeTitle)
}

// GetJSON returns an Item instance JSON format
func (i *Item) GetJSON() (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
