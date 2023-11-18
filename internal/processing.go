package internal

import (
	"encoding/json"
	"log"
	"os/exec"
)

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string) error {
	data := ParseDemo(demoPath)
	demo := Demo{Path: demoPath}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		return err
	}

	p := Player{Username: demo.Header.Nick, MapName: demo.Header.Map, Demo: &demo}

	demo.Player = p

	log.Println("Processing kills")
	err = demo.Player.processKills()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CutDemo(demoPath string, startTick int32) error {
	command := exec.Command(`bin\cut_demo.exe`, demoPath, string(startTick))
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}
