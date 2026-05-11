package graph

import (
	"errors"
	"sort"

	"lem-in/structs"
)

// BuildGraph creates a graph (map) of the ant farm using the list of rooms and tunnels.
func BuildGraph(roomList []structs.Room, connections []structs.Tunnel) (*structs.Graph, error) {
	graphData := &structs.Graph{
		Rooms:     make(map[string]*structs.Room),
		Neighbors: make(map[string][]string),
	}

	// Add each room
	for i := range roomList {
		r := roomList[i]
		graphData.Rooms[r.Name] = &r
	}

	// Add each tunnel (undirected)
	for _, t := range connections {
		if _, ok := graphData.Rooms[t.RoomA]; !ok {
			return nil, errors.New("ERROR: tunnel refers to unknown room " + t.RoomA)
		}
		if _, ok := graphData.Rooms[t.RoomB]; !ok {
			return nil, errors.New("ERROR: tunnel refers to unknown room " + t.RoomB)
		}
		graphData.Neighbors[t.RoomA] = append(graphData.Neighbors[t.RoomA], t.RoomB)
		graphData.Neighbors[t.RoomB] = append(graphData.Neighbors[t.RoomB], t.RoomA)
	}

	return graphData, nil
}

// GetOptimalPaths returns the maximum set of simple paths from the start room
// to the end room, with no shared intermediate rooms.
func GetOptimalPaths(farmGraph *structs.Graph) ([][]string, error) {
	startRoom, endRoom := findEndpoints(farmGraph)
	if startRoom == "" || endRoom == "" {
		return nil, errors.New("missing start or end room")
	}

	neighborMap := farmGraph.Neighbors
	routeCandidates := enumerateRoutes(neighborMap, startRoom, endRoom)
	if len(routeCandidates) == 0 {
		return nil, errors.New("no paths found")
	}

	selectedRoutes := pickSeparateRoutes(routeCandidates)
	if len(selectedRoutes) == 0 {
		return nil, errors.New("no disjoint paths found")
	}
	return selectedRoutes, nil
}

// findEndpoints locates and returns the names of the start and end rooms.
func findEndpoints(farmGraph *structs.Graph) (string, string) {
	var startRoom, endRoom string
	for roomName, room := range farmGraph.Rooms {
		if room.IsStart {
			startRoom = roomName
		}
		if room.IsEnd {
			endRoom = roomName
		}
	}
	return startRoom, endRoom
}

// enumerateRoutes uses a stack-based search to find every simple path
// from startRoom to endRoom.
func enumerateRoutes(neighborMap map[string][]string, startRoom, endRoom string) [][]string {
	type stackFrame struct {
		currentRoom string
		nextIndex   int
	}

	var allRoutes [][]string
	visited := make(map[string]bool)
	visited[startRoom] = true

	currentPath := []string{startRoom}
	stack := []stackFrame{{currentRoom: startRoom, nextIndex: 0}}

	for len(stack) > 0 {
		frame := &stack[len(stack)-1]
		room := frame.currentRoom

		if room == endRoom {
			// record currentPath
			route := make([]string, len(currentPath))
			copy(route, currentPath)
			allRoutes = append(allRoutes, route)

			// backtrack
			visited[room] = false
			currentPath = currentPath[:len(currentPath)-1]
			stack = stack[:len(stack)-1]
			continue
		}

		if frame.nextIndex >= len(neighborMap[room]) {
			// no neighbors left, backtrack
			visited[room] = false
			currentPath = currentPath[:len(currentPath)-1]
			stack = stack[:len(stack)-1]
			continue
		}

		// explore next neighbor
		nextRoom := neighborMap[room][frame.nextIndex]
		frame.nextIndex++
		if visited[nextRoom] {
			continue
		}

		visited[nextRoom] = true
		currentPath = append(currentPath, nextRoom)
		stack = append(stack, stackFrame{currentRoom: nextRoom, nextIndex: 0})
	}

	return allRoutes
}

// pickSeparateRoutes scores each candidate path by how often its intermediate rooms
// appear, then picks routes in increasing order of that score (ties by shorter length),
// ensuring no room is used twice.
func pickSeparateRoutes(routes [][]string) [][]string {
	// count how often each room appears in the middle of routes
	roomCount := make(map[string]int)
	for _, route := range routes {
		for _, room := range route[1 : len(route)-1] {
			roomCount[room]++
		}
	}

	// build a list of scored routes
	type rankedRoute struct {
		rooms    []string
		crowding int
		length   int
	}
	ranked := make([]rankedRoute, 0, len(routes))
	for _, route := range routes {
		score := 0
		for _, room := range route[1 : len(route)-1] {
			score += roomCount[room]
		}
		ranked = append(ranked, rankedRoute{
			rooms:    route,
			crowding: score,
			length:   len(route),
		})
	}

	// sort by (lower crowding) then (shorter route)
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].crowding != ranked[j].crowding {
			return ranked[i].crowding < ranked[j].crowding
		}
		return ranked[i].length < ranked[j].length
	})

	// pick routes, avoiding reuse of intermediate rooms
	usedRooms := make(map[string]bool)
	var selected [][]string
	for _, rr := range ranked {
		ok := true
		for _, room := range rr.rooms[1 : len(rr.rooms)-1] {
			if usedRooms[room] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		for _, room := range rr.rooms[1 : len(rr.rooms)-1] {
			usedRooms[room] = true
		}
		selected = append(selected, rr.rooms)
	}

	return selected
}
