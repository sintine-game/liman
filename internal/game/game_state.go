package game

type GameState struct {
	Rooms map[string]*Room
}

func NewGame() *GameState {
	return &GameState{
		Rooms: make(map[string]*Room),
	}
}

func (g *GameState) AddRoom(room *Room) {
	g.Rooms[room.ID] = room
}

func (g *GameState) RemoveRoom(roomID string) {
	delete(g.Rooms, roomID)
}

func (g *GameState) GetRoom(roomID string) *Room {
	return g.Rooms[roomID]
}
