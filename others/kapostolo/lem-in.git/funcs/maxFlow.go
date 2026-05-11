package funcs

// This file implements a max flow solution to find vertex-disjoint paths
// using the standard technique of vertex splitting and the Edmonds–Karp algorithm.

type Edge struct {
	from     string
	to       string
	capacity int
	flow     int
	rev      *Edge
}

// addEdge adds a forward edge (with capacity) and a corresponding reverse edge (with 0 capacity)
// to the flow network.
// Adds a forward and reverse edge between two nodes
func addEdge(graph map[string][]*Edge, from, to string, capacity int) {
	forward := &Edge{from: from, to: to, capacity: capacity, flow: 0}
	reverse := &Edge{from: to, to: from, capacity: 0, flow: 0}
	forward.rev = reverse
	reverse.rev = forward
	graph[from] = append(graph[from], forward)
	graph[to] = append(graph[to], reverse)
}

// Returns "_in" or "_out" versions unless node is start/end
func transform(node, start, end, suffix string) string {
	if node == start || node == end {
		return node
	}
	return node + "_" + suffix
}

// BuildFlowNetwork constructs the flow graph with vertex splitting
func BuildFlowNetwork(start, end string) (map[string][]*Edge, map[string]string) {
	network := make(map[string][]*Edge)
	nodeMap := make(map[string]string)
	vertexSet := make(map[string]bool)

	// Collect all unique nodes
	for u, neighbors := range connections {
		vertexSet[u] = true
		for _, v := range neighbors {
			vertexSet[v] = true
		}
	}

	// Split non-start/end nodes and add internal edge
	for node := range vertexSet {
		if node == start || node == end {
			network[node] = []*Edge{}
			nodeMap[node] = node
		} else {
			in, out := transform(node, start, end, "in"), transform(node, start, end, "out")
			network[in] = []*Edge{}
			network[out] = []*Edge{}
			addEdge(network, in, out, 1)
			nodeMap[in], nodeMap[out] = node, node
		}
	}

	// Add undirected connections with proper node transforms
	seen := make(map[string]bool)
	for u, neighbors := range connections {
		for _, v := range neighbors {
			key := u + "_" + v
			if u > v {
				key = v + "_" + u
			}
			if seen[key] {
				continue
			}
			seen[key] = true

			uOut := transform(u, start, end, "out")
			vIn := transform(v, start, end, "in")
			addEdge(network, uOut, vIn, 1)

			vOut := transform(v, start, end, "out")
			uIn := transform(u, start, end, "in")
			addEdge(network, vOut, uIn, 1)
		}
	}

	return network, nodeMap
}

// bfs finds an augmenting path in the flow network and fills in the parent mapping.
// Returns true if a path from source to sink is found.
func bfs(network map[string][]*Edge, source, sink string, parent map[string]*Edge) bool {
	// Clear the parent map.
	for key := range parent {
		delete(parent, key)
	}
	queue := []string{source}
	visited := make(map[string]bool)
	visited[source] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, edge := range network[u] {
			// Compute the residual capacity.
			if edge.capacity-edge.flow > 0 && !visited[edge.to] {
				visited[edge.to] = true
				parent[edge.to] = edge
				if edge.to == sink {
					return true
				}
				queue = append(queue, edge.to)
			}
		}
	}
	return false
}

// MaxFlow computes the maximum flow from source to sink using the Edmonds–Karp algorithm.
func MaxFlow(network map[string][]*Edge, source, sink string) int {
	parent := make(map[string]*Edge)
	maxFlow := 0

	// While there is an augmenting path, send flow (which is always 1 in this construction).
	for bfs(network, source, sink, parent) {
		flow := 1
		// Traverse the path from sink back to source and update flows.
		for v := sink; v != source; {
			edge := parent[v]
			edge.flow += flow
			edge.rev.flow -= flow
			v = edge.from
		}
		maxFlow += flow
	}
	return maxFlow
}

// dfsExtract uses depth-first search to extract one path from source to sink along edges with flow > 0.
// As a path is found, the flow on the used edges is decremented to avoid using them again.
func dfsExtract(network map[string][]*Edge, u, sink string, path *[]string) bool {
	// Append current node to path.
	*path = append(*path, u)
	if u == sink {
		return true
	}
	for _, edge := range network[u] {
		// Follow only edges that have flow remaining.
		if edge.flow > 0 {
			// Decrement the flow to mark this edge as used.
			edge.flow--
			if dfsExtract(network, edge.to, sink, path) {
				return true
			}
			// Backtrack if the path did not lead to sink.
		}
	}
	// Remove u from path if no valid continuation exists.
	*path = (*path)[:len(*path)-1]
	return false
}

// ExtractPaths retrieves all vertex-disjoint paths from the flow network.
// It uses the computed flows (each unit of flow corresponds to one path),
// and converts the transformed node names back to original room names.
func ExtractPaths(network map[string][]*Edge, nodeMap map[string]string, source, sink string, flow int) {
	for i := 0; i < flow; i++ {
		var path []string
		if dfsExtract(network, source, sink, &path) {
			// Convert transformed node names to the original names.
			var originalPath []string
			for _, node := range path {
				if node == source || node == sink {
					originalPath = append(originalPath, node)
				} else {
					originalPath = append(originalPath, nodeMap[node])
				}
			}
			// Clean up consecutive duplicates that may result from the splitting.
			cleanedPath := []string{}
			for j, name := range originalPath {
				if j == 0 || name != originalPath[j-1] {
					cleanedPath = append(cleanedPath, name)
				}
			}
			maxFlowPaths = append(maxFlowPaths, cleanedPath)
		}
	}
}

// VertexDisjointPaths computes the vertex-disjoint paths using the max flow approach.
func VertexDisjointPaths() [][]string {
	// Build the flow network with vertex splitting.
	network, nodeMap := BuildFlowNetwork(start, end)
	// Compute the maximum flow. The flow value equals the number of vertex-disjoint paths.
	maxFlowValue := MaxFlow(network, start, end)
	// Extract the actual paths from the flow network.
	ExtractPaths(network, nodeMap, start, end, maxFlowValue)

	return maxFlowPaths
}
