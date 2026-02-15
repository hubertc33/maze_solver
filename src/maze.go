package main

type Cell uint8

const (
	cellSize = 40
	uiHeight = 40
	cols     = 23
	rows     = 13
)

const (
	Floor Cell = iota
	Start
	Finish
	Wall
	Up
	Down
	Trap
)

type Pos struct{ X, Y, Z int }

type Maze struct {
	W, H, F int
	G       []Cell

	Start Pos
	Goal  Pos
}

func newMaze(floors, h, w int) *Maze {
	return &Maze{
		W: w, H: h, F: floors,
		G:     make([]Cell, floors*h*w),
		Start: Pos{1, 1, 0},
		Goal:  Pos{w - 2, h - 2, floors - 1},
	}
}

func (m *Maze) Id(p Pos) int {
	return (p.Z*m.H+p.Y)*m.W + p.X
}
func (m *Maze) In(p Pos) bool {
	return p.X >= 0 && p.X < m.W && p.Y >= 0 && p.Y < m.H && p.Z >= 0 && p.Z < m.F
}
func (m *Maze) At(p Pos) Cell {
	return m.G[m.Id(p)]
}
func (m *Maze) Set(p Pos, c Cell) {
	m.G[m.Id(p)] = c
}

func (m *Maze) Walkable(p Pos) bool {
	return m.At(p) != Wall
}

func (m *Maze) Cost(a, b Pos) int {
	if a.Z != b.Z {
		return 3
	}
	if m.At(b) == Trap {
		return 6
	}
	return 1
}

func (m *Maze) Neighbors(p Pos) []Pos {
	out := make([]Pos, 0, 6)

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, d := range dirs {
		n := Pos{X: p.X + d[0], Y: p.Y + d[1], Z: p.Z}
		if m.In(n) && m.Walkable(n) {
			out = append(out, n)
		}
	}

	switch m.At(p) {
	case Up:
		n := Pos{X: p.X, Y: p.Y, Z: p.Z + 1}
		if m.In(n) && m.Walkable(n) {
			out = append(out, n)
		}
	case Down:
		n := Pos{X: p.X, Y: p.Y, Z: p.Z - 1}
		if m.In(n) && m.Walkable(n) {
			out = append(out, n)
		}
	}

	return out
}

func (m *Maze) ClearAll() {
	for i := range m.G {
		m.G[i] = Floor
	}
	addBorders(m)
	m.Set(m.Start, Start)
	m.Set(m.Goal, Finish)
}
