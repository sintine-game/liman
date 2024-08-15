package main

import (
	"log"

	"github.com/sintine-game/liman/internal/game"
	"github.com/sintine-game/liman/internal/handler"
	"github.com/sintine-game/liman/pkg/udpserver"
)

func main() {
	gameState := game.NewGame()

	s, err := udpserver.NewUDPServer(":9177")
	if err != nil {
		log.Fatal(err)
	}

	s.RegisterHandler("ping", handler.PingHandler)
	s.RegisterHandler("cr", handler.CreateRoomHandler(gameState))
	s.RegisterHandler("jr", handler.JoinRoomHandler(gameState))
	s.RegisterHandler("gp", handler.GetPlayersInRoomHandler(gameState))
	s.RegisterHandler("df", handler.DeployFleetHandler(gameState))
	s.RegisterHandler("pf", handler.GetPlayerFleetHandler(gameState))
	s.RegisterHandler("if", handler.IsRoomFullHandler(gameState))
	s.RegisterHandler("id", handler.IsOpponentDeployed(gameState))
	s.RegisterHandler("th", handler.TryHitHandler(gameState))
	s.RegisterHandler("gh", handler.GetHitHandler(gameState))

	log.Println("Server is running on :9177")

	s.Listen()
}
