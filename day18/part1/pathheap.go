package part1

import (
	"github.com/mbordner/advent_of_code_2019/day18/graph"
	hp "container/heap"
)

type Path struct {
	Value    byte
	Distance int
	Nodes    []*graph.Node
}


type Paths []*Path

func (h Paths) Len() int {
	return len(h)
}
func (h Paths) Less(i, j int) bool {
	return h[i].Distance < h[j].Distance
}
func (h Paths) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Paths) Push(nv interface{}) {
	*h = append(*h, nv.(*Path))
}

func (h *Paths) Pop() interface{} {
	nv := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return nv
}

type pathHeap struct {
	paths Paths
}

func newPathHeap(capacity int) *pathHeap {
	h := new(pathHeap)
	h.paths = make(Paths, 0, capacity)
	return h
}

func (h *pathHeap) index(p *Path) int {
	for i := range h.paths {
		if p == h.paths[i] {
			return i
		}
	}
	return -1
}

func (h *pathHeap) remove(p *Path) {
	i := h.index(p)
	if i != -1 {
		hp.Remove(&h.paths, i)
	}
}

func (h *pathHeap) fix(p *Path) {
	i := h.index(p)
	if i != -1 {
		hp.Fix(&h.paths, i)
	}
}

func (h *pathHeap) push(p *Path) {
	hp.Push(&h.paths, p)
}

func (h *pathHeap) pop() *Path {
	i := hp.Pop(&h.paths)
	return i.(*Path)
}