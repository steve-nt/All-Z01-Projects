package repositories

import (
	"testing"
)

func TestFileDataRepository_FetchData(t *testing.T) {
	repo := NewFileDataRepository("../example05.txt")
	numOfAnts, rooms, err := repo.FetchData()
	if err != nil {
		t.Errorf("Error fetching data: %v", err)
	}
	if numOfAnts != 9 {
		t.Errorf("Expected 9 ants, got %d", numOfAnts)
	}
	if len(rooms) != 27 {
		t.Errorf("Expected 27 rooms, got %d", len(rooms))
	}
}
