package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsHttpLink(t *testing.T) {
	assert.False(t, IsValidHTTPLink("http/example.org"))
	assert.True(t, IsValidHTTPLink("http://example.org/rss.xml"))
	assert.True(t, IsValidHTTPLink("https://localhost/rs.xml"))
	assert.False(t, IsValidHTTPLink("/tmp/a.xml"))
	assert.False(t, IsValidHTTPLink("ftp://ftp.example.org/"))
	assert.False(t, IsValidHTTPLink("!@#$%^&*()_+"))
}
