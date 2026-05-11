package main

// FindPath uses BFS to find the shortest path between start and end in the graph.
// Returns the path as a slice of room names, or nil if no path is found.
func FindPath(graph Graph, start, end string) []string {
	if start == end {
		return []string{start}
	}

	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := []string{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				prev[neighbor] = current
				queue = append(queue, neighbor)

				if neighbor == end {
					// Reconstruct path
					var path []string
					for at := end; at != ""; at = prev[at] {
						path = append([]string{at}, path...)
					}
					return path
				}
			}
		}
	}

	return nil // No path found
}
