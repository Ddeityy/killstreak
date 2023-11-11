package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"k8s.io/utils/inotify"
)

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string) {
	data := ParseDemo(demoPath)
	demo := Demo{Path: demoPath}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		log.Println("Parse error:", err)
	}

	p := Player{Username: demo.Header.Nick, MapName: demo.Header.Map, Demo: &demo}

	demo.Player = p

	log.Println("Processing kills")
	demo.Player.processKills()
}

func (p *Player) processKills() {
	p.GetUserId()
	p.GetPlayerKills()
	p.GetUserKillstreaks()
	if len(p.Killstreaks) == 0 {
		log.Println("No killstreaks found.")
		p.WriteKillstreaksToEvents()
		return
	}
	log.Println("Writing killstreaks")
	p.WriteKillstreaksToEvents()
}

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

// Parses demo and returns a JSON string to unmarshal
func ParseDemo(demoPath string) string {
	parserPath, err := filepath.Abs("parse_demo")
	if err != nil {
		log.Println(err)
	}
	command := exec.Command(parserPath, demoPath)
	var out bytes.Buffer

	command.Stdout = &out
	err = command.Run()
	if err != nil {
		log.Println(err)
	}
	return out.String()
}
