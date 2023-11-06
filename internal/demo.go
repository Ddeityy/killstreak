package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/ddeityy/steamlocate-go"
)

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
	IntervalPerTick float32       `json:"intervalPerTick"` // Seconds
}

func (d *Demo) GetPlayerId(steamId string) int {
	for _, v := range d.Users {
		if v.SteamId == steamId {
			return v.UserId
		}
	}
	return 0
}

func (d *Demo) GetPlayer(playerId int) Player {
	var userKills []Kill
	for _, v := range d.Deaths {
		if v.Killer == playerId {
			userKills = append(userKills, Kill{Tick: v.Tick})
		}
	}
	if len(userKills) == 0 {
		log.Panicln("No kills found")
	}
	return Player{Kills: userKills, UserId: playerId}
}

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

func GetUserSteamId() string {
	s := steamlocate.SteamDir{}
	s.Locate()
	file, err := os.ReadDir(path.Join(s.Path, "userdata"))
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("[U:1:%s]", file[0].Name())
}
