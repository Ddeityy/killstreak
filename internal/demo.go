package internal

// Main struct of the entire demo file
type Demo struct {
	Header Header `json:"header"`
	State  State  `json:"state"`
}

type Header struct {
	DemoType string  `json:"demo_type"`
	Version  int     `json:"version"`
	Protocol int     `json:"protocol"`
	Server   string  `json:"server"`
	Nick     string  `json:"nick"`
	Map      string  `json:"map"`
	Game     string  `json:"game"`
	Duration float64 `json:"duration"`
	Ticks    int     `json:"ticks"`
	Frames   int     `json:"frames"`
	Signon   int     `json:"signon"`
}

type State struct {
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
func (d *Demo) GetUserId() int {
	for _, v := range d.State.Users {
		if v.Name == d.Header.Nick {
			return v.UserId
		}
	}
	return 0
}

// Class enum given by demo parser
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
	var result int
	for _, user := range d.State.Users {
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
