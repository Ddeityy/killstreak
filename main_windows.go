//go:build windows
// +build windows

package main

import "flag"

var cut bool

func main() {
	autoCut := flag.Bool("cut", true, "Automatically cut the demo")
	flag.Parse()
	cut = *autoCut
	WatchDemosDir()
}
