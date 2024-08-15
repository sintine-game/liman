package game

type Room struct {
	ID      string
	Players map[string]*Player
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
	}
}

func (r *Room) AddPlayer(player *Player) {
	r.Players[player.ID] = player
}

func (r *Room) RemovePlayer(playerID string) {
	delete(r.Players, playerID)
}

func (r *Room) GetPlayer(playerID string) *Player {
	return r.Players[playerID]
}

func (r *Room) GetPlayers() map[string]*Player {
	return r.Players
}

func (r *Room) IsFull() bool {
	return len(r.Players) >= 2
}

func (r *Room) IsReadyToStart() bool {
	if len(r.Players) <= 2 {
		return false
	}

	for _, player := range r.Players {
		if len(player.Fleet) <= 4 {
			return false
		}
	}

	return true
}

func (r *Room) IsDefeated(playerID string) bool {
	return r.Players[playerID].IsDefeated()
}

func (r *Room) TryHit(playerID string, aim Position) bool {
	return r.Players[playerID].TryHit(aim)
}

func (r *Room) GetOpponent(playerID string) *Player {
	for id, player := range r.Players {
		if id != playerID {
			return player
		}
	}
	return nil
}

func (r *Room) GetOpponentID(playerID string) string {
	for id := range r.Players {
		if id != playerID {
			return id
		}
	}
	return ""
}
