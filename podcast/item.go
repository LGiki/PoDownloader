package podcast

import (
	"path"
	"time"
)

type Item struct {
	Title       string               `json:"title,omitempty"`
	SafeTitle   string               `json:"safeTitle,omitempty"`
	Description string               `json:"description,omitempty"`
	PubDate     *time.Time           `json:"pubDate,omitempty"`
	GUID        string               `json:"guid,omitempty"`
	ITunesExt   *ITunesItemExtension `json:"iTunesExt,omitempty"`
	Enclosures  []*Enclosure         `json:"enclosures,omitempty"`
}

type ITunesItemExtension struct {
	Author   string `json:"author,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Image    string `json:"image,omitempty"`
	Duration string `json:"duration,omitempty"`
	Order    string `json:"order,omitempty"`
}

func (i *Item) GetItemDownloadDestDir(podcastDir string) string {
	return path.Join(podcastDir, i.SafeTitle)
}
