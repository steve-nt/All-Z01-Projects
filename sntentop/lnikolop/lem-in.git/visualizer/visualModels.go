package visualizer

import (
	"lem-in/models"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
)

type Visualizer struct {
	app    *app.Application
	scene  *core.Node
	camera *camera.Camera
	// rooms         map[*models.Room]*graphic.Mesh
	ants     map[int]*AntMesh
	start    *models.Room
	end      *models.Room
	solution []*Round
	// solutionPaths map[int][]*models.Room
	currentStep int
	animating   bool
	// progress      float32
	roundTimer    float32
	roundDuration float32
}

type AntMesh struct {
	graphic *graphic.Mesh
	path    []*models.Room
	current int
	// speed     float32
	startTime time.Time
}

type Round struct {
	AntPositions map[int]*models.Room
}
