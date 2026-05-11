package funcs

import (
	"fmt"
	"sort"
	"strings"
)

type AntAssignment struct {
	antID     int // The global ant number (starting at 1)
	pathIndex int // Index into the paths slice (after sorting)
	order     int // The order of this ant on the path (1 means first ant on that path)
}

// sortPaths sorts the available paths in ascending order by their length.
func sortPaths() {

	sort.Slice(maxFlowPaths, func(i, j int) bool {
		return len(maxFlowPaths[i]) < len(maxFlowPaths[j])
	})
}

// assignAnts distributes the ants among the available paths using a greedy strategy.
// It returns a slice of AntAssignment.
func assignAnts() []AntAssignment {

	assignments := make([]AntAssignment, 0, totalAnts)
	// assignedCounts[i] holds the number of ants already assigned to paths[i]
	assignedCounts := make([]int, len(maxFlowPaths))

	for ant := 1; ant <= totalAnts; ant++ {
		// Find the path where the candidate finish time is minimal.
		bestIdx := 0
		// Candidate finish time = (assignedCounts + 1) + (len(path) - 1) = assignedCounts + len(path)
		bestFinish := assignedCounts[0] + len(maxFlowPaths[0])
		for i := 1; i < len(maxFlowPaths); i++ {
			candidate := assignedCounts[i] + len(maxFlowPaths[i])
			if candidate < bestFinish {
				bestFinish = candidate
				bestIdx = i
			}
		}
		// Assign this ant to the chosen path.
		assignedCounts[bestIdx]++
		assignments = append(assignments, AntAssignment{
			antID:     ant,
			pathIndex: bestIdx,
			order:     assignedCounts[bestIdx], // The order number on the chosen path.
		})
	}

	return assignments
}

// SimulateAnts runs a turn-by-turn simulation of ant movement along their assigned paths.
// It returns a slice of strings, where each string is one turn's moves in the format "L<antID>-<room>".
func SimulateAnts(assignments []AntAssignment) []string {
	// Determine the maximum number of turns required.
	maxTurn := 0
	for _, a := range assignments {
		finish := a.order + len(maxFlowPaths[a.pathIndex]) - 1
		if finish > maxTurn {
			maxTurn = finish
		}
	}
	// Prepare a map for quick lookup of an ant's assignment by its antID.
	antAssignMap := make(map[int]AntAssignment)
	for _, a := range assignments {
		antAssignMap[a.antID] = a
	}
	// We'll simulate from turn 1 to maxTurn.
	resultLines := make([]string, 0)
	// To keep moves in a consistent order, we'll iterate over ants in order of their ID.
	antIDs := make([]int, totalAnts)
	for i := 1; i <= totalAnts; i++ {
		antIDs[i-1] = i
	}
	// Simulate each turn.
	for t := 1; t <= maxTurn; t++ {
		moves := make([]string, 0)
		for _, antID := range antIDs {
			a := antAssignMap[antID]
			// An ant can only start moving from its launch turn (which is its order).
			if t >= a.order {
				// Its position index on the path (starting at 0 for the start room) is:
				// pos = t - a.order + 1, because on its launch turn it moves from start to the first room.
				pos := t - a.order + 1
				// Only print a move if the ant is still en route.
				// When pos equals len(path), the ant has reached the end.
				if pos < len(maxFlowPaths[a.pathIndex]) {
					room := maxFlowPaths[a.pathIndex][pos]
					moves = append(moves, fmt.Sprintf("L%d-%s", antID, room))
				}
			}
		}
		// If any moves occurred during this turn, record them.
		if len(moves) > 0 {
			resultLines = append(resultLines, strings.Join(moves, " "))
		}
	}
	return resultLines
}

// OptimalAntDistribution combines the above steps:
// It sorts the given paths, assigns ants to them using a greedy strategy,
// simulates their movements turn by turn, and returns the list of turn strings.
func OptimalAntDistribution() {

	sortPaths()

	assignments := assignAnts()

	result := SimulateAnts(assignments)
	for _, line := range result {
		fmt.Println(line)
	}
}
