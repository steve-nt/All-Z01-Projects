package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/structs"
)

func ParseInputFile(filePath string) (int, []structs.Room, []structs.Tunnel, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// 1) No input at all?
	if !scanner.Scan() {
		return 0, nil, nil, errors.New("\nERROR: invalid data format\nNo input found...? ")
	}
	firstLine := strings.TrimSpace(scanner.Text())
	// 2a) Invalid ant count syntax
	antTotal, err := strconv.Atoi(firstLine)
	if err != nil {
		return 0, nil, nil, errors.New("\nERROR: invalid data format\nInvalid number in ant count entry")
	}
	// 2b) Non-positive ant count
	if antTotal <= 0 {
		return 0, nil, nil, errors.New("\nERROR: invalid data format\nInvalid ant number value given")
	}

	var (
		rooms          []structs.Room
		tunnels        []structs.Tunnel
		seenNames      = make(map[string]bool)
		seenCoords     = make(map[string]bool) // "x,y"
		seenTunnels    = make(map[string]bool) // "A-B" sorted
		startDirCount  int
		endDirCount    int
		startRoomCount int
		endRoomCount   int
		nextIsStart    bool
		nextIsEnd      bool
		prevWasDir     string // "start" or "end" or ""
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Comments / directives
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				startDirCount++
				if startDirCount > 1 {
					return 0, nil, nil, errors.New("\nERROR: invalid data format\nOnly use ##start once please")
				}
				if prevWasDir == "end" {
					return 0, nil, nil, errors.New("\nERROR: invalid data format\nDon't put ##end ##start next to eachother")
				}
				nextIsStart = true
				prevWasDir = "start"
				continue
			}
			if line == "##end" {
				endDirCount++
				if endDirCount > 1 {
					return 0, nil, nil, errors.New("\nERROR: invalid data format\nOnly use ##end once please")
				}
				if prevWasDir == "start" {
					return 0, nil, nil, errors.New("\nERROR: invalid data format\nDon't put ##end ##start next to eachother")
				}
				nextIsEnd = true
				prevWasDir = "end"
				continue
			}
			// ordinary comment
			prevWasDir = ""
			continue
		}

		// Room definition?
		parts := strings.Fields(line)
		if len(parts) == 3 {
			name, xs, ys := parts[0], parts[1], parts[2]
			// 4) Name must not start with 'L'
			if strings.HasPrefix(name, "L") {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nRoom names can't start with L, at this line: %s", line)
			}
			// 5) Duplicate room name?
			if seenNames[name] {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nDuplicate room entry found, at this line: %s", line)
			}
			seenNames[name] = true

			// 6) Parse coordinates
			x, errX := strconv.Atoi(xs)
			y, errY := strconv.Atoi(ys)
			if errX != nil || errY != nil {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nInvalid room coord number, at this line: %s", line)
			}
			coordKey := fmt.Sprintf("%d,%d", x, y)
			// 7) Duplicate coordinates?
			if seenCoords[coordKey] {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nFound multiple rooms with identical coordinates, at this line: %s", line)
			}
			seenCoords[coordKey] = true

			// Build and append
			rooms = append(rooms, structs.Room{
				Name:    name,
				X:       x,
				Y:       y,
				IsStart: nextIsStart,
				IsEnd:   nextIsEnd,
			})
			if nextIsStart {
				startRoomCount++
				if startRoomCount > 1 {
					return 0, nil, nil, errors.New(
						"\nERROR: invalid data format\nMultiple start rooms are not allowed")
				}
			}
			if nextIsEnd {
				endRoomCount++
				if endRoomCount > 1 {
					return 0, nil, nil, errors.New(
						"\nERROR: invalid data format\nMultiple end rooms are not allowed")
				}
			}

			// reset flags
			nextIsStart = false
			nextIsEnd = false
			prevWasDir = ""
			continue
		}

		// Tunnel definition?
		if strings.Contains(line, "-") {
			pair := strings.Split(line, "-")
			if len(pair) != 2 {
				return 0, nil, nil, errors.New(
					"\nERROR: invalid data format\nSomething invalid in a line of input...? at this line: " + line)
			}
			a, b := pair[0], pair[1]
			// 8) Self-loop?
			if a == b {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nCan't connect a room with itself, at this line: %s", line)
			}
			// 9) Both rooms must already exist:
			if !seenNames[a] || !seenNames[b] {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nConnection referenced a non existing room, at this line: %s", line)
			}
			// 10) Duplicate tunnel? (order-independent)
			key1, key2 := a+"-"+b, b+"-"+a
			if seenTunnels[key1] || seenTunnels[key2] {
				return 0, nil, nil, fmt.Errorf(
					"\nERROR: invalid data format\nRepeated connection found, at this line: %s", line)
			}
			seenTunnels[key1] = true

			tunnels = append(tunnels, structs.Tunnel{RoomA: a, RoomB: b})
			prevWasDir = ""
			continue
		}

		// Anything else is not valid
		return 0, nil, nil, fmt.Errorf(
			"\nERROR: invalid data format\nSomething invalid in a line of input...? at this line: %s", line)
	}

	// 11) Make sure we actually got one start and one end
	if startRoomCount == 0 {
		return 0, nil, nil, errors.New("\nERROR: invalid data format\nStart room entry missing")
	}
	if endRoomCount == 0 {
		return 0, nil, nil, errors.New("\nERROR: invalid data format\nEnd room entry missing")
	}

	return antTotal, rooms, tunnels, nil
}
