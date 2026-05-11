package models

import (
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "in-memory database",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitDB(":memory:")
			if err != nil {
				t.Errorf("InitDB failed: %v", err)
			}
		})
	}
}
