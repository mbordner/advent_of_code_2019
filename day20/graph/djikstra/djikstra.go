package djikstra

import (
	"github.com/mbordner/advent_of_code_2019/day20/graph"
	hp "container/heap"
	"math"
)

type NodeValue struct {
	Node         *graph.Node
	Value        float64
	PreviousNode *graph.Node
	visited      bool
}

// this map will hold all of the nodes, visited and unvisited
type ShortestPaths map[interface{}]*NodeValue

func (sps ShortestPaths) GetShortestPath(n *graph.Node) ([]*graph.Node, float64) {

	current := sps[n.GetID()].Node
	value := float64(0)
	if sps[current.GetID()].PreviousNode != nil {
		value = sps[current.GetID()].Value
	}

	nodes := make([]*graph.Node, 0, 50)

	for sps[current.GetID()].PreviousNode != nil {
		if current.IsTraversable() == false {
			return []*graph.Node{}, float64(0)
		}
		nodes = append(nodes, current)
		current = sps[current.GetID()].PreviousNode
	}

	// reverse the array
	for i, j, h := 0, len(nodes)-1, len(nodes)/2; i < h; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	return nodes, value
}

// this will hold all of the unvisited node values sorted with minimum values at the top
type nodeValues []*NodeValue

func (h nodeValues) Len() int {
	return len(h)
}
func (h nodeValues) Less(i, j int) bool {
	return h[i].Value < h[j].Value
}
func (h nodeValues) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *nodeValues) Push(nv interface{}) {
	*h = append(*h, nv.(*NodeValue))
}

func (h *nodeValues) Pop() interface{} {
	nv := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return nv
}

type heap struct {
	values nodeValues
}

func newNodeValueHeap(capacity int) *heap {
	h := new(heap)
	h.values = make(nodeValues, 0, capacity)
	return h
}

func (h *heap) index(nv *NodeValue) int {
	for i := range h.values {
		if nv == h.values[i] {
			return i
		}
	}
	return -1
}

func (h *heap) remove(nv *NodeValue) {
	i := h.index(nv)
	if i != -1 {
		hp.Remove(&h.values, i)
	}
}

func (h *heap) fix(nv *NodeValue) {
	i := h.index(nv)
	if i != -1 {
		hp.Fix(&h.values, i)
	}
}

func (h *heap) push(nv *NodeValue) {
	hp.Push(&h.values, nv)
}

func (h *heap) pop() *NodeValue {
	i := hp.Pop(&h.values)
	return i.(*NodeValue)
}

func GenerateShortestPaths(g *graph.Graph, source *graph.Node) ShortestPaths {
	// shortest paths from n to all other nodes
	sps := make(ShortestPaths)

	// node value heap used to sort current distances through nodes
	nvh := newNodeValueHeap(g.Len())

	for _, node := range g.GetTraversableNodes() {
		nv := &NodeValue{Node: node, Value: math.MaxFloat64, PreviousNode: nil}
		if node == source {
			// this is our source node, and we need to treat it differently
			nv.Value = float64(0)
		}
		// add to the heap
		nvh.push(nv)
		sps[nv.Node.GetID()] = nv
	}

	// at this point, all unvisited nodes are in the heap

	for nvh.values.Len() > 0 {
		current := nvh.pop()

		for _, e := range current.Node.GetTraversableEdges() {

			// we don't want to explore edge destinations that already have been visited, i.e. removed from nvh
			if sps[e.GetDestination().GetID()].visited == false {

				// get the edge node value (env) for this edge's destination node
				env := sps[e.GetDestination().GetID()]

				// this value is the cost up to current node + cost to destination from current
				value := current.Value + e.GetValue()

				// check if this new value is less than anything we found before
				if value < env.Value {
					// we found a shorter path from source -> e.destination through current

					env.Value = value
					env.PreviousNode = current.Node
					// need to reorder the heap after this change
					nvh.fix(env)
				}

			}

		}

		// at this point, current is marked as visited, and will remain removed
		// we are marking it visited so we don't explore it again from another node's edges.
		sps[current.Node.GetID()].visited = true
	}

	return sps
}
