package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {

	m := newMaze(3, rows, cols)
	m.Set(m.Start, Start)
	m.Set(m.Goal, Finish)
	addBorders(m)

	trapImg, _, err := ebitenutil.NewImageFromFile("trap.png")
	if err != nil {
		log.Fatal(err)
	}
	upImg, _, err := ebitenutil.NewImageFromFile("stairsup.png")
	if err != nil {
		log.Fatal(err)
	}
	downImg, _, err := ebitenutil.NewImageFromFile("stairsdown.png")
	if err != nil {
		log.Fatal(err)
	}

	sim := &Simulation{
		m:             m,
		mode:          3,
		floor:         0,
		extraPercent:  0.0,
		algo:          AlgoAStar,
		imgTrap:       trapImg,
		imgStairsUp:   upImg,
		imgStairsDown: downImg,
	}

	w, h := sim.Layout(0, 0)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Maze simulation")

	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
