package helpers

import (
	"net/url"
	"os"
	"strings"
)

// If the user provides the main server domain to abuse the system
func CheckForApplicationMatch(url_str string) bool {
	applicationHost := os.Getenv("APPLICATION_HOST")
	applicationPort := os.Getenv("APPLICATION_PORT")

	parsedUrl, err := url.Parse(url_str)
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

// Check if a string is a valid URL
func IsURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
