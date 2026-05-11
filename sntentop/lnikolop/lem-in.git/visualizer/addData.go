package visualizer

import (
	"lem-in/models"
	"time"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

func (v *Visualizer) CreateRooms(rooms []*models.Room) {
	for _, room := range rooms {
		v.AddRoom(room)
	}
}

func (v *Visualizer) AddRoom(room *models.Room) {
	var geom *geometry.Geometry
	var mat *material.Standard

	switch {
	case room == v.start:
		geom = geometry.NewBox(2, 2, 2)
		mat = material.NewStandard(math32.NewColor("lime"))
	case room == v.end:
		geom = geometry.NewBox(2, 2, 2)
		mat = material.NewStandard(math32.NewColor("orange"))
	default:
		geom = geometry.NewBox(2, 2, 2)
		mat = material.NewStandard(math32.NewColor("white"))
	}

	for _, link := range room.Links {
		v.AddTunnel(room, link)
	}

	mat.SetWireframe(false)
	mat.SetOpacity(0.5)

	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(room.PosX, room.PosY, 0.1)
	v.scene.Add(mesh)
}

func (v *Visualizer) CreateAnts(numAnts int) {
	paths := v.GeneratePaths(v.start)
	for antID := 1; antID <= numAnts; antID++ {
		path, exists := paths[antID]
		if !exists || len(path) == 0 {
			continue
		}
		v.AddAnt(antID, path)
	}
}

func (v *Visualizer) AddAnt(antID int, room []*models.Room) {
	geom := geometry.NewCylinder(0.3, 1, 16, 1, true, false)
	mat := material.NewStandard(math32.NewColor("yellow"))
	mesh := graphic.NewMesh(geom, mat)

	// Create antennae
	antennaGeom := geometry.NewCylinder(0.05, 0.5, 8, 1, true, false)
	antennaMat := material.NewStandard(math32.NewColor("red"))

	// Position antennae
	leftAntenna := graphic.NewMesh(antennaGeom, antennaMat)
	leftAntenna.SetPosition(-0.2, 0.5, 0)
	leftAntenna.SetRotationZ(math32.Pi / 4)
	mesh.Add(leftAntenna)

	rightAntenna := graphic.NewMesh(antennaGeom, antennaMat)
	rightAntenna.SetPosition(0.2, 0.5, 0)
	rightAntenna.SetRotationZ(-math32.Pi / 4)
	mesh.Add(rightAntenna)

	if len(room) > 0 {
		mesh.SetPosition(room[0].PosX, room[0].PosY, 0.1)
	}

	am := &AntMesh{
		graphic: mesh,
		// speed:     2.5,
		current:   0,
		path:      room,
		startTime: time.Now(),
	}
	v.scene.Add(mesh)
	v.ants[antID] = am
}

func (v *Visualizer) AddTunnel(r1, r2 *models.Room) {
	points := []math32.Vector3{
		{X: r1.PosX, Y: r1.PosY, Z: 0.2},
		{X: r2.PosX, Y: r2.PosY, Z: 0.2},
	}

	geom := geometry.NewGeometry()
	positions := math32.NewArrayF32(0, 0)

	for _, point := range points {
		positions = append(positions, point.X, point.Y, point.Z)
	}

	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	mat := material.NewStandard(math32.NewColor("skyblue"))
	mat.SetLineWidth(3)

	line := graphic.NewLines(geom, mat)
	v.scene.Add(line)
}

// func (v *Visualizer) AddAnt(ant *models.Ant) {
// 	geom := geometry.NewCylinder(0.3, 1, 16, 1, true, false)
// 	mat := material.NewStandard(math32.NewColor("yellow"))
// 	mesh := graphic.NewMesh(geom, mat)

// 	antennaGeom := geometry.NewCylinder(0.1, 0.5, 16, 1, true, false)
// 	antennaMat := material.NewStandard(math32.NewColor("red"))

// 	leftAntenna := graphic.NewMesh(antennaGeom, antennaMat)
// 	leftAntenna.SetPosition(-0.2, 0.5, 0)
// 	leftAntenna.SetRotationZ(math32.Pi / 4)
// 	mesh.Add(leftAntenna)

// 	rightAntenna := graphic.NewMesh(antennaGeom, antennaMat)
// 	rightAntenna.SetPosition(0.2, 0.5, 0)
// 	rightAntenna.SetRotationZ(-math32.Pi / 4)
// 	mesh.Add(rightAntenna)

// 	v.scene.Add(mesh)

// 	var pathRooms []*models.Room
// 	if ant.Path != nil {
// 		pathRooms = ant.Path.Rooms
// 	}

// 	am := &AntMesh{
// 		graphic: mesh,
// 		ant:     ant,
// 		path:    pathRooms,
// 		current: 0,
// 		speed:   2.5,
// 	}

// 	if len(pathRooms) > 0 {
// 		am.graphic.SetPosition(pathRooms[0].PosX, pathRooms[0].PosY, 1.2)
// 	}

// 	v.ants = append(v.ants, am)
// }
