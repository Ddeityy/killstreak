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
	"time"

	"github.com/bep/debounce"
	"github.com/ddeityy/steamlocate-go"
	"github.com/fsnotify/fsnotify"
)

func WatchDemosDir() {
	demosDir := getDemosDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// refactor to use something like demo.Parse(path) to use the debouncer?
	dbounce := debounce.New(1000 * time.Millisecond)

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
					demo := Demo{Path: event.Name}
					dbounce(demo.ParseDemo)

					if demo.Header.Duration == 0.0 {
						break
					}

					log.Println("Processing demo:", TrimDemoName(event.Name))
					err := demo.ProcessDemo()
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
					var ticks string
					if cut {
						CutDemo(p.Demo.Path, k.StartTick)
						ticks = fmt.Sprintf("playdemo cut_%v;", p.Demo.Name)
					} else {
						ticks = fmt.Sprintf("playdemo %v; demo_gototick %v 0 1", p.Demo.Name, k.StartTick)
					}
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

func (d *Demo) ParseDemo() {
	command := exec.Command(".\\parse_demo.exe", d.Path)

	var out bytes.Buffer

	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(out.String()), &d)
	if err != nil {
		log.Println(err)
	}
}

func CutDemo(demoPath string, startTick int) {
	command := exec.Command(".\\cut_demo.exe", demoPath, string(startTick))

	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}

// Process demo and write the result to _events.txt
func (d *Demo) ProcessDemo() error {
	p := Player{Username: d.Header.Nick, Demo: d}

	d.Player = p

	log.Println("Processing kills.")
	err := d.Player.processKills()
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	return nil
}
