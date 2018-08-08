package utils

import (
	"fmt"
	"net/http"
	"os"
)

// GetDefaultBaseDir return PWD
func GetDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}

// GetScheme return http scheme
func GetScheme(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	}
	return "https"
}

// GetBaseURL return schyeme with host
func GetBaseURL(r *http.Request) string {
	return fmt.Sprintf("%s://%s", GetScheme(r), r.Host)
}
