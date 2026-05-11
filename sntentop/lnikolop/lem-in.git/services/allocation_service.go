package services

import (
	"lem-in/models"
	"sort"
)

type AllocationService struct {
	numAnts int
	paths   []*models.Path
}

func NewAllocationService(numAnts int, paths []*models.Path) *AllocationService {
	return &AllocationService{numAnts, paths}
}

func (a *AllocationService) FindBestAllocation() ([]int, int) {
	if len(a.paths) == 0 {
		return nil, 0
	}

	sortedPaths := a.createSortedPathInfo()
	optimalTime := a.findOptimalTime(sortedPaths, a.numAnts)
	initialAllocation := a.calculateInitialAllocation(sortedPaths, optimalTime)
	finalAllocation := a.adjustForSurplus(initialAllocation, a.numAnts)

	return a.mapToOriginalOrder(finalAllocation, sortedPaths), optimalTime
}

func (a *AllocationService) createSortedPathInfo() []PathInfo {
	pathInfo := make([]PathInfo, len(a.paths))
	for i, path := range a.paths {
		pathInfo[i] = PathInfo{
			originalIndex: i,
			pathLength:    len(path.Rooms),
		}
	}

	sort.Slice(pathInfo, func(i, j int) bool {
		return pathInfo[i].pathLength < pathInfo[j].pathLength
	})
	return pathInfo
}

// findOptimalTime calculates the minimum time needed for all ants to reach the end
func (a *AllocationService) findOptimalTime(sortedPaths []PathInfo, totalAnts int) int {
	optimalTime := sortedPaths[0].pathLength
	for {
		possibleAllocation := a.calculateTotalAllocation(sortedPaths, optimalTime)
		if possibleAllocation >= totalAnts {
			break
		}
		optimalTime++
	}
	return optimalTime
}

// calculateTotalAllocation calculates how many ants can be allocated in given time
func (a *AllocationService) calculateTotalAllocation(paths []PathInfo, time int) int {
	total := 0
	for _, path := range paths {
		if time >= path.pathLength {
			total += (time - path.pathLength + 1)
		}
	}
	return total
}

// calculateInitialAllocation creates initial ant allocation for each path
func (a *AllocationService) calculateInitialAllocation(paths []PathInfo, time int) []int {
	allocation := make([]int, len(paths))
	for i, path := range paths {
		if time >= path.pathLength {
			allocation[i] = time - path.pathLength + 1
		}
	}
	return allocation
}

// adjustForSurplus removes excess ants from longer paths
func (a *AllocationService) adjustForSurplus(allocation []int, totalAnts int) []int {
	currentTotal := 0
	for _, count := range allocation {
		currentTotal += count
	}

	surplus := currentTotal - totalAnts
	adjustedAllocation := make([]int, len(allocation))
	copy(adjustedAllocation, allocation)

	// Remove surplus starting from longest paths
	for i := len(adjustedAllocation) - 1; i >= 0 && surplus > 0; i-- {
		if adjustedAllocation[i] > 0 {
			if adjustedAllocation[i] <= surplus {
				surplus -= adjustedAllocation[i]
				adjustedAllocation[i] = 0
			} else {
				adjustedAllocation[i] -= surplus
				surplus = 0
			}
		}
	}
	return adjustedAllocation
}

// mapToOriginalOrder restores the original path ordering
func (a *AllocationService) mapToOriginalOrder(allocation []int, paths []PathInfo) []int {
	finalAllocation := make([]int, len(allocation))
	for i, path := range paths {
		finalAllocation[path.originalIndex] = allocation[i]
	}
	return finalAllocation
}

type PathInfo struct {
	originalIndex int
	pathLength    int
}
