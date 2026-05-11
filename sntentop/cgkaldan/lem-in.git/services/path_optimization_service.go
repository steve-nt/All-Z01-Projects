package services

import (
	"lem-in/models"
	"math"
)

type PathOptimizationService struct {
	subsets [][]*models.Path
}

func NewPathOptimizationService(subsets [][]*models.Path) *PathOptimizationService {
	return &PathOptimizationService{subsets: subsets}
}

func (s *PathOptimizationService) Optimize(numOfAnts int) ([]*models.Path, []int, int) {
	bestSubset := s.subsets[0]
	bestTime := math.MaxInt
	bestAlloc := []int{}
	for _, subset := range s.subsets {
		allocationService := NewAllocationService(numOfAnts, subset)
		alloc, time := allocationService.FindBestAllocation()
		if time < bestTime {
			bestTime = time
			bestSubset = subset
			bestAlloc = alloc
		}
	}
	return bestSubset, bestAlloc, bestTime
}
