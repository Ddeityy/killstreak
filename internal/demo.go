package internal

import "errors"

type Chat struct {
	Kind string  `json:"kind"`
	From string  `json:"from"`
	Text string  `json:"text"`
	Tick float64 `json:"tick"`
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
	Weapon   string  `json:"weapon"`
	Victim   int     `json:"victim"`
	Assister int     `json:"assister"`
	Killer   int     `json:"killer"`
	Tick     float64 `json:"tick"`
}
type Rounds struct {
	Winner  string  `json:"winner"`
	Length  float64 `json:"length"`
	EndTick int     `json:"end_tick"`
}

type Demo struct {
	Chat      []Chat        `json:"chat"`
	Users     map[int]Users `json:"users"`
	Deaths    []Deaths      `json:"deaths"`
	Rounds    []Rounds      `json:"rounds"`
	StartTick float64       `json:"startTick"`
}

func (d *Demo) GetPlayerId(steamId string) (int, error) {
	for _, v := range d.Users {
		if v.SteamId == steamId {
			return v.UserId, nil
		}
	}
	return 0, errors.New("could match steamID to demo")
}
