package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Kill struct {
	Tick float64
}

type Player struct {
	DemoName    string
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

func (p *Player) GetPlayerKills(d Demo, demoPath string) {
	var userKills []Kill
	for _, v := range d.Deaths {
		if v.Killer != v.Victim {
			if v.Killer == p.UserId {
				userKills = append(userKills, Kill{Tick: v.Tick - d.StartTick})
			}
		}
	}
	p.Kills = userKills
	p.DemoName = trimDemoName(demoPath)
}

func trimDemoName(demoPath string) string {
	demoName := strings.Split(demoPath, "/")
	demoName = demoName[len(demoName)-1:]
	demoNameStrip := demoName[0]
	return strings.TrimSuffix(demoNameStrip, ".dem")
}

func NewPlayer(playerId int) Player {
	return Player{UserId: playerId}
}

func (p *Player) FindKillstreaks() {

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
	p.printKillstreaks()
}

func (p *Player) WriteKillstreaksToEvents() {
	demosDir := GetDemosDir()
	eventsFile := path.Join(demosDir, "_events.txt")
	file, err := os.ReadFile(eventsFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	lines := strings.Split(string(file), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Killstreak") {
			if strings.Contains(line, p.DemoName) {
				prefix := line[:18]
				for _, k := range p.Killstreaks {
					lines[i] = fmt.Sprintf(`%s Killstreak %v ("%v" %v - %v [%.2f seconds])`, prefix, len(k.Kills), p.DemoName, k.StartTick, k.EndTick, k.Length)
				}
			}
		}
	}
	lines = removeDuplicateLines(lines)

	for i, line := range lines {
		fmt.Println(i, line)
	}

	output := strings.Join(lines, "\n")

	err = os.WriteFile(eventsFile, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeDuplicateLines(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
			delete(inResult, ">")
		}
	}
	return result
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
