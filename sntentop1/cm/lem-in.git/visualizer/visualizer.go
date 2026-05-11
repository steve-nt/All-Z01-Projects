package visualizer

import (
	"fmt"
	"os"
	"strings"

	"lem-in/structs"
)

func buildRawInput(antTotal int, roomList []structs.Room, tunnelList []structs.Tunnel) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d\n", antTotal))
	for _, room := range roomList {
		if room.IsStart {
			builder.WriteString("##start\n")
		}
		if room.IsEnd {
			builder.WriteString("##end\n")
		}
		builder.WriteString(fmt.Sprintf("%s %d %d\n", room.Name, room.X, room.Y))
	}
	for _, tunnel := range tunnelList {
		builder.WriteString(fmt.Sprintf("%s-%s\n", tunnel.RoomA, tunnel.RoomB))
	}
	return builder.String()
}

func buildSummary(antTotal int, roomList []structs.Room, tunnelList []structs.Tunnel) string {
	var builder strings.Builder
	builder.WriteString("----------- Summary -----------\n")
	builder.WriteString(fmt.Sprintf("Number of ants: %d\n", antTotal))
	builder.WriteString(fmt.Sprintf("Number of rooms: %d\n", len(roomList)))
	builder.WriteString(fmt.Sprintf("Number of tunnels: %d\n", len(tunnelList)))
	var startRoom, endRoom string
	for _, room := range roomList {
		if room.IsStart {
			startRoom = room.Name
		}
		if room.IsEnd {
			endRoom = room.Name
		}
	}
	builder.WriteString(fmt.Sprintf("Start room: %s\n", startRoom))
	builder.WriteString(fmt.Sprintf("End room: %s\n", endRoom))
	return builder.String()
}

func buildAllPaths(pathList [][]string) string {
	var builder strings.Builder
	builder.WriteString("---------- All Found Paths ----------\n")
	builder.WriteString(fmt.Sprintf("Number of possible paths: %d\n", len(pathList)))
	for i, path := range pathList {
		builder.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(path, " -> ")))
	}
	return builder.String()
}

func buildSelectedPaths(assignment structs.PathAssignment) string {
	var builder strings.Builder
	builder.WriteString("---------- Selected Paths ----------\n")
	for i, path := range assignment.Paths {
		builder.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(path, " -> ")))
	}
	return builder.String()
}

// PrintExtraInfo builds the header with input, summary, and path info.
func PrintExtraInfo(antTotal int, roomList []structs.Room, tunnelList []structs.Tunnel,
	pathList [][]string, assignment structs.PathAssignment) string {
	var builder strings.Builder
	builder.WriteString(buildRawInput(antTotal, roomList, tunnelList))
	builder.WriteString("\n")
	builder.WriteString(buildSummary(antTotal, roomList, tunnelList))
	builder.WriteString("\n")
	builder.WriteString(buildAllPaths(pathList))
	builder.WriteString("\n")
	builder.WriteString(buildSelectedPaths(assignment))
	builder.WriteString("\n")
	return builder.String()
}

// GeneratePathGrid renders one path, marking any ants present.
func GeneratePathGrid(sim structs.PathSim) string {
	var builder strings.Builder
	for i, room := range sim.Path {
		var antLabels []string
		for j, pos := range sim.Positions {
			if pos == i {
				antLabels = append(antLabels, fmt.Sprintf("L%d", sim.AntIDs[j]))
			}
		}
		if len(antLabels) > 0 {
			builder.WriteString(fmt.Sprintf("[ %s (%s) ]", room, strings.Join(antLabels, ", ")))
		} else {
			builder.WriteString(fmt.Sprintf("[ %s ]", room))
		}
		if i < len(sim.Path)-1 {
			builder.WriteString(" ---> ")
		}
	}
	return builder.String()
}

// WriteSimulationOutput writes the full simulation grid to a file.
func WriteSimulationOutput(filename string, headerInfo string,
	turnGrids []string, totalTurns int) error {
	var builder strings.Builder
	builder.WriteString(headerInfo)
	builder.WriteString("\n\n")
	for i, grid := range turnGrids {
		builder.WriteString(fmt.Sprintf("TURN %d\n", i+1))
		builder.WriteString(grid)
		builder.WriteString("\n")
	}
	builder.WriteString(fmt.Sprintf("Total turns: %d\n", totalTurns))
	return os.WriteFile(filename, []byte(builder.String()), 0644)
}

// PrintTerminalOutput prints concise move info per turn.
func PrintTerminalOutput(moveList []string) {
	for i, moves := range moveList {
		fmt.Printf("Turn %d: %s\n", i+1, moves)
	}
}
