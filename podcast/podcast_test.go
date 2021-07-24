package podcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPodcast_GetItemCount(t *testing.T) {
	podcast := &Podcast{
		Items: []*Item{
			{
				Title: "foobar",
			},
			{
				Title: "foobar",
			},
			{
				Title: "foobar",
			},
		},
	}
	assert.Equal(t, 3, podcast.GetItemCount())
}

func TestPodcast_GetPodcastDownloadDestDir(t *testing.T) {
	podcast := &Podcast{SafeTitle: "foobar"}
	assert.Equal(t, "/tmp/foobar", podcast.GetPodcastDownloadDestDir("/tmp"))
}
