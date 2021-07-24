package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStringSliceContainText(t *testing.T) {
	s := []string{"foo", "bar", "hello", "world"}
	assert.Equal(t, true, IsStringSliceContainText(s, "foo"))
	assert.Equal(t, true, IsStringSliceContainText(s, "bar"))
	assert.Equal(t, true, IsStringSliceContainText(s, "hello"))
	assert.Equal(t, true, IsStringSliceContainText(s, "world"))
	assert.Equal(t, false, IsStringSliceContainText(s, "example"))
	assert.Equal(t, false, IsStringSliceContainText(s, ""))
}
