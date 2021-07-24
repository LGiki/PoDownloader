package opml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var opmlXML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
    <head>
        <title>example</title>
        <dateCreated>Sat, 24 Jul 2021 01:32:00 GMT</dateCreated>
        <dateModified>Sat, 24 Jul 2021 01:32:00 GMT</dateModified>
        <ownerName>foobar</ownerName>
        <ownerEmail>foobar@example.org</ownerEmail>
        <ownerId>foobar</ownerId>
    </head>
    <body>
        <outline text="Example Podcast" title="Example Podcast" type="rss" description="Example Podcast" xmlUrl="https://example.org/rss" htmlUrl="https://example.org" />
    </body>
</opml>`

func TestParseOPMLFromText(t *testing.T) {
	// Test correct OPML
	opml, err := ParseOPMLFromText(opmlXML)
	assert.Nil(t, err)
	assert.Equal(t, "example", opml.Head.Title)
	assert.Equal(t, "Sat, 24 Jul 2021 01:32:00 GMT", opml.Head.DateCreated)
	assert.Equal(t, "Sat, 24 Jul 2021 01:32:00 GMT", opml.Head.DateModified)
	assert.Equal(t, "foobar", opml.Head.OwnerName)
	assert.Equal(t, "foobar@example.org", opml.Head.OwnerEmail)
	assert.Equal(t, "foobar", opml.Head.OwnerId)
	assert.Equal(t, "Example Podcast", opml.Body.Outlines[0].Text)
	assert.Equal(t, "Example Podcast", opml.Body.Outlines[0].Title)
	assert.Equal(t, "rss", opml.Body.Outlines[0].Type)
	assert.Equal(t, "Example Podcast", opml.Body.Outlines[0].Description)
	assert.Equal(t, "https://example.org/rss", opml.Body.Outlines[0].XMLUrl)
	assert.Equal(t, "https://example.org", opml.Body.Outlines[0].HtmlUrl)

	// Test wrong OPML
	opml, err = ParseOPMLFromText("test")
	assert.NotNil(t, err)
	assert.Nil(t, opml)
}
