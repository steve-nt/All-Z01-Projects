package services

import (
	"lem-in/models"
)

type SimulationService struct {
	ants []*models.Ant
}

func NewSimulationService(ants []*models.Ant) *SimulationService {
	return &SimulationService{ants}
}

func (s *SimulationService) Simulate(time int) [][]string {
	var finalResult [][]string
	for i := 0; i < time; i++ {
		finalResult = append(finalResult, s.PrintSimulation(s.ants))
	}
	return finalResult
}

func (s *SimulationService) PrintSimulation(ants []*models.Ant) []string {
	tunnels := make([]*models.Tunnel, len(ants))
	for i, ant := range ants {

		tunnel := ant.PickTunnel()

		tunnels[i] = tunnel
	}
	var moves []string

	for i, ant := range ants {
		if tunnels[i] != nil {
			moves = append(moves, ant.TakeTunnel(tunnels[i]))
		}
		// if ant is finished, remove it
		if ant.Finished {
			if len(ants) > 1 {
				ants[i] = &models.Ant{
					Id:       -1,
					Finished: true,
				}
			}
		}
	}
	return moves
}
