package main

// Main struct of the entire demo file
type Demo struct {
	Name   string
	Path   string
	Header Header `json:"header"`
	State  State  `json:"state"`
	Player Player
}

type Header struct {
	Nick string `json:"nick"`
	Map  string `json:"map"`
}

type State struct {
	Users     map[int]Users `json:"users"`
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

// Returns player's userId in the demo
func (d *Demo) GetUserId(steamId string) int {
	for _, v := range d.State.Users {
		if v.SteamId == d.Player.Username {
			return v.UserId
		}
	}
	return 0
}

// Returns player's most used class
func (d *Demo) getPlayerClass() string {
	maxNum := 0
	var result int
	for _, user := range d.State.Users {
		if user.UserId == d.Player.UserId {
			for k, v := range user.Classes {
				if v > maxNum {
					maxNum = v
					result = k
				}
			}
		}
	}
	return classes[result]
}
