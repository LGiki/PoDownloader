package podcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItem_GetItemDownloadDestDir(t *testing.T) {
	item := &Item{SafeTitle: "foobar"}
	assert.Equal(t, "/tmp/foobar", item.GetItemDownloadDestDir("/tmp"))
}
