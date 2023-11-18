//go:build windows
// +build windows

package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/ddeityy/steamlocate-go"
	"github.com/fsnotify/fsnotify"
)

// On WRITE event wait a bit and:
//  1. Timer based WRITE check - bad
//  2. Try to read/write/copy a demo being written - lock?
func WatchDemosDir() {
	//demosDir := getDemosDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					if event.Name[len(event.Name)-4:] != ".dem" {
						break
					}
					// Check if demo was auto-deleted by prec_delete_useless_demo
					time.Sleep(time.Millisecond * 100)
					if _, err := os.Stat(event.Name); os.IsNotExist(err) {
						log.Println("Demo deleted:", err)
						break
					}
					log.Println("Processing demo:", TrimDemoName(event.Name))
					err := ProcessDemo(event.Name)
					if err != nil {
						log.Println(err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(`test_data\windows`)
	//	err = watcher.Add(demosDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Watching test")

	<-make(chan struct{})
}

func getDemosDir() string {
	steamDir := steamlocate.SteamDir{}
	steamDir.Locate()
	demosDir := steamDir.SteamApps.Apps[440].Path
	demosDir = path.Join(demosDir, "tf")
	return demosDir
}

// TODO add "playdemo $demopath; demo_gototick $tick 0 (offset) 1 (pause)"
// Replaces default killstreak logs with custom ones in _event.txt
func (p *Player) WriteKillstreaksToEvents() {
	//demosDir := getDemosDir()
	//eventsFile := path.Join(demosDir, "KillStreaks.txt")

	eventsFile := path.Join("test_data", "windows", "KillStreaks.txt")

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

// Parses demo and returns a JSON string to unmarshal
func ParseDemo(demoPath string) string {

	command := exec.Command(`.\bin\parse_demo.exe`, demoPath)

	var out bytes.Buffer

	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
	return out.String()
}
