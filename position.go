package main

import (
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

func ExtractPositions(s string) []Position {
	positions := make([]Position, 0)
	for _, p := range strings.Split(s, ";") {
		coords := strings.Split(p, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		positions = append(positions, NewPosition(x, y))
	}
	return positions
}
