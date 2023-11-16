//go:build linux
// +build linux

package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

// Parses demo and returns a JSON string to unmarshal
func ParseDemo(demoPath string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	parserPath, err := filepath.Abs(path.Join(homeDir, ".local", "share", "parse_demo"))
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
