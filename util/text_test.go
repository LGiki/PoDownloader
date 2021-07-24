package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFirstNCharacters(t *testing.T) {
	assert.Equal(t, "", GetFirstNCharacters("", 10))
	assert.Equal(t, "012345678", GetFirstNCharacters("0123456789", 9))
	assert.Equal(t, "0123456789", GetFirstNCharacters("0123456789", 100))
	assert.Equal(t, "0", GetFirstNCharacters("0123456789", 1))
	assert.Equal(t, "你", GetFirstNCharacters("你好世界", 1))
	assert.Equal(t, "你好世界", GetFirstNCharacters("你好世界", 50))
	assert.Equal(t, "HelloWorld你好", GetFirstNCharacters("HelloWorld你好世界", 12))
}

func TestFillTextToLength(t *testing.T) {
	assert.Equal(t, "Hello ", FillTextToLength("Hello", 6))
	assert.Equal(t, "HelloWorld", FillTextToLength("HelloWorld", 6))
	assert.Equal(t, "      ", FillTextToLength("", 6))
}
