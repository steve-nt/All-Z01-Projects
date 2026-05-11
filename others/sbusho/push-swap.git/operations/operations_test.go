package operations

import (
	"reflect"
	"testing"
)

// Helper functions
func assertStack(t *testing.T, got, want []int, stackName string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", stackName, got, want)
	}
}

func assertStacks(t *testing.T, gotA, wantA, gotB, wantB []int) {
	t.Helper()
	assertStack(t, gotA, wantA, "stackA")
	assertStack(t, gotB, wantB, "stackB")
}

func copyStack(s []int) []int {
	return append([]int{}, s...)
}

// Push operations
func TestPushOperations(t *testing.T) {
	t.Run("Pa", func(t *testing.T) {
		tests := []struct {
			name   string
			stackA []int
			stackB []int
			wantA  []int
			wantB  []int
		}{
			{"empty stackB", []int{1}, []int{}, []int{1}, []int{}},
			{"normal case", []int{1}, []int{2}, []int{2, 1}, []int{}},
			{"empty stackA", []int{}, []int{3}, []int{3}, []int{}},
			{"multiple elements", []int{1, 2}, []int{3, 4}, []int{3, 1, 2}, []int{4}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, b := copyStack(tt.stackA), copyStack(tt.stackB)
				Pa(&a, &b)
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			})
		}
	})

	t.Run("Pb", func(t *testing.T) {
		tests := []struct {
			name   string
			stackA []int
			stackB []int
			wantA  []int
			wantB  []int
		}{
			{"empty stackA", []int{}, []int{1}, []int{}, []int{1}},
			{"normal case", []int{2}, []int{1}, []int{}, []int{2, 1}},
			{"multiple elements", []int{3, 4}, []int{1, 2}, []int{4}, []int{3, 1, 2}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, b := copyStack(tt.stackA), copyStack(tt.stackB)
				Pb(&a, &b)
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			})
		}
	})
}

// Swap operations
func TestSwapOperations(t *testing.T) {
	t.Run("Sa", func(t *testing.T) {
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"two elements", []int{2, 1}, []int{1, 2}},
			{"three elements", []int{3, 1, 2}, []int{1, 3, 2}},
			{"single element", []int{1}, []int{1}},
			{"empty stack", []int{}, []int{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Sa(&s)
				assertStack(t, s, tt.want, "stackA")
			})
		}
	})

	t.Run("Sb", func(t *testing.T) {
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"swap two", []int{1, 2}, []int{2, 1}},
			{"no swap", []int{1}, []int{1}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Sb(&s)
				assertStack(t, s, tt.want, "stackB")
			})
		}
	})

	t.Run("Ss", func(t *testing.T) {
		tests := []struct {
			name   string
			stackA []int
			stackB []int
			wantA  []int
			wantB  []int
		}{
			{"both stacks", []int{2, 1}, []int{4, 3}, []int{1, 2}, []int{3, 4}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, b := copyStack(tt.stackA), copyStack(tt.stackB)
				Ss(&a, &b)
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			})
		}
	})
}

// Rotate operations
func TestRotateOperations(t *testing.T) {
	t.Run("Ra", func(t *testing.T) {
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"three elements", []int{1, 2, 3}, []int{2, 3, 1}},
			{"two elements", []int{1, 2}, []int{2, 1}},
			{"no rotation", []int{1}, []int{1}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Ra(&s)
				assertStack(t, s, tt.want, "stackA")
			})
		}
	})

	t.Run("Rb", func(t *testing.T) {
		// Similar to Ra tests
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"normal case", []int{1, 2, 3}, []int{2, 3, 1}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Rb(&s)
				assertStack(t, s, tt.want, "stackB")
			})
		}
	})

	t.Run("Rr", func(t *testing.T) {
		tests := []struct {
			name   string
			stackA []int
			stackB []int
			wantA  []int
			wantB  []int
		}{
			{"both stacks", []int{1, 2, 3}, []int{4, 5, 6}, []int{2, 3, 1}, []int{5, 6, 4}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, b := copyStack(tt.stackA), copyStack(tt.stackB)
				Rr(&a, &b)
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			})
		}
	})
}

// Reverse rotate operations
func TestReverseRotateOperations(t *testing.T) {
	t.Run("Rra", func(t *testing.T) {
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"three elements", []int{1, 2, 3}, []int{3, 1, 2}},
			{"two elements", []int{1, 2}, []int{2, 1}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Rra(&s)
				assertStack(t, s, tt.want, "stackA")
			})
		}
	})

	t.Run("Rrb", func(t *testing.T) {
		// Similar to Rra tests
		tests := []struct {
			name  string
			stack []int
			want  []int
		}{
			{"normal case", []int{1, 2, 3}, []int{3, 1, 2}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := copyStack(tt.stack)
				Rrb(&s)
				assertStack(t, s, tt.want, "stackB")
			})
		}
	})

	t.Run("Rrr", func(t *testing.T) {
		tests := []struct {
			name   string
			stackA []int
			stackB []int
			wantA  []int
			wantB  []int
		}{
			{"both stacks", []int{1, 2, 3}, []int{4, 5, 6}, []int{3, 1, 2}, []int{6, 4, 5}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				a, b := copyStack(tt.stackA), copyStack(tt.stackB)
				Rrr(&a, &b)
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			})
		}
	})
}

// Command execution
func TestExecuteCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		initialA []int
		initialB []int
		wantA    []int
		wantB    []int
		wantErr  bool
	}{
		{"valid pa", "pa", []int{1}, []int{2}, []int{2, 1}, []int{}, false},
		{"valid pb", "pb", []int{2}, []int{1}, []int{}, []int{2, 1}, false},
		{"invalid cmd", "xx", []int{1}, []int{}, []int{1}, []int{}, true},
		{"empty cmd", "", []int{1}, []int{}, []int{1}, []int{}, false},
		{"with spaces", " pa ", []int{1}, []int{2}, []int{2, 1}, []int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := copyStack(tt.initialA), copyStack(tt.initialB)
			err := ExecuteCommand(tt.command, &a, &b)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assertStacks(t, a, tt.wantA, b, tt.wantB)
			}
		})
	}
}
