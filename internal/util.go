package internal

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/ddeityy/steamlocate-go"
)

// Returns the absolute path of /demos
func GetDemosDir() string {
	steamdir := steamlocate.SteamDir{}
	steamdir.Locate()
	demosDir := steamdir.SteamApps.Apps[440].Path
	demosDir = path.Join(demosDir, "tf", "demos")
	if _, err := os.Stat(demosDir); os.IsNotExist(err) {
		log.Fatalf("Demos folder doesn't exist: %v", err)
	}
	return demosDir
}

// Returns the demo name without path and extension
func trimDemoName(demoPath string) string {
	demoName := strings.Split(demoPath, "/")
	demoName = demoName[len(demoName)-1:]
	demoNameStrip := demoName[0]
	return strings.TrimSuffix(demoNameStrip, ".dem")
}
