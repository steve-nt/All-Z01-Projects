package models

import (
	"fmt"
)

type Ant struct {
	Id          int
	Path        *Path
	Room        *Room
	TunnelIndex int
	Finished    bool
}

func NewAnt(id int, room *Room) *Ant {
	ant := &Ant{
		Id:   id,
		Path: NewPath([]*Room{room}),
		Room: room,
	}
	return ant
}

func (a *Ant) TakeTunnel(tunnel *Tunnel) string {
	// (1) if the tunnel is nil or the ant is finished, return
	if tunnel == nil || a.Finished {
		return ""
	}

	// (1) free the tunnel
	tunnel.Occupied = false

	// (3) move the ant to the next room
	a.Room = tunnel.Room2
	token := fmt.Sprintf("L%d-%s ", a.Id, a.Room.Name)

	// (4) if the ant has reached the end of the path, mark it as finished
	if a.Room == a.Path.Rooms[len(a.Path.Rooms)-1] {
		a.Finished = true
	}

	return token
}

func (a *Ant) PickTunnel() *Tunnel {
	// (1) check if the ant is finished
	if a.Id == -1 || a.Finished {
		return nil
	}

	// (2) get the current tunnel position
	curPtr := a.Path.Tunnels[a.TunnelIndex]
	index := a.Path.GetTunnelIndex(curPtr)

	// (3) collision check
	if a.Path.Tunnels[index].Occupied {
		return nil
	}

	// (4) safety bounds check
	if index < 0 || index >= len(a.Path.Tunnels) {
		return nil
	}

	// (5) advance the ant's pointer to the next time
	if a.TunnelIndex != len(a.Path.Tunnels)-1 {
		a.TunnelIndex = index + 1
	}

	// (6) mark the tunnel as occupied
	a.Path.Tunnels[index].Occupied = true

	return a.Path.Tunnels[index]
}
