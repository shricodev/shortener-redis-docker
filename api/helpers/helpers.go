package helpers

import (
	"net/url"
	"os"
	"strings"
)

// If the user provides the main server domain to abuse the system
func CheckForApplicationMatch(urlStr string) bool {
	applicationHost := os.Getenv("APPLICATION_HOST")
	applicationPort := os.Getenv("APPLICATION_PORT")

	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	parsedHost := parsedUrl.Hostname()
	parsedPort := parsedUrl.Port()

	if parsedHost == applicationHost && parsedPort == applicationPort {
		return false
	}

	return true
}

// Ensures the URL has a valid HTTP or HTTPS prefix
func EnforceHTTP(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return "http://" + urlStr
	}
	return urlStr
}
