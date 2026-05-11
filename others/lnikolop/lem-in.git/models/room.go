package models

type Room struct {
	Name     string
	X        float32
	Y        float32
	PosX     float32
	PosY     float32
	Links    []*Room
	Occupied bool
}

func NewRoom(name string, x float32, y float32) *Room {
	return &Room{
		Name:     name,
		X:        x,
		Y:        y,
		PosX:     x * 2,
		PosY:     y * 2,
		Links:    []*Room{},
		Occupied: false,
	}
}
