package funcs

// buildConnections creates an adjacency list from tunnels
func BuildConnections() {

	connections = make(map[string][]string)
	for _, side := range tunnels {
		a, b := side[0], side[1]
		connections[a] = append(connections[a], b)
		connections[b] = append(connections[b], a)
	}
}
