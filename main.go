package main

import (
	"killstreak/internal"
)

func main() {
	steamdId := internal.GetUserSteamId()
	demosDir := internal.GetDemosDir()
	internal.WatchDemosDir(demosDir, steamdId)
}
