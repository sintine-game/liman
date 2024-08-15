package main

type Client struct {
	Id    string
	Fleet []ShipPart
}

func NewClient(id string) Client {
	return Client{
		Id:    id,
		Fleet: make([]ShipPart, 0),
	}
}

func (c *Client) AddShip(ship ShipPart) {
	c.Fleet = append(c.Fleet, ship)
}

func (c *Client) AddShips(ships []ShipPart) {
	c.Fleet = append(c.Fleet, ships...)
}

func (c *Client) TryHit(aim Position) bool {
	for i, ship := range c.Fleet {
		if ship.Position.X == aim.X && ship.Position.Y == aim.Y {
			c.Fleet[i].IsHit = true
			return true
		}
	}
	return false
}

func (c *Client) IsDefeated() bool {
	for _, ship := range c.Fleet {
		if !ship.IsHit {
			return false
		}
	}
	return true
}

func (c *Client) IsDeployed() bool {
	return len(c.Fleet) == 8
}
