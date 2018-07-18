package main

import (
	"fmt"
	"net/http"
	"os"
)

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}

func getScheme(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	}
	return "https"
}

func getBaseURL(r *http.Request) string {
	return fmt.Sprintf("%s://%s", getScheme(r), r.Host)
}
