package opml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOPML_GetAllXMLUrl(t *testing.T) {
	opml := &OPML{
		Head: nil,
		Body: &Body{
			Outlines: []*Outline{
				{
					XMLUrl: "https://example.org/rss.xml",
				},
				{
					XMLUrl: "https://example.com/feed.xml",
				},
				{
					XMLUrl: "https://foo.bar/rss",
				},
			},
		},
	}
	xmlUrls := opml.GetAllXMLUrl()
	assert.Equal(t, "https://example.org/rss.xml", xmlUrls[0])
	assert.Equal(t, "https://example.com/feed.xml", xmlUrls[1])
	assert.Equal(t, "https://foo.bar/rss", xmlUrls[2])
}
