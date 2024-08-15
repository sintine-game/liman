package game

type Player struct {
	ID      string
	Fleet   []ShipPart
	LastHit Position
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:      id,
		Fleet:   make([]ShipPart, 0),
		LastHit: Position{X: -1, Y: -1},
	}
}

func (p *Player) AddShipPart(part ShipPart) {
	p.Fleet = append(p.Fleet, part)
}

func (p *Player) AddShipParts(parts []ShipPart) {
	p.Fleet = append(p.Fleet, parts...)
}

func (p *Player) TryHit(aim Position) bool {
	for i, ship := range p.Fleet {
		if ship.TryHit(aim) {
			p.Fleet[i] = ship
			return true
		}
	}
	return false
}

func (p *Player) GetLastHitAndReset() Position {
	lastHit := p.LastHit
	p.LastHit = Position{X: -1, Y: -1}
	return lastHit
}

func (p *Player) IsDefeated() bool {
	for _, parts := range p.Fleet {
		if !parts.IsHit {
			return false
		}
	}
	return true
}

func (p *Player) IsFleetDeployed() bool {
	return len(p.Fleet) == 8
}
