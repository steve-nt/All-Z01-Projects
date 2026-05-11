package services

import (
	"lem-in/models"
	"sort"
)

type PathService struct {
	start *models.Room
	end   *models.Room
	paths []*models.Path
}

func NewPathService(start, end *models.Room) *PathService {
	return &PathService{start, end, []*models.Path{}}
}

func (p *PathService) FindAllPaths() []*models.Path {
	visited := make(map[*models.Room]bool)
	currentPath := []*models.Room{p.start}

	p.dfs(p.start, p.end, visited, currentPath)
	sort.Slice(p.paths, func(i, j int) bool {
		return len(p.paths[i].Rooms) < len(p.paths[j].Rooms)
	})
	return p.paths
}

// dfs performs depth-first search to find all paths
func (p *PathService) dfs(current, end *models.Room, visited map[*models.Room]bool, currentPath []*models.Room) {
	// Mark current room as visited
	visited[current] = true

	// If we reached the end room, we found a path
	if current == end {
		// Create a new path from the current path
		newPath := make([]*models.Room, len(currentPath))
		copy(newPath, currentPath)
		p.paths = append(p.paths, models.NewPath(newPath))
	} else {
		// Explore all neighbors
		for _, neighbor := range current.Links {
			if !visited[neighbor] {
				// Add neighbor to current path
				currentPath = append(currentPath, neighbor)
				// Recursive DFS call
				p.dfs(neighbor, end, visited, currentPath)
				// Backtrack: remove neighbor from current path
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}

	// Unmark current room (backtracking)
	visited[current] = false
}
