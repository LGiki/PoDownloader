package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

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

var extensionNames = []string{
	// Audio
	"aac", "mp3", "mp4", "wav", "m4a", "wav", "aiff", "ogg",
	// Image
	"bmp", "heic", "jpg", "jpeg", "png", "webp",
	// Text
	"xml",
}

// SanitizeFileName
// Remove invalid characters in filename
// See also: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
func SanitizeFileName(fileName string) string {
	fileName = regexp.MustCompile(`[:/<>"\\|?*]`).ReplaceAllString(fileName, "")
	fileName = regexp.MustCompile(`\s+`).ReplaceAllString(fileName, " ")
	return fileName
}

func GetExtensionNameByMimeType(mimeType string) (string, bool) {
	extensionName, ok := mimeTypeToExtensionName[mimeType]
	return extensionName, ok
}

func stripQueryParam(inUrl string) string {
	u, err := url.Parse(inUrl)
	if err != nil {
		return inUrl
	}
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func GetRemoteFileExtensionName(httpClient *http.Client, url string) (string, error) {
	urlSplitByDot := strings.Split(stripQueryParam(url), ".")
	lastSegmentOfUrl := urlSplitByDot[len(urlSplitByDot)-1]
	lastSegmentOfUrl = strings.ToLower(lastSegmentOfUrl)
	if IsStringSliceContainText(extensionNames, lastSegmentOfUrl) {
		return lastSegmentOfUrl, nil
	}
	resp, err := httpClient.Head(url)
	if err != nil {
		return "", err
	}
	contentType := resp.Header.Get(http.CanonicalHeaderKey("Content-Type"))
	if extensionName, ok := GetExtensionNameByMimeType(contentType); ok {
		return extensionName, nil
	}
	return "", errors.New(fmt.Sprintf("unknown mimetype: %s", contentType))
}

func GetRemoteFileSize(httpClient *http.Client, url string) (int64, error) {
	resp, err := httpClient.Head(url)
	if err != nil {
		return 0, err
	}
	return resp.ContentLength, nil
}

func IsPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileSize
// Return file size in bytes
func GetFileSize(path string) (int64, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func EnsureDir(path string) error {
	if !IsPathExist(path) {
		if err := Mkdir(path); err != nil {
			return err
		}
	}
	return nil
}

func EnsureDirAll(path string) error {
	if !IsPathExist(path) {
		if err := MkdirAll(path); err != nil {
			return err
		}
	}
	return nil
}

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
		if line != "" && IsValidHttpLink(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func WriteFile(content string, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	_, err = out.WriteString(content)
	return err
}
