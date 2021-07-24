package podcast

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewPodcastParser(t *testing.T) {
	podcastParser := NewPodcastParser(&http.Client{})
	assert.NotNil(t, podcastParser)
}
