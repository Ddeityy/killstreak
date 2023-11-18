package main

import (
	"errors"
	"fmt"
	"log"
)

// Main player struct to retrieve killstreaks
type Player struct {
	Demo        *Demo
	DemoName    string
	Username    string
	UserId      int
	MapName     string
	Kills       []Kill
	Killstreaks []Killstreak
	MainClass   string // Most spawned as class
}

type Killstreak struct {
	Kills     []Kill
	StartTick float64
	EndTick   float64
	Length    float64 // Seconds
}

type Kill struct {
	Tick float64
}

const killInterval = 15.0 // P-REC default = 15.0
const tick = 0.015        // Amount of seconds per tick

// Populates the kills, mainclass and demoname fields
func (p *Player) GetPlayerKills() error {
	var userKills []Kill
	for _, v := range p.Demo.State.Deaths {
		if v.Killer != v.Victim {
			if v.Killer == p.UserId {
				userKills = append(userKills, Kill{Tick: v.Tick - p.Demo.State.StartTick})
			}
		}
	}
	if len(p.Kills) <= 3 {
		return fmt.Errorf("Less than 3 kills found, aborting")
	}
	p.MainClass = p.Demo.getPlayerClass()
	p.Kills = userKills
	p.DemoName = TrimDemoName(p.Demo.Path)
	return nil
}

func (p *Player) processKills() error {
	p.GetUserId()
	err := p.GetPlayerKills()
	if err != nil {
		log.Println(err)
		return err
	}
	err = p.GetUserKillstreaks()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Formatting and writing killstreaks")
	p.WriteKillstreaksToEvents()
	return nil
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

	lastKill := p.Kills[0]

	killstreak := Killstreak{StartTick: lastKill.Tick}

	for _, currentKill := range p.Kills[1:] {
		timeBetweenKills := (currentKill.Tick - lastKill.Tick) * tick
		killstreak.Kills = append(killstreak.Kills, lastKill)

		if timeBetweenKills <= killInterval {
			killstreak.EndTick = currentKill.Tick
		} else {
			if len(killstreak.Kills) >= 4 {
				killstreak.Length = (killstreak.EndTick - killstreak.StartTick) * tick
				p.Killstreaks = append(p.Killstreaks, killstreak)
			}
			killstreak = Killstreak{StartTick: currentKill.Tick}
		}
		lastKill = currentKill
	}
	if len(p.Killstreaks) == 0 {
		return errors.New("No killstreaks found")
	}
	return nil
}

// Removes duplicate killstreaks and keeps the separator
func removeDuplicateLines(s []string, separator string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
			delete(inResult, separator)
		}
	}
	return result
}
