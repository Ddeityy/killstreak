package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func ProcessDemo(demoPath string, steamId string) {
	data := ParseDemo(demoPath)
	demo := Demo{}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		log.Println("Parse error:", err)
	}

	playerId := demo.GetPlayerId(steamId)

	p := Player{UserId: playerId}
	p.GetPlayerKills(demo, demoPath)
	if len(p.Kills) == 0 {
		log.Println("No kills found - bookmark")
		return
	}
	p.FindKillstreaks()
	p.WriteKillstreaksToEvents()
}

// Watch for inotify events and process new demos
func WatchDemosDir(demosDir string, steamId string) {
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
				time.Sleep(time.Millisecond * 1)
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					log.Println("Demo deleted:", err)
					break
				}
				log.Println("Processing demo:", trimDemoName(event.Name))
				ProcessDemo(event.Name, steamId)
			}
		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}
}

// Returns the steamID(V3) from userdata directory
func GetUserSteamId() string {
	s := steamlocate.SteamDir{}
	s.Locate()
	file, err := os.ReadDir(path.Join(s.Path, "userdata"))
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("[U:1:%s]", file[0].Name())
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
