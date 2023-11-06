package internal

import "log"

type Kill struct {
	Tick int
}

type Player struct {
	UserId int
	Kills  []Kill
}

type Killstreak struct {
	StartTick int
	EndTick   int
}

func (p *Player) FindKillstreaks() string {
	timeBetweenKills := 15.0
	tick := 0.015

	firstKill := p.Kills[0]

	for _, kill := range p.Kills[1:] {
		if ((float64(kill.Tick) - float64(firstKill.Tick)) * tick) <= timeBetweenKills {
			log.Println(kill.Tick)
		}
	}

	return ""
}
