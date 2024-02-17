package main

import (
	"encoding/json"
	"path"
	"strings"
)

// Main struct of the entire demo file
type Demo struct {
	Name             string
	Path             string
	Header           Header `json:"header"`
	State            State  `json:"state"`
	Player           Player
	Date             string
	EventsFile       string
	LegacyEventsFile string
	DemoDir          string
}

type Header struct {
	Nick     string  `json:"nick"`
	Map      string  `json:"map"`
	Duration float32 `json:"duration"`
}

type Message struct {
	Kind string `json:"kind"`
	From string `json:"from"`
	Text string `json:"text"`
	Tick int    `json:"tick"`
}

type State struct {
	Users     map[int]Users `json:"users"`
	Chat      []Message     `json:"chat"`
	Deaths    []Deaths      `json:"deaths"`
	StartTick float64       `json:"startTick"`
}

// All players in the demo
type Users struct {
	Classes map[int]int `json:"classes"`
	Name    string      `json:"name"`
	UserId  int         `json:"userId"`
	SteamId string      `json:"steamId"`
	Team    string      `json:"team"`
}

// All deaths in the demo
type Deaths struct {
	Weapon   string  `json:"weapon"`
	Victim   int     `json:"victim"`
	Assister int     `json:"assister"`
	Killer   int     `json:"killer"`
	Tick     float64 `json:"tick"`
}

// Class enums given by demo parser
var classes = map[int]string{
	0: "other",
	1: "scout",
	2: "sniper",
	3: "soldier",
	4: "demoman",
	5: "medic",
	6: "heavy",
	7: "pyro",
	8: "spy",
	9: "engineer",
}

func NewDemo(demoData string, demoPath string, demosDir string) (*Demo, error) {
	demo := Demo{}

	if err := json.Unmarshal([]byte(demoData), &demo); err != nil {
		return nil, err
	}

	p := Player{Username: demo.Header.Nick, Demo: &demo}
	p.GetUserId()
	p.GetClass()

	demo.Player = p
	demo.Path = demoPath
	demo.Name = TrimDemoName(p.Demo.Path)
	demo.Date = strings.Split(demo.Name, "_")[0]
	demo.DemoDir = demosDir
	demo.EventsFile = path.Join(demosDir, "events.txt")
	demo.LegacyEventsFile = path.Join(demosDir, "_events.txt")

	return &demo, nil
}
