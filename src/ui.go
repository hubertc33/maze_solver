package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Algorithm int

const (
	AlgoAStar Algorithm = iota
	AlgoDijkstra
	AlgoBFS
)

func (a Algorithm) String() string {
	switch a {
	case AlgoAStar:
		return "A*"
	case AlgoDijkstra:
		return "Dijkstra"
	default:
		return "BFS"
	}
}

type Simulation struct {
	m *Maze

	mode         int
	floor        int
	extraPercent float64
	algo         Algorithm

	path     []Pos
	pathCost int
	stats    SearchStats

	imgTrap       *ebiten.Image
	imgStairsUp   *ebiten.Image
	imgStairsDown *ebiten.Image
}

func (s *Simulation) Layout(outW, outH int) (int, int) {
	return cols * cellSize, rows*cellSize + uiHeight
}

func (s *Simulation) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		s.mode = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		s.mode = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		s.mode = 3
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		s.mode = 4
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		s.mode = 5
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		s.mode = 6
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		s.algo = (s.algo + 1) % 3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && s.floor < s.m.F-1 {
		s.floor++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && s.floor > 0 {
		s.floor--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		s.addFloor()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		s.removeFloor()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		s.m.ClearAll()
		s.path = nil
		s.pathCost = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		switch s.extraPercent {
		case 0.0:
			s.extraPercent = 0.1
		case 0.1:
			s.extraPercent = 0.2
		case 0.2:
			s.extraPercent = 0.3
		case 0.3:
			s.extraPercent = 0.4
		case 0.4:
			s.extraPercent = 0.5
		case 0.5:
			s.extraPercent = 0.6
		case 0.6:
			s.extraPercent = 0.7
		default:
			s.extraPercent = 0.0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		generateFloors(s.m)
		addExtraConnections(s.m, s.extraPercent)
		addTraps(s.m, 30)
		connectFloors(s.m, 10)

		s.m.Set(s.m.Start, Start)
		s.m.Set(s.m.Goal, Finish)

		s.path = nil
		s.pathCost = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.runSearch()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		s.path = nil
		s.pathCost = 0
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gx := x / cellSize
		gy := (y - uiHeight) / cellSize
		if gx >= 0 && gx < s.m.W && gy >= 0 && gy < s.m.H {
			s.handleClick(gx, gy)
		}
	}

	return nil
}

func (s *Simulation) runSearch() {
	start := s.m.Start
	goal := s.m.Goal

	switch s.algo {
	case AlgoBFS:
		p, st := BFS(s.m, start, goal)
		s.path = p
		s.stats = st
		if s.path != nil {
			s.pathCost = pathCost(s.m, s.path)
		}

	case AlgoDijkstra:
		p, c, st := findPath(s.m, start, goal, nil)
		s.path, s.pathCost, s.stats = p, c, st

	case AlgoAStar:
		p, c, st := findPath(s.m, start, goal, manhattan)
		s.path, s.pathCost, s.stats = p, c, st
	}
}

func (s *Simulation) handleClick(gx, gy int) {
	p := Pos{gx, gy, s.floor}
	cell := s.m.At(p)

	switch s.mode {
	case 1: // START
		if cell == Wall || cell == Finish || cell == Up || cell == Down || cell == Trap {
			return
		}
		s.m.Set(s.m.Start, Floor)
		s.m.Start = p
		s.m.Set(s.m.Start, Start)

	case 2: // FINISH
		if cell == Wall || cell == Start || cell == Up || cell == Down || cell == Trap {
			return
		}
		s.m.Set(s.m.Goal, Floor)
		s.m.Goal = p
		s.m.Set(s.m.Goal, Finish)

	case 3: // WALL
		if gx == 0 || gy == 0 || gx == s.m.W-1 || gy == s.m.H-1 {
			return
		}
		if cell == Start || cell == Finish {
			return
		}
		if cell == Wall {
			s.m.Set(p, Floor)
		} else if cell == Floor {
			s.m.Set(p, Wall)
		}

	case 4: // STAIRS UP
		if s.floor >= s.m.F-1 {
			return
		}
		if cell == Down {
			return
		}

		if cell == Up {
			s.m.Set(p, Floor)
			s.m.Set(Pos{gx, gy, s.floor + 1}, Floor)
			return
		}

		above := Pos{gx, gy, s.floor + 1}
		if s.m.At(p) != Floor || s.m.At(above) != Floor {
			return
		}
		s.m.Set(p, Up)
		s.m.Set(above, Down)

	case 5: // STAIRS DOWN
		if s.floor <= 0 {
			return
		}
		if cell == Up {
			return
		}
		if cell == Down {
			s.m.Set(p, Floor)
			s.m.Set(Pos{gx, gy, s.floor - 1}, Floor)
			return
		}
		below := Pos{gx, gy, s.floor - 1}
		if s.m.At(p) != Floor || s.m.At(below) != Floor {
			return
		}
		s.m.Set(p, Down)
		s.m.Set(below, Up)

	case 6: // TRAP
		if cell == Start || cell == Finish || cell == Wall || cell == Up || cell == Down {
			return
		}
		if cell == Trap {
			s.m.Set(p, Floor)
		} else {
			s.m.Set(p, Trap)
		}
	}
}

func (s *Simulation) addFloor() {
	old := s.m
	nm := newMaze(old.F+1, old.H, old.W)
	copy(nm.G, old.G)
	addBorders(nm)

	nm.Start = old.Start
	nm.Goal = old.Goal
	nm.Set(nm.Start, Start)
	nm.Set(nm.Goal, Finish)

	s.m = nm
}

func (s *Simulation) removeFloor() {
	if s.m.F <= 1 {
		return
	}
	old := s.m
	nm := newMaze(old.F-1, old.H, old.W)
	copy(nm.G, old.G[:(old.F-1)*old.H*old.W])
	addBorders(nm)

	nm.Start = old.Start
	if nm.Start.Z >= nm.F {
		nm.Start.Z = nm.F - 1
	}
	nm.Goal = old.Goal
	if nm.Goal.Z >= nm.F {
		nm.Goal.Z = nm.F - 1
	}
	nm.Set(nm.Start, Start)
	nm.Set(nm.Goal, Finish)

	s.m = nm
	if s.floor >= s.m.F {
		s.floor = s.m.F - 1
	}
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	screenW := cols * cellSize

	screen.Fill(color.RGBA{30, 30, 30, 255})
	vector.FillRect(screen, 0, 0, float32(screenW), float32(uiHeight), color.RGBA{15, 15, 15, 255}, false)

	modeName := map[int]string{
		1: "SET START",
		2: "SET FINISH",
		3: "SET WALL",
		4: "STAIRS UP",
		5: "STAIRS DOWN",
		6: "SET TRAP",
	}[s.mode]
	if modeName == "" {
		modeName = "SET WALL"
	}

	ebitenutil.DebugPrintAt(screen,
		"[1-6] Mods  [A] Algorithm  [^/v] Floor  [SPACE] Start  [+/-] Add/Del Floor  [M] Maze type  [G] Generate  [C] Clear  [X] Clear path",
		10, 10)

	ebitenutil.DebugPrintAt(screen,
		"Mode: "+modeName+
			" | Algorithm: "+s.algo.String()+
			" | Maze type: "+fmt.Sprintf("%.0f%%", s.extraPercent*100)+
			" | Floor: "+fmt.Sprintf("%d / %d", s.floor, s.m.F-1),
		screenW/2-300, 25)

	for y := 0; y < s.m.H; y++ {
		for x := 0; x < s.m.W; x++ {
			p := Pos{x, y, s.floor}
			cell := s.m.At(p)

			x0 := float64(x * cellSize)
			y0 := float64(y*cellSize + uiHeight)

			vector.FillRect(screen, float32(x0), float32(y0), float32(cellSize-1), float32(cellSize-1),
				color.RGBA{180, 180, 180, 255}, false)

			var img *ebiten.Image
			switch cell {
			case Trap:
				img = s.imgTrap
			case Up:
				img = s.imgStairsUp
			case Down:
				img = s.imgStairsDown
			}

			if img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(
					float64(cellSize)/float64(img.Bounds().Dx()),
					float64(cellSize)/float64(img.Bounds().Dy()),
				)
				op.GeoM.Translate(x0, y0)
				screen.DrawImage(img, op)
				continue
			}

			var clr color.Color
			switch cell {
			case Floor:
				clr = color.RGBA{180, 180, 180, 255}
			case Start:
				clr = color.RGBA{0, 200, 0, 255}
			case Finish:
				clr = color.RGBA{200, 0, 0, 255}
			case Wall:
				clr = color.RGBA{50, 50, 50, 255}
			default:
				clr = color.RGBA{255, 0, 255, 255}
			}
			vector.FillRect(screen, float32(x0), float32(y0), float32(cellSize-1), float32(cellSize-1), clr, false)
		}
	}

	if s.path != nil {
		for _, p := range s.path {
			if p.Z != s.floor {
				continue
			}
			c := s.m.At(p)
			if c == Start || c == Finish {
				continue
			}

			x0f := float32(p.X * cellSize)
			y0f := float32(p.Y*cellSize + uiHeight)

			vector.FillRect(screen, x0f, y0f, float32(cellSize-1), float32(cellSize-1),
				color.RGBA{255, 255, 0, 180}, false)

			var img *ebiten.Image
			switch c {
			case Trap:
				img = s.imgTrap
			case Up:
				img = s.imgStairsUp
			case Down:
				img = s.imgStairsDown
			}

			if img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(
					float64(cellSize)/float64(img.Bounds().Dx()),
					float64(cellSize)/float64(img.Bounds().Dy()),
				)
				op.GeoM.Translate(float64(p.X*cellSize), float64(p.Y*cellSize+uiHeight))
				screen.DrawImage(img, op)
			}
		}

		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf(
				"Path cost: %d   Steps: %d   Visited: %d   Expanded: %d",
				s.pathCost, len(s.path)-1,
				s.stats.Visited, s.stats.Expanded,
			),
			10, 55)
	}
}
