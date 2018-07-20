package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBaseDir(t *testing.T) {
	defaultBaseDir, _ := os.Getwd()
	funcDefaultBaseDir := getDefaultBaseDir()
	assert.Equal(t, defaultBaseDir, funcDefaultBaseDir, "should be equal")
}
