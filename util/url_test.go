package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsHttpLink(t *testing.T) {
	assert.False(t, IsValidHttpLink("http/example.org"))
	assert.True(t, IsValidHttpLink("http://example.org/rss.xml"))
	assert.True(t, IsValidHttpLink("https://localhost/rs.xml"))
	assert.False(t, IsValidHttpLink("/tmp/a.xml"))
	assert.False(t, IsValidHttpLink("ftp://ftp.example.org/"))
	assert.False(t, IsValidHttpLink("!@#$%^&*()_+"))
}
