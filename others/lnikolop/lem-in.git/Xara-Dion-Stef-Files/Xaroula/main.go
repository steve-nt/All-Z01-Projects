package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if a filename was provided as an argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	filename := os.Args[1]

	// Read all lines from the file
	lines, err := ReadLines(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Get the number of ants from the first line
	var antCount int
	if len(lines) > 0 {
		fmt.Sscanf(lines[0], "%d", &antCount)
	}

	// Parse all room lines from the file
	rooms, err := ParseRooms(lines)
	if err != nil {
		fmt.Println("Error parsing rooms:", err)
		return
	}

	// Print all parsed rooms
	fmt.Println("\n🔍 Rooms found:")
	for _, r := range rooms {
		fmt.Printf("- %s (%d,%d)", r.Name, r.X, r.Y)
		if r.IsStart {
			fmt.Print(" [START]")
		}
		if r.IsEnd {
			fmt.Print(" [END]")
		}
		fmt.Println()
	}

	// Optionally: print the full raw file content
	fmt.Println("\n📄 Raw input file content:")
	for _, line := range lines {
		fmt.Println(line)
	}

	graph := ParseConnections(lines, rooms)

	fmt.Println("\n🔗 Room Connections:")
	for room, connections := range graph {
		fmt.Printf("- %s is connected to: %v\n", room, connections)
	}

	// Find start and end room names
	var startName, endName string
	for _, room := range rooms {
		if room.IsStart {
			startName = room.Name
		}
		if room.IsEnd {
			endName = room.Name
		}
	}

	path := FindPath(graph, startName, endName)
	if path == nil {
		fmt.Println("\n🚫 No path found from start to end.")
	} else {
		fmt.Println("\n✅ Shortest path from start to end:")
		for i, p := range path {
			if i > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(p)
		}
		fmt.Println()
	}

	if path != nil {
		fmt.Println("\n🐜 Ant movements:")
		SimulateAnts(path, antCount)
	}

}
