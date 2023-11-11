package main

import "killstreak/internal"

func main() {
	demosDir := internal.GetDemosDir()
	internal.WatchDemosDir(demosDir)
}
