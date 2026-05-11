package colony

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Room struct {
	Name        string
	X, Y        int
	Connections []*Room
}

type Colony struct {
	Rooms    map[string]*Room
	Start    *Room
	End      *Room
	AntCount int
}

func NewColony(input string) (*Colony, error) {
	c := &Colony{Rooms: make(map[string]*Room)}
	if err := c.parse(input); err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format: %w", err)
	}
	return c, nil
}

func (c *Colony) parse(input string) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return errors.New("empty input")
	}

	if err := c.parseAnts(lines[0]); err != nil {
		return err
	}

	var expectStart, expectEnd bool
	coordinates := make(map[string]struct{})

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case line == "##start":
			if c.Start != nil {
				return errors.New("multiple start rooms")
			}
			expectStart = true
		case line == "##end":
			if c.End != nil {
				return errors.New("multiple end rooms")
			}
			expectEnd = true
		case strings.HasPrefix(line, "#"):
			continue
		case strings.Contains(line, "-"):
			if err := c.parseTunnel(line); err != nil {
				return err
			}
		default:
			if err := c.parseRoom(line, &expectStart, &expectEnd, coordinates); err != nil {
				return err
			}
		}
	}

	if c.Start == nil || c.End == nil {
		return errors.New("missing start or end room")
	}
	return nil
}

func (c *Colony) parseAnts(line string) error {
	antCount, err := strconv.Atoi(line)
	if err != nil || antCount < 1 {
		return fmt.Errorf("invalid ant count: %q", line)
	}
	c.AntCount = antCount
	return nil
}

func (c *Colony) parseRoom(line string, expectStart, expectEnd *bool, coords map[string]struct{}) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid room format: %q", line)
	}

	name, xStr, yStr := parts[0], parts[1], parts[2]

	if invalidName(name) {
		return fmt.Errorf("invalid room name: %q", name)
	}

	if _, exists := c.Rooms[name]; exists {
		return fmt.Errorf("duplicate room: %q", name)
	}

	x, err := strconv.Atoi(xStr)
	if err != nil {
		return fmt.Errorf("invalid X coordinate: %q", xStr)
	}

	y, err := strconv.Atoi(yStr)
	if err != nil {
		return fmt.Errorf("invalid Y coordinate: %q", yStr)
	}

	coordKey := fmt.Sprintf("%d-%d", x, y)
	if _, exists := coords[coordKey]; exists {
		return fmt.Errorf("duplicate coordinates: %s", coordKey)
	}
	coords[coordKey] = struct{}{}

	room := &Room{Name: name, X: x, Y: y}
	c.Rooms[name] = room

	switch {
	case *expectStart:
		c.Start = room
		*expectStart = false
	case *expectEnd:
		c.End = room
		*expectEnd = false
	}

	return nil
}

func (c *Colony) parseTunnel(line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid tunnel format: %q", line)
	}

	a, b := parts[0], parts[1]
	if a == b {
		return fmt.Errorf("self-linked room: %q", a)
	}

	roomA, existsA := c.Rooms[a]
	roomB, existsB := c.Rooms[b]
	if !existsA || !existsB {
		return fmt.Errorf("undefined rooms in tunnel: %q-%q", a, b)
	}

	if slices.Contains(roomA.Connections, roomB) {
		return fmt.Errorf("duplicate tunnel: %q-%q", a, b)
	}

	roomA.Connections = append(roomA.Connections, roomB)
	roomB.Connections = append(roomB.Connections, roomA)
	return nil
}

func invalidName(name string) bool {
	return strings.HasPrefix(name, "L") ||
		strings.HasPrefix(name, "#") ||
		strings.Contains(name, " ")
}
