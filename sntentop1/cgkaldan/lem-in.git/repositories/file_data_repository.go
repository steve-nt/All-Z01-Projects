package repositories

import (
	"bufio"
	"fmt"
	"lem-in/models"
	"os"
	"strconv"
	"strings"
)

type FileDataRepository struct {
	filename string
}

func NewFileDataRepository(filename string) *FileDataRepository {
	return &FileDataRepository{filename}
}

func (r *FileDataRepository) FetchData() (int, []*models.Room, error) {
	lines, err := r.readFile()
	if err != nil {
		return 0, nil, err
	}

	numOfAnts, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, nil, err
	}
	lines = lines[1:]

	rooms := []*models.Room{}
	tunnels := []string{}
	roomsStrings := []string{}
	tunnelsStart := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case line == "##start":
			startRoom, err := processSpecialRoom(lines, &i)
			if err != nil {
				return 0, nil, err
			}
			rooms = append(rooms, startRoom)

		case line == "##end":
			endRoom, err := processSpecialRoom(lines, &i)
			if err != nil {
				return 0, nil, err
			}
			rooms = append(rooms, endRoom)

		case strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##"):
		case strings.Contains(line, "-"):
			tunnelsStart = true
			tunnels = append(tunnels, line)

		case line == "#rooms":
			// Skip this line
		default:
			if !tunnelsStart {
				roomsStrings = append(roomsStrings, line)
			}
		}
	}

	for _, roomString := range roomsStrings {
		room, err := getRoomData(roomString)
		if err != nil {
			return 0, nil, err
		}
		rooms = append(rooms, room)
	}

	rooms, err = assignTunnels(rooms, tunnels)
	if err != nil {
		return 0, nil, err
	}

	return numOfAnts, rooms, nil
}

func processSpecialRoom(lines []string, index *int) (*models.Room, error) {
	roomData := lines[*index+1]
	room, err := getRoomData(roomData)
	if err != nil {
		return nil, err
	}
	*index++ // Move to the next line after processing
	return room, nil
}

func (r *FileDataRepository) readFile() ([]string, error) {
	file, err := os.Open(r.filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getRoomData(roomString string) (*models.Room, error) {
	roomData := strings.Split(roomString, " ")
	name := roomData[0]
	if len(roomData) < 3 {
		return nil, fmt.Errorf("invalid room data: %s", roomString)
	}
	x, err := strconv.Atoi(roomData[1])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(roomData[2])
	if err != nil {
		return nil, err
	}

	room := models.NewRoom(name, float32(x), float32(y))
	return room, nil
}

func assignTunnels(rooms []*models.Room, tunnels []string) ([]*models.Room, error) {
	for _, tunnel := range tunnels {
		tunnelData := strings.Split(tunnel, "-")
		firstRoomName := tunnelData[0]
		if !roomExists(rooms, firstRoomName) {
			return nil, fmt.Errorf("room %s not found", firstRoomName)
		}
		secondRoomName := tunnelData[1]
		if !roomExists(rooms, secondRoomName) {
			return nil, fmt.Errorf("room %s not found", secondRoomName)
		}
		firstRoom := getRoomByName(rooms, firstRoomName)
		secondRoom := getRoomByName(rooms, secondRoomName)
		firstRoom.Links = append(firstRoom.Links, secondRoom)
		secondRoom.Links = append(secondRoom.Links, firstRoom)
	}
	return rooms, nil
}

func roomExists(rooms []*models.Room, name string) bool {
	for _, room := range rooms {
		if room.Name == name {
			return true
		}
	}
	return false
}

func getRoomByName(rooms []*models.Room, name string) *models.Room {
	for _, room := range rooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}
