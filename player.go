package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

// Main player struct to retrieve killstreaks
type Player struct {
	Demo        *Demo
	Username    string
	UserId      int
	Kills       []Kill
	Killstreaks []Killstreak
	Class       string // Most spawned as class
	Bookmarks   []Bookmark
}

type Bookmark struct {
	Tick int
}

type Killstreak struct {
	Kills     []Kill
	StartTick int
	EndTick   int
	Length    float64 // Seconds
}

type Kill struct {
	Tick int
}

const killInterval = 15.0 // P-REC default = 15.0
const tick = 0.015        // Amount of seconds per tick

// Populates the kills, mainclass and demoname fields
func (p *Player) GetPlayerKills() error {
	var userKills []Kill
	for _, v := range p.Demo.State.Deaths {
		if v.Killer != v.Victim {
			if v.Killer == p.UserId {
				userKills = append(userKills, Kill{Tick: int(v.Tick)})
			}
		}
	}

	p.Kills = userKills
	if len(p.Kills) <= 3 {
		return fmt.Errorf("less than 3 kills found, aborting")
	}
	return nil
}

func (p *Player) GetUserBookmarks() error {
	file, err := os.ReadFile(p.Demo.LegacyEventsFile)
	if err != nil {
		log.Printf("%v", err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.Contains(line, p.Demo.Name) {
			if strings.Contains(strings.ToLower(line), "bookmark") {
				ticks := strings.Split(line, " ")
				tick := ticks[len(ticks)-1]
				tick = strings.TrimSuffix(tick, ")")
				log.Println(tick)
				intTick, _ := strconv.Atoi(tick)
				p.Bookmarks = append(p.Bookmarks, Bookmark{Tick: intTick})
			}
		}
	}
	if len(p.Bookmarks) == 0 {
		return fmt.Errorf("no bookmarks found")
	}
	return nil
}

func (p *Player) ProcessEvents() error {
	kErr := p.GetUserKillstreaks()
	bErr := p.GetUserBookmarks()

	if kErr != nil && bErr != nil {
		return fmt.Errorf("no killstreaks or bookmarks found")
	}

	log.Println("Formatting and writing killstreaks.")
	p.WriteEvents()
	return nil
}

func NewPlayer(d *Demo) Player {
	p := Player{Username: d.Header.Nick, Demo: d}
	p.GetUserId()
	p.GetClass()
	return p
}

// Returns player's most used class
func (p *Player) GetClass() {
	maxNum := 0
	var result int
	for _, user := range p.Demo.State.Users {
		if user.UserId == p.UserId {
			for k, v := range user.Classes {
				if v > maxNum {
					maxNum = v
					result = k
				}
			}
		}
	}
	p.Class = classes[result]
}

func (p *Player) GetUserId() {
	for _, v := range p.Demo.State.Users {
		if v.Name == p.Username {
			p.UserId = v.UserId
		}
	}
}

// Finds all killstreaks
func (p *Player) GetUserKillstreaks() error {
	err := p.GetPlayerKills()
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	lastKill := p.Kills[0]

	killstreak := Killstreak{StartTick: lastKill.Tick}

	for _, currentKill := range p.Kills[1:] {
		timeBetweenKills := (float64(currentKill.Tick) - float64(lastKill.Tick)) * tick
		killstreak.Kills = append(killstreak.Kills, lastKill)

		if timeBetweenKills <= killInterval {
			killstreak.EndTick = currentKill.Tick
		} else {
			if len(killstreak.Kills) >= 4 {
				killstreak.Length = (float64(killstreak.EndTick) - float64(killstreak.StartTick)) * tick
				p.Killstreaks = append(p.Killstreaks, killstreak)
			}
			killstreak = Killstreak{StartTick: currentKill.Tick}
		}
		lastKill = currentKill
	}
	if len(p.Killstreaks) == 0 {
		return errors.New("no killstreaks found")
	}
	return nil
}

func (p *Player) WriteEvents() {
	file, err := os.OpenFile(p.Demo.EventsFile,
		os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	prefix := fmt.Sprintf("[%v]", p.Demo.Date)
	header := fmt.Sprintf("%v %v %v\n", prefix, p.Demo.Header.Map, p.Demo.Player.Class)
	_, playdemoPrefix := path.Split(p.Demo.DemoDir)
	if playdemoPrefix != "tf" {
		playdemoPrefix += string(os.PathSeparator)
	}

	var bookmarks []string
	for _, b := range p.Bookmarks {
		playdemo := fmt.Sprintf("playdemo %v%v; demo_gototick %v 0 1", playdemoPrefix, p.Demo.Name, b.Tick-500)
		bookmark := fmt.Sprintf(
			"%s Bookmark at %v",
			prefix,
			b.Tick,
		)
		bookmark = fmt.Sprintf("%v%v%v\n", bookmark, strings.Repeat(" ", 65-len(bookmark)), playdemo)
		bookmarks = append(bookmarks, bookmark)
	}

	var streaks []string
	for _, k := range p.Killstreaks {
		playdemo := fmt.Sprintf("playdemo %v%v; demo_gototick %v 0 1", playdemoPrefix, p.Demo.Name, k.StartTick-500)
		streak := fmt.Sprintf(
			"%s Killstreak %v %v-%v [%.2f seconds]",
			prefix,
			len(k.Kills),
			k.StartTick,
			k.EndTick,
			k.Length,
		)
		streak = fmt.Sprintf("%v%v%v\n", streak, strings.Repeat(" ", 65-len(streak)), playdemo)
		streaks = append(streaks, streak)
	}

	var lines []string
	lines = append(lines, header)
	lines = append(lines, streaks...)
	lines = append(lines, bookmarks...)

	log.Println("Writing to events")
	for _, line := range lines {
		file.WriteString(line)
	}
	file.WriteString(">\n")
	log.Println("Done")
}
