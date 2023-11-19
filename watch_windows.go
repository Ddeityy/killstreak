//go:build windows
// +build windows

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ddeityy/steamlocate-go"
	"github.com/fsnotify/fsnotify"
)

// On WRITE event wait a bit and:
//  1. Timer based WRITE check - bad
//  2. Try to read/write/copy a demo being written - lock?
func WatchDemosDir() {
	demosDir := getDemosDir()
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

					// Check if demo was auto-deleted
					demo := ParseDemo(event.Name)
					if demo == "File not found" {
						log.Println(demo)
						break
					}

					if strings.Contains(demo, `"duration":0.0`) {
						break
					}

					log.Println("Processing demo:", TrimDemoName(event.Name))
					err := ProcessDemo(event.Name)
					if err != nil {
						log.Println("Error:", err)
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
	err = watcher.Add(demosDir)
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

// Replaces default killstreak logs with custom ones in killstreaks.txt
func (p *Player) WriteKillstreaksToEvents() {
	demosDir := getDemosDir()
	eventsFile := path.Join(demosDir, "KillStreaks.txt")

	file, err := os.ReadFile(eventsFile)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Reading _events.txt")
	lines := strings.Split(string(file), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Kill Streak") {
			if strings.Contains(line, p.Demo.Name) {
				prefix := line[:18]
				for _, k := range p.Killstreaks {
					ticks := fmt.Sprintf("playdemo demos/%v; demo_gototick %v 0 1", p.Demo.Name, k.StartTick)
					header := fmt.Sprintf("%v %v %v", prefix, p.Demo.Header.Map, p.MainClass)
					streak := fmt.Sprintf(
						`%s Killstreak %v ("%v" %v-%v [%.2f seconds])`,
						prefix,
						len(k.Kills),
						p.Demo.Name,
						k.StartTick,
						k.EndTick,
						k.Length,
					)
					var l []string
					l = append(l, ticks, header, streak)
					lines[i] = strings.Join(l, "\n")
				}
			}
		}
	}
	lines = removeDuplicateLines(lines, "\n")

	output := strings.Join(lines, "\n")

	err = os.WriteFile(eventsFile, []byte(output), 0644)
	if err != nil {
		log.Println("Error:", err)
	}
	log.Printf("Finished: %+v", p.Killstreaks)
}

func ParseDemo(demoPath string) string {
	command := exec.Command(".\\parse_demo.exe", demoPath)

	var out bytes.Buffer

	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
	return out.String()
}

func CutDemo(demoPath string, startTick int32) {
	command := exec.Command(".\\cut_demo.exe", demoPath, string(startTick))

	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string) error {
	data := ParseDemo(demoPath)
	demo := Demo{Path: demoPath}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		return err
	}

	p := Player{Username: demo.Header.Nick, Demo: &demo}

	demo.Player = p

	log.Println("Processing kills.")
	err = demo.Player.processKills()
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	return nil
}
