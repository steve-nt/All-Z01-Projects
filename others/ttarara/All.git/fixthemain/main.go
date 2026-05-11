package main

import (
	"github.com/01-edu/z01"
)

// Δομή Door και καταστάσεις
type Door struct {
	state int
}

const (
	OPEN  = 1
	CLOSE = 0
)

// Συνάρτηση PrintStr
func PrintStr(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
	z01.PrintRune('\n')
}

// Συνάρτηση OpenDoor
func OpenDoor(ptrDoor *Door) bool {
	PrintStr("Door Opening...")
	ptrDoor.state = OPEN
	return true
}

// Συνάρτηση CloseDoor
func CloseDoor(ptrDoor *Door) bool {
	PrintStr("Door Closing...")
	ptrDoor.state = CLOSE
	return true
}

// Συνάρτηση IsDoorOpen
func IsDoorOpen(ptrDoor *Door) bool {
	PrintStr("is the Door opened ?")
	return ptrDoor.state == OPEN
}

// Συνάρτηση IsDoorClose
func IsDoorClose(ptrDoor *Door) bool {
	PrintStr("is the Door closed ?")
	return ptrDoor.state == CLOSE
}

// Συνάρτηση main
func main() {
	door := &Door{}

	OpenDoor(door)
	if IsDoorClose(door) {
		OpenDoor(door)
	}
	if IsDoorOpen(door) {
		CloseDoor(door)
	}
	if door.state == OPEN {
		CloseDoor(door)
	}
}
