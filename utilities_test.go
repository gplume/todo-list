package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBaseDir(t *testing.T) {
	defaultBaseDir, _ := os.Getwd()
	funcDefaultBaseDir := getDefaultBaseDir()
	assert.Equal(t, defaultBaseDir, funcDefaultBaseDir, "should be equal")
}

func TestGetScheme(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8000", nil)
	assert.Nil(t, err)
	schem := getScheme(req)
	assert.Equal(t, "http", schem, "should be equal")
}

func TestBaseUrl(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8000/test/base/url1", nil)
	assert.Nil(t, err)
	baseURL := getBaseURL(req)
	assert.Equal(t, "http://localhost:8000", baseURL, "should be equal")
}
