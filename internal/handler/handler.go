package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/sintine-game/liman/internal/game"
	"github.com/sintine-game/liman/pkg/udpserver"
)

func PingHandler(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	log.Println("ping from", addr.String())
	conn.WriteToUDP([]byte("pong"), addr)
}

// cr:room_id:player_id
func CreateRoomHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 3 {
			log.Println("invalid message format for create room: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]
		playerId := messageParts[2]

		room := game.NewRoom(roomId)
		room.AddPlayer(game.NewPlayer(playerId))
		gameState.AddRoom(room)

		conn.WriteToUDP([]byte("ok"), addr)
	}
}

// jr:room_id:player_id
func JoinRoomHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 3 {
			log.Println("invalid message format for join room: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]
		playerId := messageParts[2]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		room.AddPlayer(game.NewPlayer(playerId))
		conn.WriteToUDP([]byte("ok"), addr)
	}
}

// gp:room_id
func GetPlayersInRoomHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 2 {
			log.Println("invalid message format for get players in room: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		players := room.GetPlayers()
		playerList := make([]string, 0, len(players))
		for _, player := range players {
			playerList = append(playerList, player.ID)
		}
		resp := strings.Join(playerList, ",")
		conn.WriteToUDP([]byte(resp), addr)
	}
}

// df:room_id:player_id:ship_parts
func DeployFleetHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 4 {
			log.Println("invalid message format for deploy fleet: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		playerId := messageParts[2]
		player := room.GetPlayer(playerId)
		if player == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		shipParts := strings.Split(messageParts[3], ";")
		if len(shipParts) != 8 {
			log.Println("invalid message format for deploy fleet: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}

		ships := make([]game.ShipPart, 0, 8)
		for _, shipPart := range shipParts {
			positions := strings.Split(shipPart, ",")
			if len(positions) != 2 {
				log.Println("invalid message format for deploy fleet: ", message)
				conn.WriteToUDP([]byte("invalid message format"), addr)
				return
			}

			x, err := strconv.Atoi(positions[0])
			if err != nil {
				log.Println("invalid message format for deploy fleet: ", message)
				conn.WriteToUDP([]byte("invalid message format"), addr)
				return
			}

			y, err := strconv.Atoi(positions[1])
			if err != nil {
				log.Println("invalid message format for deploy fleet: ", message)
				conn.WriteToUDP([]byte("invalid message format"), addr)
				return
			}

			ships = append(ships, game.NewShipPart(game.Position{X: x, Y: y}))
		}

		player.AddShipParts(ships)
		conn.WriteToUDP([]byte("ok"), addr)
	}
}

// pf:room_id:player_id
func GetPlayerFleetHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 3 {
			log.Println("invalid message format for get player fleet: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		playerId := messageParts[2]
		player := room.GetPlayer(playerId)
		if player == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		shipParts := player.Fleet
		shipPartList := make([]string, 0, len(shipParts))
		for _, shipPart := range shipParts {
			shipPartList = append(shipPartList, strconv.Itoa(shipPart.Position.X)+","+strconv.Itoa(shipPart.Position.Y))
		}
		resp := strings.Join(shipPartList, ";")
		conn.WriteToUDP([]byte(resp), addr)
	}
}

// if:room_id
func IsRoomFullHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 2 {
			log.Println("invalid message format for is room ready: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		if room.IsFull() {
			conn.WriteToUDP([]byte("yes"), addr)
		} else {
			conn.WriteToUDP([]byte("no"), addr)
		}
	}
}

// id:room_id:player_id
func IsOpponentDeployed(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 3 {
			log.Println("invalid message format for is opponent deployed: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		playerId := messageParts[2]
		player := room.GetOpponent(playerId)
		if player == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		if player.IsFleetDeployed() {
			conn.WriteToUDP([]byte("yes"), addr)
		} else {
			conn.WriteToUDP([]byte("no"), addr)
		}
	}
}

// th:room_id:player_id:x,y
func TryHitHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 4 {
			log.Println("invalid message format for try hit: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		playerId := messageParts[2]
		opponent := room.GetOpponent(playerId)
		if opponent == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		player := room.GetPlayer(playerId)
		if player == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		positions := strings.Split(messageParts[3], ",")
		if len(positions) != 2 {
			log.Println("invalid message format for try hit: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}

		x, err := strconv.Atoi(positions[0])
		if err != nil {
			log.Println("invalid message format for try hit: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}

		y, err := strconv.Atoi(positions[1])
		if err != nil {
			log.Println("invalid message format for try hit: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}

		hit := room.TryHit(opponent.ID, game.Position{X: x, Y: y})
		if hit {
			conn.WriteToUDP([]byte("hit"), addr)
		} else {
			conn.WriteToUDP([]byte("miss"), addr)
		}

		player.LastHit = game.Position{X: x, Y: y}
	}
}

// gh:room_id:player_id
func GetHitHandler(gameState *game.GameState) func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
	return func(conn *udpserver.Conn, addr *udpserver.Addr, message string) {
		messageParts := strings.Split(message, ":")
		if len(messageParts) != 3 {
			log.Println("invalid message format for get hit: ", message)
			conn.WriteToUDP([]byte("invalid message format"), addr)
			return
		}
		roomId := messageParts[1]

		room := gameState.GetRoom(roomId)
		if room == nil {
			conn.WriteToUDP([]byte("room not found"), addr)
			return
		}

		playerId := messageParts[2]
		opponent := room.GetOpponent(playerId)
		if opponent == nil {
			conn.WriteToUDP([]byte("player not found"), addr)
			return
		}

		lastHit := opponent.GetLastHitAndReset()
		conn.WriteToUDP([]byte(strconv.Itoa(lastHit.X)+","+strconv.Itoa(lastHit.Y)), addr)
	}
}
