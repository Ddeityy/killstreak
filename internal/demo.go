package internal

// Main struct of the entire demo file
type Demo struct {
	Chat      []Chat        `json:"chat"`
	Users     map[int]Users `json:"users"`
	Deaths    []Deaths      `json:"deaths"`
	Rounds    []Rounds      `json:"rounds"`
	StartTick float64       `json:"startTick"`
}

// Chat entries
type Chat struct {
	Kind string  `json:"kind"`
	From string  `json:"from"`
	Text string  `json:"text"`
	Tick float64 `json:"tick"`
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

// All rounds and who won
type Rounds struct {
	Winner  string  `json:"winner"`
	Length  float64 `json:"length"`
	EndTick int     `json:"end_tick"`
}

// Returns player's userId in the demo
func (d *Demo) GetPlayerId(steamId string) int {
	for _, v := range d.Users {
		if v.SteamId == steamId {
			return v.UserId
		}
	}
	return 0
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

// Returns player's most used class
func (d *Demo) getPlayerClass(userId int) string {
	maxNum := 0
	result := 0
	for _, user := range d.Users {
		if user.UserId == userId {
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
