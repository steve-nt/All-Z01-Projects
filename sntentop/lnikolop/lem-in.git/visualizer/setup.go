package visualizer

import (
	"fmt"
	"lem-in/models"
	"strconv"
	"strings"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

func NewVisualizer(start, end *models.Room) *Visualizer {
	v := &Visualizer{
		app:           app.App(),
		scene:         core.NewNode(),
		start:         start,
		end:           end,
		roundDuration: 100,
		ants:          make(map[int]*AntMesh),
	}
	return v
}

// main is the entry point for the application.
func (v *Visualizer) setupScene(rooms []*models.Room, numAnts int) {
	v.app.Gls().ClearColor(0.1, 0.2, 0.3, 1.0)

	width, height := v.app.GetSize()
	aspect := float32(width) / float32(height)

	v.camera = camera.NewPerspective(aspect, 0.1, 1000, 65, camera.Horizontal)
	v.camera.SetPosition(20, 15, 50)
	v.camera.LookAt(&math32.Vector3{X: 0, Y: 0, Z: 0}, &math32.Vector3{X: 0, Y: 1, Z: 0})

	ambientLight := light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.3)
	v.scene.Add(ambientLight)

	directonalLight := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.5)
	directonalLight.SetPosition(1, 1, 1)
	v.scene.Add(directonalLight)

	camControl := camera.NewOrbitControl(v.camera)
	camControl.MaxDistance = 100
	camControl.MinDistance = 5
	camControl.KeyPanSpeed = 0
	camControl.KeyRotSpeed = 0
	camControl.RotSpeed = 0.5
	camControl.ZoomSpeed = 0.5

	// Key controls
	v.CreateRooms(rooms)
	v.CreateAnts(numAnts)
	v.setupControls()
}
func (v *Visualizer) setupControls() {
	window.Get().SubscribeID(window.OnKeyDown, v, func(evname string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		switch kev.Key {
		case window.KeyRight:
			v.NextRound()
		case window.KeyLeft:
			v.PrevRound()
		case window.KeySpace:
			if v.animating {
				v.StopAnimation()
			} else {
				v.StartAnimation()
			}
		}
	})
}

func Visualization(start, end *models.Room, rooms []*models.Room, numAnts int, solution [][]string) {
	roomMap := make(map[string]*models.Room)
	for _, room := range rooms {
		roomMap[room.Name] = room
	}
	vis := NewVisualizer(start, end)
	vis.ParseSolution(solution, roomMap)
	vis.setupScene(rooms, numAnts)
	vis.Run()
}

func (v *Visualizer) ParseSolution(solution [][]string, roomMap map[string]*models.Room) {
	v.solution = make([]*Round, len(solution)+1)
	initialRound := &Round{
		AntPositions: make(map[int]*models.Room),
	}

	antIDs := make(map[int]bool)
	for _, step := range solution {
		for _, move := range step {
			parts := strings.Split(move, "-")
			antID, _ := strconv.Atoi(parts[0][1:])
			antIDs[antID] = true
		}
	}

	for antID := range antIDs {
		initialRound.AntPositions[antID] = v.start
	}

	v.solution[0] = initialRound

	for i, step := range solution {
		round := &Round{
			AntPositions: make(map[int]*models.Room),
		}

		for _, move := range step {
			parts := strings.Split(move, "-")

			antID, _ := strconv.Atoi(parts[0][1:])
			roomName := strings.TrimSpace(parts[1])

			room, exists := roomMap[roomName]
			if !exists {
				fmt.Printf("Unknown room: %s in move: %s\n", roomName, move)
				continue
			}
			round.AntPositions[antID] = room
		}

		v.solution[i+1] = round
	}
}

func (v *Visualizer) GeneratePaths(startRoom *models.Room) map[int][]*models.Room {
	paths := make(map[int][]*models.Room)

	// Find all ant IDs
	antIDs := make(map[int]bool)
	for _, round := range v.solution {
		for antID := range round.AntPositions {
			antIDs[antID] = true
		}
	}

	// Initialize paths with start room
	for antID := range antIDs {
		paths[antID] = []*models.Room{startRoom}
	}

	// Build paths from solution
	for _, round := range v.solution {
		for antID := range antIDs {
			currentPath := paths[antID]
			lastPos := currentPath[len(currentPath)-1]

			if newRoom, exists := round.AntPositions[antID]; exists {
				paths[antID] = append(currentPath, newRoom)
			} else {
				paths[antID] = append(currentPath, lastPos)
			}
		}
	}

	return paths
}

func (v *Visualizer) Run() {
	v.app.Gls().Enable(gls.MULTISAMPLE)
	lastTime := time.Now()

	v.app.Run(func(renderer *renderer.Renderer, _ time.Duration) {
		now := time.Now()
		delta := float32(now.Sub(lastTime).Seconds())
		lastTime = now
		v.UpdateAnts(delta)
		v.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(v.scene, v.camera)
	})
}
