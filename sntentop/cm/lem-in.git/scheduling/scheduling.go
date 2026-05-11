package scheduling

import (
	"sort"

	"lem-in/structs"
)

// AssignAnts distributes ants among paths by minimizing cost = len+assigned-1.
func AssignAnts(antCount int, paths [][]string) structs.PathAssignment {
	numPaths := len(paths)
	antsPerPath := make([]int, numPaths)

	for i := 0; i < antCount; i++ {
		type pathCost struct {
			index int
			cost  int
		}
		pathCosts := make([]pathCost, numPaths)
		for j := 0; j < numPaths; j++ {
			pathCosts[j] = pathCost{index: j, cost: len(paths[j]) + antsPerPath[j] - 1}
		}
		sort.Slice(pathCosts, func(a, b int) bool {
			return pathCosts[a].cost < pathCosts[b].cost
		})
		antsPerPath[pathCosts[0].index]++
	}

	return structs.PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
