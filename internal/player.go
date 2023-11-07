package internal

import (
	"log"
)

type Kill struct {
	Tick float64
}

type Player struct {
	UserId      int
	Kills       []Kill
	Killstreaks []Killstreak
}

type Killstreak struct {
	Kills     []Kill
	StartTick float64
	EndTick   float64
	Length    float64
}

const killInterval = 15.0 // P-REC default = 15.0
const tick = 0.015

func (p *Player) GetPlayerKills(d Demo) {
	var userKills []Kill
	for _, v := range d.Deaths {
		if v.Killer != v.Victim {
			if v.Killer == p.UserId {
				userKills = append(userKills, Kill{Tick: v.Tick - d.StartTick})
			}
		}
	}
	if len(userKills) == 0 {
		log.Panicln("No kills found")
	}
	p.Kills = userKills
}

func NewPlayer(playerId int) Player {
	return Player{UserId: playerId}
}

func (p *Player) FindKillstreaks() {

	lastKill := p.Kills[0]

	killstreak := Killstreak{StartTick: lastKill.Tick}

	for _, currentKill := range p.Kills[1:] {

		timeBetweenKills := (currentKill.Tick - lastKill.Tick) * tick
		log.Println("Seconds since last kill: ", timeBetweenKills)

		killstreak.Kills = append(killstreak.Kills, lastKill)
		log.Println(len(killstreak.Kills), " kills")
		if timeBetweenKills <= killInterval {
			killstreak.EndTick = currentKill.Tick
		} else {
			if len(killstreak.Kills) >= 4 {
				killstreak.Length = (killstreak.EndTick - killstreak.StartTick) * tick
				p.Killstreaks = append(p.Killstreaks, killstreak)
				log.Printf("Added killstreak: %+v", killstreak)
			}
			killstreak = Killstreak{StartTick: currentKill.Tick}
			log.Println("Resetting killstreak")
		}
		lastKill = currentKill
	}
	p.printKillstreaks()
}

func (p *Player) printKillstreaks() {
	for i, v := range p.Killstreaks {
		log.Println("-------------")
		log.Printf("Killstreak %v \n", i+1)
		log.Printf("%v seconds long\n", v.Length)
		lastKill := 0.0
		for i, kill := range v.Kills {
			if lastKill == 0.0 {
				log.Printf("Kill %v - %v [-0 seconds]", i+1, kill.Tick)
				lastKill = kill.Tick
				continue
			}
			log.Printf("Kill %v - %v [-%.2f seconds]", i+1, kill.Tick, (kill.Tick-lastKill)*tick)
			lastKill = kill.Tick
		}
	}
}
