package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// mimeTypeToExtensionName maps mime type name to file extension name
// See also:
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
// - https://www.iana.org/assignments/media-types/media-types.xhtml
var mimeTypeToExtensionName = map[string]string{
	// Audio
	"audio/aac":    "aac",
	"audio/mpeg":   "mp3",
	"audio/mp4":    "mp4",
	"audio/wav":    "wav",
	"audio/x-m4a":  "m4a",
	"audio/x-aiff": "aiff",
	"audio/ogg":    "ogg",
	// Image
	"image/bmp":  "bmp",
	"image/heic": "heic",
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
	// Text
	"text/xml": "xml",
}

// extensionNames is a string slice that contains all extensionNames appeared in mimeTypeToExtensionName
var extensionNames = []string{
	// Audio
	"aac", "mp3", "mp4", "wav", "m4a", "wav", "aiff", "ogg",
	// Image
	"bmp", "heic", "jpg", "jpeg", "png", "webp",
	// Text
	"xml",
}

// SanitizeFileName returns file name with invalid characters removed
// See also: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
func SanitizeFileName(fileName string) string {
	fileName = regexp.MustCompile(`[:/<>"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")
	return fileName
}

// GetExtensionNameByMimeType returns extension name that matches the specified mime type
func GetExtensionNameByMimeType(mimeType string) (string, bool) {
	extensionName, ok := mimeTypeToExtensionName[mimeType]
	return extensionName, ok
}

// GetRemoteFileExtensionName returns the extension name of specified URL
// Try to determine the file extension name based on the string after the last dot in the URL first,
// if can not determine the file extension name based on that, an HTTP HEAD request will be sent, then
// determine the file extension name based on the response Content-type field
func GetRemoteFileExtensionName(httpClient *http.Client, url string) (string, error) {
	urlSplitByDot := strings.Split(StripQueryParam(url), ".")
	lastSegmentOfURL := urlSplitByDot[len(urlSplitByDot)-1]
	lastSegmentOfURL = strings.ToLower(lastSegmentOfURL)
	if IsStringSliceContainText(extensionNames, lastSegmentOfURL) {
		return lastSegmentOfURL, nil
	}
	resp, err := httpClient.Head(url)
	if err != nil {
		return "", err
	}
	contentType := resp.Header.Get(http.CanonicalHeaderKey("Content-Type"))
	if extensionName, ok := GetExtensionNameByMimeType(contentType); ok {
		return extensionName, nil
	}
	return "", fmt.Errorf("unknown mimetype: %s", contentType)
}

// GetRemoteFileSize returns the file size in bytes corresponding to the specified URL
func GetRemoteFileSize(httpClient *http.Client, url string) (int64, error) {
	resp, err := httpClient.Head(url)
	if err != nil {
		return 0, err
	}
	return resp.ContentLength, nil
}

// IsPathExist returns true if specified path is exists, otherwise returns false
func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileSize returns file size in bytes
func GetFileSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

// Mkdir creates a new directory with the specified path and permission 0755
func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

// MkdirAll creates directories with the specified path and permission 0755,
// along with any necessary parents
func MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

// EnsureDir ensures the directory with the specified path is exists
func EnsureDir(path string) error {
	if !IsPathExist(path) {
		if err := Mkdir(path); err != nil {
			return err
		}
	}
	return nil
}

// EnsureDirAll ensures the directories with the specified path is exists
func EnsureDirAll(path string) error {
	if !IsPathExist(path) {
		if err := MkdirAll(path); err != nil {
			return err
		}
	}
	return nil
}

// GetFileContent returns file content in string
func GetFileContent(filePath string) (string, error) {
	if !IsPathExist(filePath) {
		return "", errors.New("file does not exist")
	}
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

// GetRSSListByTextFile returns http links from specified file path,
// the text file must contain one RSS link per line
func GetRSSListByTextFile(filePath string) ([]string, error) {
	content, err := GetFileContent(filePath)
	if err != nil {
		return nil, err
	}
	// The line break of Windows is \r\n, so I need to replace \r\n with \n
	// then I can use \n to split the file
	content = strings.ReplaceAll(content, "\r\n", "\n")
	contentSplit := strings.Split(content, "\n")
	var lines []string
	for _, line := range contentSplit {
		if line != "" && IsValidHTTPLink(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// WriteContentToFile writes specified content to specified destination file path
func WriteContentToFile(content string, destFilePath string) error {
	out, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	_, err = out.WriteString(content)
	return err
}
