package structs

// Room holds a room's data.
type Room struct {
	Name    string
	X       int
	Y       int
	IsStart bool
	IsEnd   bool
}

// Tunnel represents a connection between two rooms.
type Tunnel struct {
	RoomA string
	RoomB string
}

// Graph stores rooms and adjacency.
type Graph struct {
	Rooms     map[string]*Room
	Neighbors map[string][]string
}

// PathAssignment maps paths to ant counts.
type PathAssignment struct {
	Paths       [][]string
	AntsPerPath []int
}

// PathSim tracks ants on a path.
type PathSim struct {
	Path      []string
	Positions []int
	AntIDs    []int
}
