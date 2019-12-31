package part2

import (
	"errors"
	"fmt"
	"github.com/mbordner/advent_of_code_2019/day20/geom"
	"github.com/mbordner/advent_of_code_2019/day20/graph"
	"github.com/mbordner/advent_of_code_2019/day20/graph/djikstra"
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

type PortalNodeType int

const (
	Inner PortalNodeType = iota
	Outer
)

type PortalData struct {
	ID    string
	Nodes [2]*graph.Node
}

func (pd *PortalData) GetPathNode(t PortalNodeType) *graph.Node {
	if n := pd.GetPortalNode(t); n != nil {
		for _, e := range n.GetEdges() {
			ot := e.GetDestination().GetProperty("type").(ObjectType)
			if ot == Path {
				return e.GetDestination()
			}
		}
	}
	return nil
}

func (pd *PortalData) PairNode(n *graph.Node) {
	pd.Nodes[1] = n
	n.AddProperty("pid", pd.ID)
	n.AddProperty("value", pd.ID)
	for i := range pd.Nodes {
		pd.Nodes[i].AddProperty("type", Portal)
	}

	pd.Nodes[0].AddEdge(pd.Nodes[1], float64(0))
	pd.Nodes[1].AddEdge(pd.Nodes[0], float64(0))
}

func (pd *PortalData) GetPortalNode(t PortalNodeType) *graph.Node {
	if pd.Nodes[0].GetProperty("portalType").(PortalNodeType) == t {
		return pd.Nodes[0]
	} else if pd.Nodes[1] != nil && pd.Nodes[1].GetProperty("portalType").(PortalNodeType) == t {
		return pd.Nodes[1]
	}
	return nil
}

func (pd *PortalData) SetPortalNodeTraversable(t PortalNodeType, b bool) {
	if n := pd.GetPortalNode(t); n != nil {
		n.SetTraversable(b)
	}
}

func (pd *PortalData) SetPortalTraversable(b bool) {
	for i := range pd.Nodes {
		if pd.Nodes[i] != nil {
			for _, e := range pd.Nodes[i].GetEdges() {
				t := e.GetDestination().GetProperty("type").(ObjectType)
				if t == Portal {
					e.SetTraversable(b)
				}
			}
		}
	}
}

func NewPortalData(id string, n *graph.Node) *PortalData {
	pd := new(PortalData)

	pd.ID = id
	n.AddProperty("type", HalfPortal)
	n.AddProperty("pid", id)
	n.AddProperty("value", id)
	n.AddProperty("portalType", Outer)
	pd.Nodes[0] = n

	return pd
}

type Level struct {
	ID        int
	GameGraph *graph.Graph
	portals   map[string]*PortalData
	bb        geom.BoundingBox
}

func NewLevel(id int, chars [][]byte) *Level {
	l := new(Level)
	l.ID = id
	l.GameGraph = graph.NewGraph()
	l.portals = make(map[string]*PortalData)
	l.init(chars)
	return l
}

func (l *Level) GetPortalData(p string) *PortalData {
	if pd, ok := l.portals[p]; ok {
		return pd
	}
	return nil
}

func (l *Level) GetPortals() []*PortalData {
	pds := make([]*PortalData, len(l.portals), len(l.portals))
	i := 0
	for _, pd := range l.portals {
		pds[i] = pd
		i++
	}
	return pds
}

func (l *Level) ShortestPath(p1, p2 string) ([]*graph.Node, int) {
	pd1 := l.GetPortalData(p1)
	pd2 := l.GetPortalData(p2)

	start := pd1.GetPathNode(Outer)
	end := pd2.GetPathNode(Outer)

	shortestPaths := djikstra.GenerateShortestPaths(l.GameGraph, start)
	path, distance := shortestPaths.GetShortestPath(end)

	return path, int(distance)
}

func (l *Level) init(chars [][]byte) {
	for y, row := range chars {
		for x, char := range row {
			pos := geom.Pos{X: x, Y: y, Z: l.ID}

			var objType ObjectType

			if char >= 'A' && char <= 'Z' {
				objType = Letter
			} else if char == '#' {
				objType = Wall
			} else if char == '.' {
				objType = Path
			}

			if objType == Letter || objType == Path {
				n := l.GameGraph.CreateNode(pos)
				n.AddProperty("value", char)
				n.AddProperty("type", objType)

				l.bb.Extend(pos)
			}

		}
	}

	nodes := l.GameGraph.GetNodes()

	for _, n := range nodes {

		objType := n.GetProperty("type").(ObjectType)
		pos := n.GetID().(geom.Pos)

		if objType == Path {
			var ot ObjectType

			// check for right node
			o := l.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y, Z: l.ID})
			ot, o = l.transform(geom.East, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))

					if l.bb.DistanceFromEdge(o.GetID().(geom.Pos)) == 1 {
						o.AddProperty("portalType", Outer)
					} else {
						o.AddProperty("portalType", Inner)
					}
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for left node
			o = l.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y, Z: l.ID})
			ot, o = l.transform(geom.West, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))

					if l.bb.DistanceFromEdge(o.GetID().(geom.Pos)) == 1 {
						o.AddProperty("portalType", Outer)
					} else {
						o.AddProperty("portalType", Inner)
					}
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for node above
			o = l.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1, Z: l.ID})
			ot, o = l.transform(geom.North, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))

					if l.bb.DistanceFromEdge(o.GetID().(geom.Pos)) == 1 {
						o.AddProperty("portalType", Outer)
					} else {
						o.AddProperty("portalType", Inner)
					}
				} else {
					n.AddEdge(o, float64(1))
				}
			}
			// check for node below
			o = l.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1, Z: l.ID})
			ot, o = l.transform(geom.South, o)
			if o != nil {
				if ot == Portal || ot == HalfPortal {
					n.AddEdge(o, float64(0))
					o.AddEdge(n, float64(1))

					if l.bb.DistanceFromEdge(o.GetID().(geom.Pos)) == 1 {
						o.AddProperty("portalType", Outer)
					} else {
						o.AddProperty("portalType", Inner)
					}
				} else {
					n.AddEdge(o, float64(1))
				}
			}

		}

	}

	pds := l.GetPortals()
	for i := range pds {
		if pds[i].ID == "AA" || pds[i].ID == "ZZ" {
			if l.ID == 0 {
				pds[i].SetPortalNodeTraversable(Outer, true)
			} else {
				pds[i].SetPortalNodeTraversable(Outer, false)
			}
		} else {
			if l.ID == 0 {
				pds[i].SetPortalNodeTraversable(Outer, false)
			} else {
				pds[i].SetPortalNodeTraversable(Outer, true)
			}
		}
	}
}

func (l *Level) transform(dir geom.Direction, n *graph.Node) (ObjectType, *graph.Node) {
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
				o = l.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1, Z: l.ID})
			case geom.East:
				o = l.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y, Z: l.ID})
			case geom.South:
				o = l.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1, Z: l.ID})
			case geom.West:
				o = l.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y, Z: l.ID})
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
			if pd, ok := l.portals[pid]; ok {
				portalData = pd
				portalData.PairNode(n)
			} else {
				portalData = NewPortalData(pid, n)
				l.portals[pid] = portalData
			}

			objType = n.GetProperty("type").(ObjectType)

		}
	}
	return objType, n
}

type Game struct {
	chars        [][]byte
	levels       []*Level
	linkedLevels map[string]int
	spsOuter     map[string]djikstra.ShortestPaths
	spsInner     map[string]djikstra.ShortestPaths
}

func NewGame(chars [][]byte) *Game {
	g := new(Game)
	g.chars = chars
	g.levels = make([]*Level, 0, 10)
	g.linkedLevels = make(map[string]int)
	g.spsOuter = make(map[string]djikstra.ShortestPaths)
	g.spsInner = make(map[string]djikstra.ShortestPaths)

	g.levels = append(g.levels, NewLevel(0, g.chars))

	pds := g.levels[0].GetPortals()
	for _, pd := range pds {
		pd.SetPortalTraversable(false)
	}
	for _, pd := range pds {
		if n := pd.GetPathNode(Inner); n != nil {
			g.spsInner[pd.ID] = djikstra.GenerateShortestPaths(g.levels[0].GameGraph, n)
		}
		if n := pd.GetPathNode(Outer); n != nil {
			g.spsOuter[pd.ID] = djikstra.GenerateShortestPaths(g.levels[0].GameGraph, n)
		}
	}

	for _, pd := range pds {
		pd.SetPortalTraversable(true)
	}

	g.linkLevels(0, "AA", 0)

	fmt.Println("merging graphs")

	for i := 1; i < len(g.levels); i++ {
		g.levels[0].GameGraph.Merge(g.levels[i].GameGraph)
	}

	fmt.Println("done merging graphs")

	return g
}

func (g *Game) linkLevels(currentLevelID int, originPortalID string, sourceLevelId int) {
	if currentLevelID == 16 && originPortalID == "ZH" {
		fmt.Println("about to fail")
	}
	originPortalType := Outer
	sourcePortalType := Inner
	sps := g.spsOuter[originPortalID]

	if currentLevelID < sourceLevelId {
		originPortalType = Inner
		sourcePortalType = Outer
		sps = g.spsInner[originPortalID]
	}

	linkId := fmt.Sprintf("%d,%d,%s", int(originPortalType), currentLevelID, originPortalID)
	//fmt.Println(sourceLevelId,"->",currentLevelID, " from ",originPortalID)

	if _, ok := g.linkedLevels[linkId]; ok {
		//fmt.Println("cycle detected, already came to ",originPortalID," from ", v, " and now coming from ",sourceLevelId)
		return
	}

	g.linkedLevels[linkId] = sourceLevelId

	level := g.GetLevel(currentLevelID)
	pds := level.GetPortals()

	if currentLevelID > len(pds)*4 {
		return
	}

	type portalLink struct {
		destinationLevel int
		portalData       *PortalData
	}

	iportals := make([]portalLink, 0, 10)
	oportals := make([]portalLink, 0, 10)

	linked := false

	for _, pd := range pds {

		if pd.ID == originPortalID {

			if currentLevelID != sourceLevelId {
				destLevelPortal := level.GetPortalData(originPortalID)
				sourceLevelPortal := g.GetLevel(sourceLevelId).GetPortalData(originPortalID)
				destLevelNode := destLevelPortal.GetPortalNode(originPortalType)
				sourceLevelNode := sourceLevelPortal.GetPortalNode(sourcePortalType)

				found := false
				for _, e := range destLevelNode.GetEdges() {
					if e.GetDestination().GetProperty("type").(ObjectType) == Portal {
						e.SetDestination(sourceLevelNode)
						found = true
						break
					}
				}
				if !found {
					panic(errors.New("why isn't this linked to a portal"))
				}
				found = false
				for _, e := range sourceLevelNode.GetEdges() {
					if e.GetDestination().GetProperty("type").(ObjectType) == Portal {
						e.SetDestination(destLevelNode)
						found = true
						break
					}
				}
				if !found {
					panic(errors.New("why isn't this linked to a portal"))
				}
			}

			linked = true
		} else {

			if currentLevelID > 0 {
				if n := pd.GetPortalNode(Outer); n != nil && n.IsTraversable() {
					if n := pd.GetPathNode(Outer); n != nil {
						pos := n.GetID().(geom.Pos)
						o := g.levels[0].GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y, Z: 0})
						_, distance := sps.GetShortestPath(o)

						if distance > 0 {
							pl := portalLink{
								destinationLevel: currentLevelID - 1,
								portalData:       pd,
							}

							oportals = append(oportals, pl)
						}
					}
				}
			}

			if n := pd.GetPortalNode(Inner); n != nil && n.IsTraversable() {
				if n := pd.GetPathNode(Inner); n != nil {
					pos := n.GetID().(geom.Pos)
					o := g.levels[0].GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y, Z: 0})
					_, distance := sps.GetShortestPath(o)

					if distance > 0 {
						pl := portalLink{
							destinationLevel: currentLevelID + 1,
							portalData:       pd,
						}

						iportals = append(iportals, pl)
					}
				}
			}

		}

	}

	if !linked {
		panic(errors.New("we didn't link"))
	}

	if len(iportals)+len(oportals) == 0 {
		//fmt.Println("won't be linking further")
	}

	for _, pl := range oportals {
		fmt.Println("linking ", originPortalID, "(", currentLevelID, ") to ", pl.portalData.ID, "(", pl.destinationLevel, ")")
		g.linkLevels(pl.destinationLevel, pl.portalData.ID, currentLevelID)
	}

	for _, pl := range iportals {
		fmt.Println("linking ", originPortalID, "(", currentLevelID, ") to ", pl.portalData.ID, "(", pl.destinationLevel, ")")
		g.linkLevels(pl.destinationLevel, pl.portalData.ID, currentLevelID)
	}

}

func (g *Game) GetLevel(id int) *Level {
	if id < len(g.levels) && g.levels[id] != nil {
		return g.levels[id]
	}
	if id != len(g.levels) {
		panic(errors.New("we are spawning new levels recursively, and we shouldn't be skipping levels"))
	}
	fmt.Println("spawned level ", id)
	l := NewLevel(id, g.chars)
	g.levels = append(g.levels, l)
	return l
}

func (g *Game) GetPortalData(p string) *PortalData {
	return g.levels[0].GetPortalData(p)
}

func (g *Game) ShortestPath(p1, p2 string) ([]*graph.Node, int) {
	return g.levels[0].ShortestPath(p1, p2)
}
