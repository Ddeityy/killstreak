package main

import (
	"path"
	"strings"
)

// Returns the demo name without path and extension
func TrimDemoName(demoPath string) string {
	_, demoName := path.Split(demoPath)
	return strings.TrimSuffix(demoName, ".dem")
}
