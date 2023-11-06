package main

import (
	"encoding/json"
	"killstreak/internal"
	"log"
)

// one - 16000-18500
// two - 74000-75000
// get user's steamID using tf2 appmanifest: LastUser

func main() {
	data := internal.ParseDemo("demos/one.dem")
	// _ = os.WriteFile("demo.json", []byte(data), 0644)
	demo := internal.Demo{}
	err := json.Unmarshal([]byte(data), &demo)
	if err != nil {
		log.Fatalf("Failed to unmarshal demo: %s", err)
	}
	steamId := internal.GetUserSteamId()
	playerId := demo.GetPlayerId(steamId)
	player := demo.GetPlayer(playerId)
	for _, kill := range player.Kills {
		log.Println(kill.Tick)
	}
}
