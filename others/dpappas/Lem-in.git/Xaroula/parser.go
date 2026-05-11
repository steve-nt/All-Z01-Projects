package main

import (
	"strconv"
	"strings"
)

// Room represents a room in the ant colony.
type Room struct {
	Name    string
	X, Y    int
	IsStart bool
	IsEnd   bool
}

// ParseRooms receives a slice of lines and extracts all valid rooms.
// It detects ##start and ##end flags and marks the next room accordingly.
func ParseRooms(lines []string) ([]Room, error) {
	var rooms []Room
	expectStart := false
	expectEnd := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and handle special commands
		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				expectStart = true
			} else if line == "##end" {
				expectEnd = true
			}
			continue
		}

		// A valid room line must have exactly 3 parts: name, x, y
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])

			if err1 == nil && err2 == nil {
				room := Room{
					Name:    parts[0],
					X:       x,
					Y:       y,
					IsStart: expectStart,
					IsEnd:   expectEnd,
				}
				rooms = append(rooms, room)

				// Reset special flags so only the next room is marked
				expectStart = false
				expectEnd = false
			}
		}
	}

	return rooms, nil
}

// Graph represents the connectivity between rooms
type Graph map[string][]string

// ParseConnections takes all lines and builds a graph of room links (edges)
// It assumes that room names are valid and known.
func ParseConnections(lines []string, rooms []Room) Graph {
	graph := make(Graph)

	// Create a map for fast lookup to verify room existence
	roomMap := make(map[string]bool)
	for _, room := range rooms {
		roomMap[room.Name] = true
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Connection lines have the form: name1-name2
		if strings.Count(line, "-") == 1 && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "-")
			a, b := parts[0], parts[1]

			// Skip if one of the rooms doesn't exist
			if !roomMap[a] || !roomMap[b] {
				continue
			}

			// Undirected graph: add both directions
			graph[a] = append(graph[a], b)
			graph[b] = append(graph[b], a)
		}
	}

	return graph
}
