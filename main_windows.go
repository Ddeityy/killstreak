//go:build windows
// +build windows

package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

var cut bool

func main() {
	autoCut := flag.Bool("cut", true, "Automatically cut the demo")
	flag.Parse()
	cut = *autoCut
	WatchDemosDir()
}

func formatDemos() {
	demosDir := getDemosDir()
	demos, _ := os.ReadDir(demosDir)
	for _, demo := range demos {
		if strings.Contains(demo.Name(), ".dem") {
			log.Println("------------------------------------------------")
			log.Println("Processing:", demo.Name())
			demo := Demo{Path: path.Join(demosDir, demo.Name())}
			demo.ProcessDemo()
		}
	}
}
