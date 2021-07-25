package util

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSanitizeFileName(t *testing.T) {
	sanitizeResult := SanitizeFileName("<>:\"/\\|?*")
	assert.Equal(t, "", sanitizeResult)
}

func TestGetExtensionNameByMimeType(t *testing.T) {
	extensionName, ok := GetExtensionNameByMimeType("image/jpeg")
	assert.True(t, ok)
	assert.Equal(t, "jpg", extensionName)
}

func TestGetRemoteFileExtensionName(t *testing.T) {
	jpegServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/jpeg")
	}))
	defer jpegServer.Close()
	extensionName, err := GetRemoteFileExtensionName(&http.Client{}, "https://example.org/foobar.jpg")
	assert.Nil(t, err)
	assert.Equal(t, "jpg", extensionName)
	extensionName, err = GetRemoteFileExtensionName(&http.Client{}, jpegServer.URL)
	assert.Nil(t, err)
	assert.Equal(t, "jpg", extensionName)
	htmlServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	}))
	defer htmlServer.Close()
	extensionName, err = GetRemoteFileExtensionName(&http.Client{}, htmlServer.URL)
	assert.NotNil(t, err)
	assert.Equal(t, "", extensionName)
	extensionName, err = GetRemoteFileExtensionName(&http.Client{}, "!@#$%^&*()")
	assert.NotNil(t, err)
	assert.Equal(t, "", extensionName)
}

func TestGetRemoteFileSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("HelloWorld"))
	}))
	defer server.Close()
	fileSize, err := GetRemoteFileSize(&http.Client{}, server.URL)
	assert.Nil(t, err)
	assert.Equal(t, int64(len([]byte("HelloWorld"))), fileSize)
	fileSize, err = GetRemoteFileSize(&http.Client{}, "!@#$%^&*()")
	assert.NotNil(t, err)
}
