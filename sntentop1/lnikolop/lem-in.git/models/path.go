package models

import (
	"strings"
)

type Path struct {
	Rooms   []*Room
	Tunnels []*Tunnel
}

func NewPath(rooms []*Room) *Path {
	path := &Path{
		Rooms:   rooms,
		Tunnels: []*Tunnel{},
	}
	path.assignTunnels()
	return path
}

func (p *Path) intermediateRooms() []*Room {
	if len(p.Rooms) <= 2 {
		return []*Room{}
	}
	return p.Rooms[1 : len(p.Rooms)-1]
}

func (p *Path) CheckCommonRooms(path *Path) bool {
	intermediate := p.intermediateRooms()
	for _, room := range intermediate {
		for _, room2 := range path.intermediateRooms() {
			if room.Name == room2.Name {
				return true
			}
		}
	}
	return false
}

func (p *Path) String() string {
	rooms := []string{}
	for _, room := range p.Rooms {
		rooms = append(rooms, room.Name)
	}
	return strings.Join(rooms, "-")
}

func (p *Path) assignTunnels() {
	for i, room := range p.Rooms {
		if i == len(p.Rooms)-1 {
			break
		}
		tunnel := NewTunnel(room, p.Rooms[i+1])
		p.Tunnels = append(p.Tunnels, tunnel)
	}
}

func (p *Path) GetTunnelIndex(tunnel *Tunnel) int {
	for i, t := range p.Tunnels {
		if t == tunnel {
			return i
		}
	}
	return -1
}
