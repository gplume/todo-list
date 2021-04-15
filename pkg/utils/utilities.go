package utils

import (
	"fmt"
	"io/ioutil"
	"log"
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

// SearchDir dir is the parent directory you want to search
func SearchDir(dirPath, dir string) bool {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() && file.Name() == dir {
			return true
		}
	}
	return false
}
