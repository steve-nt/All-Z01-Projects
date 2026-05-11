package utils

import (
	"testing"
)

func TestSortFinal(t *testing.T) {
	tests := []struct {
		name         string
		stackA       []int
		wantCommands []string
	}{
		{
			name:         "empty stack",
			stackA:       []int{},
			wantCommands: []string{},
		},
		{
			name:         "single element",
			stackA:       []int{1},
			wantCommands: []string{},
		},
		{
			name:         "two elements sorted",
			stackA:       []int{1, 2},
			wantCommands: []string{},
		},
		{
			name:         "two elements unsorted",
			stackA:       []int{2, 1},
			wantCommands: []string{"sa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stackA := make([]int, len(tt.stackA))
			copy(stackA, tt.stackA)
			stackB := []int{}

			gotCommands := SortFinal(&stackA, &stackB)

			// Compare lengths first for better error reporting
			if len(gotCommands) != len(tt.wantCommands) {
				t.Errorf("SortFinal() command count = %d, want %d", len(gotCommands), len(tt.wantCommands))
				return
			}

			// Compare contents if lengths match
			for i := range gotCommands {
				if gotCommands[i] != tt.wantCommands[i] {
					t.Errorf("command %d = %q, want %q", i, gotCommands[i], tt.wantCommands[i])
				}
			}
		})
	}
}

func TestSortThree(t *testing.T) {
	tests := []struct {
		name         string
		stackA       []int
		wantA        []int
		wantCommands []string
	}{
		{
			name:         "case 2 1 3",
			stackA:       []int{2, 1, 3},
			wantA:        []int{1, 2, 3},
			wantCommands: []string{"sa"},
		},
		{
			name:         "case 3 2 1",
			stackA:       []int{3, 2, 1},
			wantA:        []int{1, 2, 3},
			wantCommands: []string{"sa", "rra"},
		},
		{
			name:         "case 3 1 2",
			stackA:       []int{3, 1, 2},
			wantA:        []int{1, 2, 3},
			wantCommands: []string{"ra"},
		},
		{
			name:         "case 1 3 2",
			stackA:       []int{1, 3, 2},
			wantA:        []int{1, 2, 3},
			wantCommands: []string{"sa", "ra"},
		},
		{
			name:         "case 2 3 1",
			stackA:       []int{2, 3, 1},
			wantA:        []int{1, 2, 3},
			wantCommands: []string{"rra"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stackA := make([]int, len(tt.stackA))
			copy(stackA, tt.stackA)
			stackB := []int{}
			var commands []string

			SortThree(&stackA, &stackB, &commands)

			// Compare stackA contents
			if len(stackA) != len(tt.wantA) {
				t.Errorf("stackA length = %d, want %d", len(stackA), len(tt.wantA))
				return
			}
			for i := range stackA {
				if stackA[i] != tt.wantA[i] {
					t.Errorf("stackA[%d] = %d, want %d", i, stackA[i], tt.wantA[i])
				}
			}

			// Compare commands
			if len(commands) != len(tt.wantCommands) {
				t.Errorf("command count = %d, want %d", len(commands), len(tt.wantCommands))
				return
			}
			for i := range commands {
				if commands[i] != tt.wantCommands[i] {
					t.Errorf("command %d = %q, want %q", i, commands[i], tt.wantCommands[i])
				}
			}
		})
	}
}
