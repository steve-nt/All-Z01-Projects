package funcs

import (
	"fmt"
	"os"
)

// We use BFS for quick lookup to find at least one valid connection start -> end
func StartEndConnection() {
	visited := make(map[string]bool) // keep track of rooms visited
	queue := []string{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			return // success: path exists
		}

		visited[current] = true
		for _, neighbor := range connections[current] { // gives you all rooms connected to the current room
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}

	fmt.Println("[ERROR]: valid paths not found")
	os.Exit(0)
}
