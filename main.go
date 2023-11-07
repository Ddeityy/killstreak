package main

import (
	"encoding/json"
	"killstreak/internal"
	"log"
)

func main() {
	data := internal.ParseDemo("demos/two.dem")
	demo := internal.Demo{}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		log.Fatalf("Failed to unmarshal demo: %s", err)
	}
	steamId := internal.GetUserSteamId()
	playerId := demo.GetPlayerId(steamId)
	player := internal.NewPlayer(playerId)
	player.GetPlayerKills(demo)
	player.FindKillstreaks()
}
