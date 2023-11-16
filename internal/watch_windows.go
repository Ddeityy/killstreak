//go:build windows
// +build windows

package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ddeityy/steamlocate-go"
)

// On WRITE event wait a bit and:
//  1. Timer based WRITE check - bad
//  2. Try to read/write/copy a demo being written - lock?
func WatchDemosDir()

// TODO add "playdemo $demopath; demo_gototick $tick 0 (offset) 1 (pause)"
// Replaces default killstreak logs with custom ones in _event.txt
func (p *Player) WriteKillstreaksToEvents() {
	steamDir := steamlocate.SteamDir{}
	steamDir.Locate()
	demosDir := steamdir.SteamApps.Apps[440].Path
	demosDir = path.Join(demosDir, "tf")

	eventsFile := path.Join(demosDir, "KillStreaks.txt")

	file, err := os.ReadFile(eventsFile)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Reading _events.txt")
	lines := strings.Split(string(file), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Kill Streak") {
			if strings.Contains(line, p.DemoName) {
				prefix := line[:18]
				for _, k := range p.Killstreaks {
					lines[i-1] = fmt.Sprintf(">\n%v %v %v", prefix, p.MapName, p.MainClass)
					lines[i] = fmt.Sprintf(
						`%s Kill Streak: %v ("%v" %v-%v [%.2f seconds])`,
						prefix,
						len(k.Kills),
						p.DemoName,
						k.StartTick,
						k.EndTick,
						k.Length,
					)
				}
			}
		}
	}
	lines = removeDuplicateLines(lines, "\n")

	output := strings.Join(lines, "\n")

	err = os.WriteFile(eventsFile, []byte(output), 0644)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Finished: %+v", p.Killstreaks)
}
