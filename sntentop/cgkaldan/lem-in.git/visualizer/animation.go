package visualizer

import (
	"time"

	"github.com/g3n/engine/math32"
)

func (v *Visualizer) UpdateAnts(delta float32) {
	if !v.animating || v.currentStep >= len(v.solution) {
		return
	}

	v.roundTimer += delta
	currentRoundData := v.solution[v.currentStep]

	allArived := true
	for antID, targetRoom := range currentRoundData.AntPositions {
		am, exists := v.ants[antID]
		if !exists {
			continue
		}

		currentPos := am.graphic.Position()
		targetPos := math32.NewVector3(targetRoom.PosX, targetRoom.PosY, 0.1)

		progress := v.roundTimer / v.roundDuration
		// if progress > 1 {
		// 	progress = 1
		// }

		newX := currentPos.X + (targetPos.X-currentPos.X)*progress
		newY := currentPos.Y + (targetPos.Y-currentPos.Y)*progress

		bounce := math32.Sin(float32(time.Since(am.startTime).Milliseconds())*0.001) * 0.5
		am.graphic.SetPosition(newX, newY, 0.1+bounce)

		dx := targetPos.X - currentPos.X
		dy := targetPos.Y - currentPos.Y
		angle := math32.Atan2(dy, dx)
		am.graphic.SetRotationZ(angle - math32.Pi/2)

		if progress < 1 {
			allArived = false
		}

		if allArived {
			v.currentStep++
			v.roundTimer = 0
			v.animating = v.currentStep < len(v.solution)
		}
	}
}

func (v *Visualizer) StartAnimation() {
	if v.currentStep == 0 && v.roundTimer == 0 && !v.animating {
		v.currentStep++
	}
	v.animating = true
	v.roundTimer = 0
}

func (v *Visualizer) StopAnimation() {
	v.animating = false
}

func (v *Visualizer) NextRound() {
	if v.currentStep < len(v.solution)-1 && v.animating {
		v.SnapToRound(v.currentStep)
		v.currentStep++
		v.roundTimer = 0
	}
}

func (v *Visualizer) PrevRound() {
	if v.currentStep > 0 && v.animating {
		v.SnapToRound(v.currentStep)
		v.currentStep--
		v.roundTimer = 0
	}
}

func (v *Visualizer) SnapToRound(round int) {
	if round < 0 || round >= len(v.solution) || !v.animating {
		return
	}

	roundData := v.solution[round]
	for antID, targetRoom := range roundData.AntPositions {
		am, exists := v.ants[antID]
		if !exists {
			continue
		}
		curPos := am.graphic.Position()
		tx := targetRoom.PosX
		ty := targetRoom.PosY

		if tx != curPos.X || ty != curPos.Y {
			am.graphic.SetPosition(tx, ty, 0.1)
		}
	}
}
