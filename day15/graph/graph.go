package graph

import "github.com/mbordner/advent_of_code_2019/day15/geom"

var (
	Nodes = make(map[geom.Pos]*Node)
)

type Node struct {
	Pos         geom.Pos
	North       *Node
	West        *Node
	South       *Node
	East        *Node
	Traversable bool
}

func NewNode(x, y int) (n *Node) {
	pos := geom.Pos{X: x, Y: y}

	if o, ok := Nodes[pos]; ok {
		n = o
	} else {
		n = new(Node)
		n.Pos = pos
		Nodes[pos] = n
	}

	return n
}

func (n *Node) Unlink(o *Node) {
	if n.North == o {
		n.North = nil
	} else if n.East == o {
		n.East = nil
	} else if n.South == o {
		n.South = nil
	} else if n.West == o {
		n.West = nil
	}
}

func (n *Node) GetDirectionTo(o *Node) geom.Direction {
	if n.North == o {
		return geom.North
	}
	if n.East == o {
		return geom.East
	}
	if n.South == o {
		return geom.South
	}
	if n.West == o {
		return geom.West
	}
	return geom.Unknown
}

func (n *Node) SetTraversable(b bool) {
	n.Traversable = b
}

func (n *Node) GetNextNode() (*Node, geom.Direction) {
	var o *Node
	var dir geom.Direction

	if n.North == nil {
		o = NewNode(n.Pos.X, n.Pos.Y-1)
		n.North = o
		dir = geom.North
		o.South = n
	} else if n.West == nil {
		o = NewNode(n.Pos.X-1, n.Pos.Y)
		n.West = o
		dir = geom.West
		o.East = n
	} else if n.South == nil {
		o = NewNode(n.Pos.X, n.Pos.Y+1)
		n.South = o
		dir = geom.South
		o.North = n
	} else if n.East == nil {
		o = NewNode(n.Pos.X+1, n.Pos.Y)
		n.East = o
		dir = geom.East
		o.West = n
	}

	return o, dir
}

func (n *Node) ShortestPath(g *Node) []*Node {

	if n == g {
		return []*Node{n}
	}

	t := []*Node{}

	n.SetTraversable(false)

	if n.North != nil && n.North.Traversable {
		u := n.North.ShortestPath(g)
		if len(t) == 0 || (len(u) > 0 && len(u) < len(t)) {
			t = u
		}
	}

	if n.South != nil && n.South.Traversable {
		u := n.South.ShortestPath(g)
		if len(t) == 0 || (len(u) > 0 && len(u) < len(t)) {
			t = u
		}
	}

	if n.East != nil && n.East.Traversable {
		u := n.East.ShortestPath(g)
		if len(t) == 0 || (len(u) > 0 && len(u) < len(t)) {
			t = u
		}
	}

	if n.West != nil && n.West.Traversable {
		u := n.West.ShortestPath(g)
		if len(t) == 0 || (len(u) > 0 && len(u) < len(t)) {
			t = u
		}
	}

	n.SetTraversable(true)

	if len(t) > 0 {
		return append([]*Node{n}, t...)
	}
	return []*Node{}
}

type Graph struct {
	Start    *Node
	Pointer  *Node
	Goal     *Node
	previous []*Node
}

func (g *Graph) GetStart() *Node {
	return g.Start
}

func (g *Graph) GetGoal() *Node {
	return g.Goal
}

func (g *Graph) GetPointer() *Node {
	return g.Pointer
}

func (g *Graph) GenerateShortestPath() []*Node {
	p := g.Start.ShortestPath(g.Goal)
	return p
}

func (g *Graph) RemoveImpassable() {
	for p, n := range Nodes {
		if n.Traversable == false {
			delete(Nodes, p)
			if n.North != nil {
				n.North.Unlink(n)
				n.North = nil
			}
			if n.South != nil {
				n.South.Unlink(n)
				n.South = nil
			}
			if n.East != nil {
				n.East.Unlink(n)
				n.East = nil
			}
			if n.West != nil {
				n.West.Unlink(n)
				n.West = nil
			}
		}
	}
}

func NewGraph() *Graph {
	g := new(Graph)
	g.Start = NewNode(0, 0)
	g.Pointer = g.Start
	g.previous = make([]*Node, 0, 100)
	g.previous = append(g.previous, g.Start)
	g.SetTraversable(true)
	return g
}

func (g *Graph) SetGoal() {
	g.Goal = g.Pointer
}

func (g *Graph) GetNextDirection() geom.Direction {
	n, d := g.Pointer.GetNextNode()
	l := len(g.previous)
	if n == nil && l > 1 {
		prevNode := g.previous[l-1]
		g.previous = g.previous[:l-1]
		d = g.Pointer.GetDirectionTo(prevNode)
		g.Pointer = prevNode
	} else if d != geom.Unknown {
		g.previous = append(g.previous, g.Pointer)
		g.Pointer = n
	}

	return d
}

func (g *Graph) SetTraversable(b bool) {
	g.Pointer.SetTraversable(b)
	if b == false {
		l := len(g.previous)
		prevNode := g.previous[l-1]
		g.previous = g.previous[:l-1]
		g.Pointer = prevNode
	}
}
