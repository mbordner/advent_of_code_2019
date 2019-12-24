package part2

import (
	"github.com/mbordner/advent_of_code_2019/day20/geom"
	"github.com/mbordner/advent_of_code_2019/day20/graph"
	"github.com/mbordner/advent_of_code_2019/day20/graph/djikstra"
	"errors"
)

type ObjectType int

const (
	Empty ObjectType = iota
	Letter
	Path
	Wall
	Portal
	HalfPortal
)

type PortalData struct {
	ID    string
	Nodes [2]*graph.Node
}

func (pd *PortalData) GetPathNode() *graph.Node {
	for _, e := range pd.Nodes[0].GetEdges() {
		ot := e.GetDestination().GetProperty("type").(ObjectType)
		if ot == Path {
			return e.GetDestination()
		}
	}
	return nil
}

func (pd *PortalData) PairNode(n *graph.Node) {
	pd.Nodes[1] = n
	n.AddProperty("pid", pd.ID)
	for i := range pd.Nodes {
		pd.Nodes[i].AddProperty("type", Portal)
	}

	pd.Nodes[0].AddEdge(pd.Nodes[1], float64(0))
	pd.Nodes[1].AddEdge(pd.Nodes[0], float64(0))
}

func NewPortalData(id string, n *graph.Node) *PortalData {
	pd := new(PortalData)

	pd.ID = id
	n.AddProperty("type", HalfPortal)
	n.AddProperty("pid", id)
	pd.Nodes[0] = n

	return pd
}

type Game struct {
	chars     [][]byte
	GameGraph *graph.Graph
	portals   map[string]*PortalData
}

func NewGame(chars [][]byte) *Game {
	g := new(Game)
	g.chars = chars
	g.GameGraph = graph.NewGraph()
	g.portals = make(map[string]*PortalData)
	g.init()

	return g
}

func (g *Game) GetPortalData(p string) *PortalData {
	if pd, ok := g.portals[p]; ok {
		return pd
	}
	return nil
}

func (g *Game) ShortestPath(p1, p2 string) ([]*graph.Node, int) {
	pd1 := g.GetPortalData(p1)
	pd2 := g.GetPortalData(p2)

	start := pd1.GetPathNode()
	end := pd2.GetPathNode()

	shortestPaths := djikstra.GenerateShortestPaths(g.GameGraph, start)
	path, distance := shortestPaths.GetShortestPath(end)

	return path, int(distance)
}

func (g *Game) init() {
	for y, row := range g.chars {
		for x, char := range row {
			pos := geom.Pos{X: x, Y: y}

			var objType ObjectType

			if char >= 'A' && char <= 'Z' {
				objType = Letter
			} else if char == '#' {
				objType = Wall
			} else if char == '.' {
				objType = Path
			}

			if objType == Letter || objType == Path {
				n := g.GameGraph.CreateNode(pos)
				n.AddProperty("value", char)
				n.AddProperty("type", objType)
			}

		}
	}

	nodes := g.GameGraph.GetNodes()

	for _, n := range nodes {

		objType := n.GetProperty("type").(ObjectType)
		pos := n.GetID().(geom.Pos)

		if objType == Path {
			var ot ObjectType

			// check for right node
			o := g.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y})
			ot, o = g.transform(geom.East, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for left node
			o = g.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y})
			ot, o = g.transform(geom.West, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for node above
			o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1})
			ot, o = g.transform(geom.North, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for node below
			o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1})
			ot, o = g.transform(geom.South, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))
				} else {
					n.AddEdge(o, float64(1))
				}
			}

		}

	}
}

func (g *Game) transform(dir geom.Direction, n *graph.Node) (ObjectType, *graph.Node) {
	var objType ObjectType

	if n != nil {
		objType = n.GetProperty("type").(ObjectType)
		if objType == Letter {

			portalId := make([]byte, 2, 2)
			pos := n.GetID().(geom.Pos)
			nb := n.GetProperty("value").(byte)

			var o *graph.Node

			switch dir {
			case geom.North:
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1})
			case geom.East:
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y})
			case geom.South:
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1})
			case geom.West:
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y})
			}

			if o == nil {
				panic(errors.New("missing expected node"))
			}
			if o.GetProperty("type").(ObjectType) != Letter {
				panic(errors.New("unexpected node type"))
			}
			ob := o.GetProperty("value").(byte)

			switch dir {
			case geom.North:
				fallthrough
			case geom.West:
				portalId[0] = ob
				portalId[1] = nb
			case geom.East:
				fallthrough
			case geom.South:
				portalId[0] = nb
				portalId[1] = ob
			}

			pid := string(portalId)

			var portalData *PortalData
			if pd, ok := g.portals[pid]; ok {
				portalData = pd
				portalData.PairNode(n)
			} else {
				portalData = NewPortalData(pid, n)
				g.portals[pid] = portalData
			}

			objType = n.GetProperty("type").(ObjectType)

		}
	}
	return objType, n
}
