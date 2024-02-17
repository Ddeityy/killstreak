package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ddeityy/steamlocate-go"
	"k8s.io/utils/inotify"
)

// Returns the demo name without path and extension
func TrimDemoName(demoPath string) string {
	_, demoName := path.Split(demoPath)
	return strings.TrimSuffix(demoName, ".dem")
}

// Returns the absolute path of the directory with most .dem files
func GetDemosDir() (string, error) {
	steamDir := steamlocate.SteamDir{}
	steamDir.Locate()

	tfDir := steamDir.SteamApps.Apps[440].Path
	tfDir = path.Join(tfDir, "tf")
	demosDir := path.Join(tfDir, "demos")

	tf, err := countDemos(tfDir)
	if err != nil {
		return "", err
	}

	demos, err := countDemos(demosDir)
	if err != nil {
		return "", err
	}

	if tf == 0 && demos == 0 {
		return "", errors.New("no demos found in either folder")
	}

	if tf > demos {
		return tfDir, nil
	} else {
		return demosDir, nil
	}
}

func countDemos(dir string) (int, error) {
	demos := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".dem") {
			demos++
		}
	}

	return demos, nil
}

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string, demosDir string) error {
	log.Println("Processing demo.")
	data := RustParseDemo(demoPath)
	demo, err := NewDemo(data, demoPath, demosDir)
	if err != nil {
		return err
	}

	log.Println("Processing kills.")
	if err = demo.Player.ProcessEvents(); err != nil {
		return fmt.Errorf("demo: %w", err)
	}

	return nil
}

// Watch for inotify events and process new demos
func WatchDemosDir() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	demosDir, err := GetDemosDir()
	if err != nil {
		log.Fatal(err)
	}

	if err = watcher.Watch(demosDir); err != nil {
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
				log.Println("Finished writing", event.Name)

				// Check if demo was auto-deleted by ds_stop
				time.Sleep(time.Millisecond * 100)
				if _, err := os.Stat(event.Name); os.IsNotExist(err) {
					log.Println("demo deleted:", err)
					break
				}

				log.Println("Processing demo", TrimDemoName(event.Name))
				if err := ProcessDemo(event.Name, demosDir); err != nil {
					log.Println("demo processing:", err)
				}

			}
		case err := <-watcher.Error:
			log.Println("watcher:", err)
		}
	}
}

func FormatDemos() {
	demosDir, err := GetDemosDir()
	if err != nil {
		panic(err)
	}

	demos, _ := os.ReadDir(demosDir)
	for _, demo := range demos {
		if strings.Contains(demo.Name(), ".dem") {
			log.Println("------------------------------------------------")
			log.Println("Processing", demo.Name())
			ProcessDemo(path.Join(demosDir, demo.Name()), demosDir)
		}
	}
}
