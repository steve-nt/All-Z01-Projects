package models

import "fmt"

type Tunnel struct {
	Room1    *Room
	Room2    *Room
	Occupied bool
}

func NewTunnel(room1 *Room, room2 *Room) *Tunnel {
	return &Tunnel{
		Room1:    room1,
		Room2:    room2,
		Occupied: false,
	}
}

func (t *Tunnel) String() string {
	return fmt.Sprintf("%s-%s", t.Room1.Name, t.Room2.Name)
}
