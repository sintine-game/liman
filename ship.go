package main

type ShipPart struct {
	Position Position
	IsHit    bool
}

func NewShip(x, y int) ShipPart {
	return ShipPart{
		Position: Position{
			X: x,
			Y: y,
		},
		IsHit: false,
	}
}
