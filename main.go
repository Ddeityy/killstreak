package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// one - 16000-18500
// two - 74000-75000

type Chat struct {
	Kind string `json:"kind"`
	From string `json:"from"`
	Text string `json:"text"`
	Tick int    `json:"tick"`
}

type Classes map[int]int

type Users struct {
	Classes Classes `json:"classes"`
	Name    string  `json:"name"`
	UserId  int     `json:"userId"`
	SteamId string  `json:"steamId"`
	Team    string  `json:"team"`
}

type Deaths struct {
	Weapon   string `json:"weapon"`
	Victim   int    `json:"victim"`
	Assister int    `json:"assister"`
	Killer   int    `json:"killer"`
	Tick     int    `json:"tick"`
}
type Rounds struct {
	Winner  string  `json:"winner"`
	Length  float64 `json:"length"`
	EndTick int     `json:"end_tick"`
}

type Demo struct {
	Chat            []Chat        `json:"chat"`
	Users           map[int]Users `json:"users"`
	Deaths          []Deaths      `json:"deaths"`
	Rounds          []Rounds      `json:"rounds"`
	StartTick       int           `json:"startTick"`
	IntervalPerTick float32       `json:"intervalPerTick"`
}

func parseDemo(demoPath string) string {
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

func findSteamId() string {
	homeDir := os.Getenv("HOME")
	steamPath := path.Join(homeDir, ".steam", "steam")
	file, err := os.ReadDir(path.Join(steamPath, "userdata"))
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("[U:1:%s]", file[0].Name())
}

func findPlayerId(demo Demo, steamId string) int {
	for _, v := range demo.Users {
		if v.SteamId == steamId {
			return v.UserId
		}
	}
	return 0
}

// get user's steamID using tf2 appmanifest: LastUser

func main() {
	data := parseDemo("demos/one.dem")
	// _ = os.WriteFile("demo.json", []byte(data), 0644)
	demo := Demo{}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		return
	}
	log.Println(demo.Users[3].Name)
	log.Println(demo.Users[3].SteamId)
	log.Println(demo.Deaths[0].Killer, demo.Deaths[0].Weapon, demo.Deaths[0].Victim, demo.Deaths[0].Tick)
	steamid := findSteamId()
	userid := findPlayerId(demo, steamid)
	fmt.Println(userid)
}
