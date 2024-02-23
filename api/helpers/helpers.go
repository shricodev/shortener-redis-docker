package helpers

import (
	"net/url"
	"os"
	"strings"
)

// If the user provides the main server domain to abuse the system
func CheckForApplicationMatch(url_str string) bool {
	application_host := os.Getenv("APPLICATION_HOST")
	application_port := os.Getenv("APPLICATION_PORT")

	parsed_url, err := url.Parse(url_str)
	if err != nil {
		return false
	}

	parsed_host := parsed_url.Hostname()
	parsed_port := parsed_url.Port()

	if parsed_host == application_host && parsed_port == application_port {
		return false
	}

	return true
}

// Ensures the URL has a valid HTTP or HTTPS prefix
func EnforceHTTP(url_str string) string {
	if !strings.HasPrefix(url_str, "http://") && !strings.HasPrefix(url_str, "https://") {
		return "http://" + url_str
	}
	return url_str
}
