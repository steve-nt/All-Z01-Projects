package services

import "lem-in/models"

type SubsetService struct {
	paths     []*models.Path
	subsets   [][]*models.Path
	numOfAnts int
}

func NewSubsetService(paths []*models.Path, numOfAnts int) *SubsetService {
	return &SubsetService{paths, [][]*models.Path{}, numOfAnts}
}

func (s *SubsetService) GetGoodSubsets() [][]*models.Path {
	s.generateNonOverlappingSubsets()
	maxLength := 0
	for _, subset := range s.subsets {
		if len(subset) > maxLength {
			maxLength = len(subset)
		}
	}

	goodSubsets := [][]*models.Path{}

	for _, subset := range s.subsets {
		if len(subset) == maxLength {
			goodSubsets = append(goodSubsets, subset)
		}
	}
	return goodSubsets
}

func (s *SubsetService) generateNonOverlappingSubsets() {
	var result [][]*models.Path

	var backtrack func(index int, current []*models.Path)
	backtrack = func(index int, current []*models.Path) {
		// Save the current valid subset
		result = append(result, append([]*models.Path{}, current...))

		times := s.numOfAnts
		if len(s.paths) < times {
			times = len(s.paths)
		}

		for i := index; i < times; i++ {
			conflict := false
			for _, selected := range current {
				if selected.CheckCommonRooms(s.paths[i]) {
					conflict = true
					break
				}
			}
			if !conflict {
				backtrack(i+1, append(current, s.paths[i]))
			}
		}
	}

	backtrack(0, []*models.Path{})
	s.subsets = result
}
