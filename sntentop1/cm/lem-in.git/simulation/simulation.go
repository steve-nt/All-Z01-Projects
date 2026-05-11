package simulation

import (
	"fmt"
	"strings"

	"lem-in/structs"
	"lem-in/visualizer"
)

// initSimulation prepares simulation state for each path.
func initSimulation(pathList [][]string, assignment structs.PathAssignment) []structs.PathSim {
	simStates := make([]structs.PathSim, len(pathList))
	antIDCounter := 1

	for i, path := range pathList {
		antCountForPath := assignment.AntsPerPath[i]
		positions := make([]int, antCountForPath)
		for j := range positions {
			positions[j] = -1
		}
		antIDs := make([]int, antCountForPath)
		for j := 0; j < antCountForPath; j++ {
			antIDs[j] = antIDCounter
			antIDCounter++
		}
		simStates[i] = structs.PathSim{
			Path:      path,
			Positions: positions,
			AntIDs:    antIDs,
		}
	}
	return simStates
}

// isRoomOccupied checks if any ant is at the given index.
func isRoomOccupied(positions []int, roomIndex int) bool {
	for _, pos := range positions {
		if pos == roomIndex {
			return true
		}
	}
	return false
}

// processTurn moves ants one step along each path.
func processTurn(simStates []structs.PathSim) ([]string, string) {
	var moveDescriptions []string
	var gridBuilder strings.Builder

	for idx := range simStates {
		simState := &simStates[idx]
		pathLength := len(simState.Path)
		newPositions := make([]int, len(simState.Positions))
		copy(newPositions, simState.Positions)

		if pathLength == 2 {
			// direct path: inject one ant
			for j := range simState.Positions {
				if simState.Positions[j] == -1 {
					newPositions[j] = 1
					moveDescriptions = append(moveDescriptions,
						fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1]))
					break
				}
			}
		} else {
			// longer path: move existing ants first
			for j := len(simState.Positions) - 1; j >= 0; j-- {
				if simState.Positions[j] == -1 {
					if !isRoomOccupied(newPositions, 1) {
						newPositions[j] = 1
						moveDescriptions = append(moveDescriptions,
							fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1]))
					}
				} else if simState.Positions[j] < pathLength-1 {
					nextIndex := simState.Positions[j] + 1
					if nextIndex == pathLength-1 || !isRoomOccupied(newPositions, nextIndex) {
						newPositions[j] = nextIndex
						moveDescriptions = append(moveDescriptions,
							fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[nextIndex]))
					}
				}
			}
		}

		copy(simState.Positions, newPositions)
		gridBuilder.WriteString(visualizer.GeneratePathGrid(*simState) + "\n")
	}

	return moveDescriptions, gridBuilder.String()
}

// SimulateMultiPath runs the simulation until completion.
func SimulateMultiPath(antTotal int, pathList [][]string, assignment structs.PathAssignment, headerInfo string) {
	simStates := initSimulation(pathList, assignment)
	var moveOutputs, gridOutputs []string
	turnCount := 0

	for {
		moves, grid := processTurn(simStates)
		if len(moves) == 0 {
			break
		}
		moveOutputs = append(moveOutputs, strings.Join(moves, " "))
		gridOutputs = append(gridOutputs, grid)
		turnCount++
	}

	err := visualizer.WriteSimulationOutput("simulation_output.txt", headerInfo, gridOutputs, turnCount)
	if err != nil {
		fmt.Println("Error writing simulation output:", err)
	}
	visualizer.PrintTerminalOutput(moveOutputs)
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
