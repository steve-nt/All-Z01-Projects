package main

import (
	"fmt"
	"strings"
)

// SimulateAnts moves ants along the given path according to lem-in rules.
// Only ants that move are shown. One move per ant per turn. Rooms can hold only one ant,
// except for start (unlimited) and end (unlimited).
func SimulateAnts(path []string, antCount int) {
	type Ant struct {
		ID    int
		Index int // index in the path (starts from 1)
	}

	var active []Ant // ants currently moving on the path
	antID := 1       // ID to assign to next ant
	pathLen := len(path)
	occupancy := make([]int, pathLen) // keeps track of which room is occupied (by ant ID), excluding start/end

	for {
		var moves []string
		var newActive []Ant

		// Move ants that are already on the path
		for _, ant := range active {
			nextIndex := ant.Index + 1
			if nextIndex == pathLen || occupancy[nextIndex] == 0 {
				// free to move
				if ant.Index < pathLen-1 {
					occupancy[ant.Index] = 0 // leave current room
				}
				if nextIndex < pathLen-1 {
					occupancy[nextIndex] = ant.ID // occupy next room
				}
				ant.Index = nextIndex
				if ant.Index < pathLen {
					newActive = append(newActive, ant)
				}
				// only print if still in valid range
				if ant.Index < len(path) {
					moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, path[ant.Index]))
				}
			} else {
				// Room occupied, wait
				newActive = append(newActive, ant)
			}
		}

		// Start a new ant if available
		if antID <= antCount && len(path) > 1 && occupancy[1] == 0 {
			occupancy[1] = antID
			newActive = append(newActive, Ant{ID: antID, Index: 1})
			moves = append(moves, fmt.Sprintf("L%d-%s", antID, path[1]))
			antID++
		}

		// If no moves were made, we're done
		if len(moves) == 0 {
			break
		}

		// Print the current move step
		fmt.Println(strings.Join(moves, " "))
		active = newActive
	}
}
