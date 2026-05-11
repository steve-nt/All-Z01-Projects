package services

import "lem-in/models"

type AntAllocationService struct {
	ants  []*models.Ant
	paths []*models.Path
	alloc []int
}

func NewAntAllocationService(ants []*models.Ant, paths []*models.Path, alloc []int) *AntAllocationService {
	return &AntAllocationService{ants, paths, alloc}
}

func (s *AntAllocationService) Allocate() {
	allocIndex := 0
	for _, ant := range s.ants {
		for s.alloc[allocIndex] == 0 {
			allocIndex++
			allocIndex %= len(s.alloc)
		}
		path := s.paths[allocIndex]
		ant.Path = path

		s.alloc[allocIndex]--
		allocIndex++
		allocIndex %= len(s.alloc)
	}
}
