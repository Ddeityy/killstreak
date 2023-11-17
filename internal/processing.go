package internal

import (
	"encoding/json"
	"log"
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
