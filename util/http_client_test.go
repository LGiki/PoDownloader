package util

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	httpClient := NewHTTPClient("Test User Agent")
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "Test User Agent", request.Header.Get("User-Agent"))
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	_, _ = httpClient.Get(server.URL)
	assert.NotNil(t, httpClient)
}
