package main

import (
	"container/heap"
)

type Heuristic func(a, b Pos) int

type pqItem struct {
	P        Pos
	Priority int
}
type pqueue []*pqItem

func (pq pqueue) Len() int {
	return len(pq)
}
func (pq pqueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *pqueue) Push(x any) {
	*pq = append(*pq, x.(*pqItem))
}
func (pq *pqueue) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type SearchStats struct {
	Visited  int
	Expanded int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(a, b Pos) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)*3
}

func findPath(m *Maze, start, goal Pos, h Heuristic) (path []Pos, cost int, st SearchStats) {
	pq := &pqueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{P: start, Priority: 0})

	came := map[Pos]Pos{}
	dist := map[Pos]int{start: 0}
	st.Visited = 1

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem).P
		st.Expanded++

		if cur == goal {
			return reconstruct(came, start, goal), dist[goal], st
		}

		for _, nb := range m.Neighbors(cur) {
			nd := dist[cur] + m.Cost(cur, nb)
			old, ok := dist[nb]
			if !ok || nd < old {
				if !ok {
					st.Visited++
				}
				dist[nb] = nd
				came[nb] = cur

				pr := nd
				if h != nil {
					pr += h(nb, goal)
				}
				heap.Push(pq, &pqItem{P: nb, Priority: pr})
			}
		}
	}
	return nil, 0, st
}

func BFS(m *Maze, start, goal Pos) (path []Pos, st SearchStats) {
	q := make([]Pos, 0, 256)
	q = append(q, start)
	head := 0

	seen := map[Pos]bool{start: true}
	came := map[Pos]Pos{}
	st.Visited = 1

	for head < len(q) {
		cur := q[head]
		head++
		st.Expanded++

		if cur == goal {
			return reconstruct(came, start, goal), st
		}

		for _, nb := range m.Neighbors(cur) {
			if !seen[nb] {
				seen[nb] = true
				st.Visited++
				came[nb] = cur
				q = append(q, nb)
			}
		}
	}
	return nil, st
}

func reconstruct(came map[Pos]Pos, start, goal Pos) []Pos {
	var rev []Pos
	cur := goal
	for {
		rev = append(rev, cur)
		if cur == start {
			break
		}
		cur = came[cur]
	}
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}

func pathCost(m *Maze, path []Pos) int {
	if len(path) <= 1 {
		return 0
	}
	sum := 0
	for i := 1; i < len(path); i++ {
		sum += m.Cost(path[i-1], path[i])
	}
	return sum
}
