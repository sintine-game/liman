package main

import "fmt"

type Room struct {
	Id           string
	Clients      []Client
	ReadyCh      chan bool
	FleetReadyCh chan bool
}

func NewRoom(id string) *Room {
	return &Room{
		Id:           id,
		Clients:      make([]Client, 0),
		ReadyCh:      make(chan bool),
		FleetReadyCh: make(chan bool),
	}
}

func (r *Room) AddClient(client Client) error {
	if r.Clients == nil {
		r.Clients = make([]Client, 0)
	}

	if len(r.Clients) < 2 {
		r.Clients = append(r.Clients, client)
		if r.IsReady() {
			r.ReadyCh <- true
		}
		return nil
	}

	return fmt.Errorf("Room is full")
}

func (r *Room) IsReady() bool {
	return len(r.Clients) == 2
}

func (r *Room) IsFleetReady() bool {
	for _, client := range r.Clients {
		if !client.IsDeployed() {
			return false
		}
	}

	return true
}
