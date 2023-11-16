//go:build linux
// +build linux

package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"k8s.io/utils/inotify"
)

// Watch for inotify events and process new demos
func WatchDemosDir() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	demosDir := GetDemosDir()

	err = watcher.Watch(demosDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching", demosDir)
	for {
		select {
		case event := <-watcher.Event:
			if event.Mask == inotify.InCloseWrite {
				if event.Name[len(event.Name)-4:] != ".dem" {
					break
				}
				log.Println("Finished writing:", event.Name)
				// Check if demo was auto-deleted by ds_stop
				time.Sleep(time.Millisecond * 100)
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					log.Println("Demo deleted:", err)
					break
				}
				log.Println("Processing demo:", trimDemoName(event.Name))
				ProcessDemo(event.Name)
			}
		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}
}

// TODO add "playdemo $demopath; demo_gototick $tick 0 (offset) 1 (pause)"
// Replaces default killstreak logs with custom ones in _event.txt
func (p *Player) WriteKillstreaksToEvents() {
	demosDir := GetDemosDir()
	eventsFile := path.Join(demosDir, "_events.txt")

	file, err := os.ReadFile(eventsFile)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Reading _events.txt")
	lines := strings.Split(string(file), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Killstreak") {
			if strings.Contains(line, p.DemoName) {
				prefix := line[:18]
				for _, k := range p.Killstreaks {
					lines[i-1] = fmt.Sprintf(">\n%v %v %v", prefix, p.MapName, p.MainClass)
					lines[i] = fmt.Sprintf(
						`%s Killstreak %v ("%v" %v-%v [%.2f seconds])`,
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
	lines = removeDuplicateLines(lines, ">")

	output := strings.Join(lines, "\n")

	err = os.WriteFile(eventsFile, []byte(output), 0644)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Finished: %+v", p.Killstreaks)
}
