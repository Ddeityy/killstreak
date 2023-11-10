package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/ddeityy/steamlocate-go"
	"k8s.io/utils/inotify"
)

// get user's steamID using tf2 appmanifest: LastUser

// Returns the absolute path of /demos
func GetDemosDir() string {
	steamdir := steamlocate.SteamDir{}
	steamdir.Locate()
	demosDir := steamdir.SteamApps.Apps[440].Path
	demosDir = path.Join(demosDir, "tf", "demos")
	if _, err := os.Stat(demosDir); os.IsNotExist(err) {
		log.Fatalf("Demos folder doesn't exist: %v", err)
	}
	return demosDir
}

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string) {
	data := ParseDemo(demoPath)
	demo := Demo{}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		log.Println("Parse error:", err)
	}

	p := Player{Username: demo.Header.Nick, MapName: demo.Header.Map, UserId: demo.GetUserId()}
	p.GetPlayerKills(demo, demoPath)
	if len(p.Kills) == 0 {
		log.Println("Only bookmards found.")
		return
	}
	log.Println("Gettinng killstreaks")
	p.FindKillstreaks()
	log.Println("Writing killstreaks")
	p.WriteKillstreaksToEvents()
}

// Watch for inotify events and process new demos
func WatchDemosDir(demosDir string) {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

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
