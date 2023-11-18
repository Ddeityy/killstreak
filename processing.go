package main

import (
	"encoding/json"
	"log"
)

// Process demo and write the result to _events.txt
func ProcessDemo(demoPath string) error {
	data := RustParseDemo(demoPath)
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
