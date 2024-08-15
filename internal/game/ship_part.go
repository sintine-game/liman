package game

type ShipPart struct {
	Position Position
	IsHit    bool
}

func NewShipPart(position Position) ShipPart {
	return ShipPart{
		Position: position,
	}
}

func (sp *ShipPart) TryHit(aim Position) bool {
	if sp.Position.X == aim.X && sp.Position.Y == aim.Y {
		sp.IsHit = true
		return true
	}
	return false
}
