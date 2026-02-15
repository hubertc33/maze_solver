package main

import "math/rand"

func addBorders(m *Maze) {
	for z := 0; z < m.F; z++ {
		for y := 0; y < m.H; y++ {
			for x := 0; x < m.W; x++ {
				if x == 0 || y == 0 || x == m.W-1 || y == m.H-1 {
					m.Set(Pos{x, y, z}, Wall)
				}
			}
		}
	}
}

func fill(m *Maze, c Cell) {
	for i := range m.G {
		m.G[i] = c
	}
}

func generateFloors(m *Maze) {
	fill(m, Wall)
	addBorders(m)

	for z := 0; z < m.F; z++ {
		start := Pos{1, 1, z}
		m.Set(start, Floor)
		carveDFS(m, start)
	}
}

func carveDFS(m *Maze, p Pos) {
	dirs := []Pos{{2, 0, 0}, {-2, 0, 0}, {0, 2, 0}, {0, -2, 0}}
	rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })

	for _, d := range dirs {
		n := Pos{p.X + d.X, p.Y + d.Y, p.Z}
		if !m.In(n) {
			continue
		}
		if m.At(n) != Wall {
			continue
		}
		mid := Pos{p.X + d.X/2, p.Y + d.Y/2, p.Z}
		m.Set(mid, Floor)
		m.Set(n, Floor)
		carveDFS(m, n)
	}
}

func addExtraConnections(m *Maze, prob float64) {
	for z := 0; z < m.F; z++ {
		for y := 1; y < m.H-1; y++ {
			for x := 1; x < m.W-1; x++ {
				p := Pos{x, y, z}
				if m.At(p) != Wall || rand.Float64() >= prob {
					continue
				}
				ver := m.At(Pos{x, y - 1, z}) == Floor && m.At(Pos{x, y + 1, z}) == Floor
				hor := m.At(Pos{x - 1, y, z}) == Floor && m.At(Pos{x + 1, y, z}) == Floor
				if ver || hor {
					m.Set(p, Floor)
				}
			}
		}
	}
}

func addTraps(m *Maze, perFloor int) {
	for z := 0; z < m.F; z++ {
		for i := 0; i < perFloor; i++ {
			for tries := 0; tries < 10_000; tries++ {
				x := rand.Intn(m.W-2) + 1
				y := rand.Intn(m.H-2) + 1
				p := Pos{x, y, z}
				if m.At(p) == Floor {
					m.Set(p, Trap)
					break
				}
			}
		}
	}
}

func connectFloors(m *Maze, stairsPerLevel int) {
	for z := 0; z < m.F-1; z++ {
		for i := 0; i < stairsPerLevel; i++ {
			for tries := 0; tries < 10_000; tries++ {
				x := rand.Intn(m.W-2) + 1
				y := rand.Intn(m.H-2) + 1
				a := Pos{x, y, z}
				b := Pos{x, y, z + 1}
				if m.At(a) == Floor && m.At(b) == Floor {
					m.Set(a, Up)
					m.Set(b, Down)
					break
				}
			}
		}
	}
}
