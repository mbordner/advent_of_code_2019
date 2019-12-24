package part2

import (
	"github.com/mbordner/advent_of_code_2019/day18/geom"
	"github.com/mbordner/advent_of_code_2019/day18/graph"
	"github.com/mbordner/advent_of_code_2019/day18/graph/djikstra"
	"errors"
	"fmt"
	"math"
	"sort"
)

type ObjectType int

const (
	Empty ObjectType = iota
	Wall
	Door
	Key
	Start
)

type DistanceCacheResults struct {
	permutation string
	distance    int
}

type DistanceCache struct {
	cache map[string]DistanceCacheResults
}

func NewDistanceCache() *DistanceCache {
	dc := new(DistanceCache)
	dc.cache = make(map[string]DistanceCacheResults)
	return dc
}

func (dc *DistanceCache) SortAcquiredKeys(keys []byte) {
	acquiredKeys := keys[:len(keys)-1]
	sort.Slice(acquiredKeys, func(i int, j int) bool { return acquiredKeys[i] < acquiredKeys[j] })
	copy(keys, acquiredKeys)
}

func (dc *DistanceCache) GetResults(locations []byte,keysTaken []byte) *DistanceCacheResults {
	keys := make([]byte,len(keysTaken),len(keysTaken))
	copy(keys,keysTaken)
	dc.SortAcquiredKeys(keys)

	if results, ok := dc.cache[string(locations)+string(keys)]; ok {
		return &results
	}

	return nil
}

func (dc *DistanceCache) CacheResults(locations []byte,keysTaken []byte, permutation string, distance int) {
	keys := make([]byte,len(keysTaken),len(keysTaken))
	copy(keys,keysTaken)
	dc.SortAcquiredKeys(keys)

	dc.cache[string(locations)+string(keys)] = DistanceCacheResults{permutation: permutation, distance: distance}
}

type KeyDistances struct {
	keyDistance map[byte]map[byte]int
	keyRequires map[byte]map[byte]map[byte]bool
}

func (kd *KeyDistances) AddPath(a byte, b byte, path []*graph.Node, distance float64) {
	if _, ok := kd.keyDistance[a]; !ok {
		kd.keyDistance[a] = make(map[byte]int)
		kd.keyRequires[a] = make(map[byte]map[byte]bool)
	}
	if _, ok := kd.keyRequires[a][b]; !ok {
		kd.keyRequires[a][b] = make(map[byte]bool)
	}
	kd.keyDistance[a][b] = int(distance)
	for i, n := range path {
		if n.GetProperty("type").(ObjectType) == Door {
			requiredKey := n.GetProperty("value").(byte)
			kd.keyRequires[a][b][requiredKey+32] = true
		} else if n.GetProperty("type").(ObjectType) == Key && i != len(path) - 1{
			requiredKey := n.GetProperty("value").(byte)
			kd.keyRequires[a][b][requiredKey] = true
		}
	}
}

func (kd *KeyDistances) GetDistance(a byte, b byte, acquiredKeys map[byte]bool) int {
	if d, ok := kd.keyDistance[a][b]; ok {
		for rk, v := range kd.keyRequires[a][b] {
			if v {
				if acquired, ok := acquiredKeys[rk]; !ok || !acquired {
					return 0
				}
			}
		}

		return d
	}

	return 0
}

func NewKeyDistances() *KeyDistances {
	kd := new(KeyDistances)
	kd.keyDistance = make(map[byte]map[byte]int)
	kd.keyRequires = make(map[byte]map[byte]map[byte]bool)
	return kd
}

type Game struct {
	GameGraph          *graph.Graph
	keys               map[byte]*graph.Node
	doors              map[byte]*graph.Node
	originalStart      *graph.Node
	starts             []*graph.Node
	startShortestPaths []djikstra.ShortestPaths
	keyShortestPaths   map[byte]djikstra.ShortestPaths
	keyDistances       *KeyDistances
	resultsCache       *DistanceCache
	keyAccessibleFrom  map[byte]int
}

func NewGame(chars [][]byte) *Game {
	g := new(Game)
	g.keys = make(map[byte]*graph.Node)
	g.doors = make(map[byte]*graph.Node)
	g.keyShortestPaths = make(map[byte]djikstra.ShortestPaths)
	g.startShortestPaths = make([]djikstra.ShortestPaths, 4, 4)
	g.keyDistances = NewKeyDistances()
	g.GameGraph = graph.NewGraph()
	g.resultsCache = NewDistanceCache()
	g.keyAccessibleFrom = make(map[byte]int)

	for y, row := range chars[1 : len(chars)-1] {
		for x, char := range row[1 : len(row)-1] {
			pos := geom.Pos{X: x + 1, Y: y + 1} // offset by 1 because we're skipping the edges

			objType := Wall

			if char >= 'A' && char <= 'Z' {
				objType = Door
			} else if char >= 'a' && char <= 'z' {
				objType = Key
			} else if char == '.' {
				objType = Empty
			} else if char == '@' {
				objType = Start
			}

			if objType != Wall {
				n := g.GameGraph.CreateNode(pos)
				n.AddProperty("value", char)
				n.AddProperty("type", objType)

				switch objType {
				case Door:
					g.doors[char] = n
				case Key:
					g.keys[char] = n
				case Start:
					g.originalStart = n
				}
			}
		}
	}

	newWalls := make([]*graph.Node, 0, 4)
	newStarts := make([]*graph.Node, 0, 4)

	// we need to modify this map for part 2
	startPos := g.originalStart.GetID().(geom.Pos)
	newWalls = append(newWalls, g.GameGraph.GetNode(geom.Pos{X: startPos.X - 1, Y: startPos.Y})) //left
	newWalls = append(newWalls, g.GameGraph.GetNode(geom.Pos{X: startPos.X + 1, Y: startPos.Y})) //right
	newWalls = append(newWalls, g.GameGraph.GetNode(geom.Pos{X: startPos.X, Y: startPos.Y + 1})) //below
	newWalls = append(newWalls, g.GameGraph.GetNode(geom.Pos{X: startPos.X, Y: startPos.Y - 1})) //above

	newStarts = append(newStarts, g.GameGraph.GetNode(geom.Pos{X: startPos.X - 1, Y: startPos.Y - 1})) // above left
	newStarts = append(newStarts, g.GameGraph.GetNode(geom.Pos{X: startPos.X + 1, Y: startPos.Y - 1})) // above right
	newStarts = append(newStarts, g.GameGraph.GetNode(geom.Pos{X: startPos.X + 1, Y: startPos.Y + 1})) // below right
	newStarts = append(newStarts, g.GameGraph.GetNode(geom.Pos{X: startPos.X - 1, Y: startPos.Y + 1})) // below left


	for _, t := range newWalls {
		t.AddProperty("type", Wall)
		t.AddProperty("value", '#')
		t.SetTraversable(false)
	}

	g.originalStart.AddProperty("type", Wall)
	g.originalStart.AddProperty("value", '#')
	g.originalStart.SetTraversable(false)

	g.starts = newStarts

	for _, t := range newStarts {
		t.AddProperty("type", Start)
		t.AddProperty("value", '@')
	}

	for y, row := range chars[1 : len(chars)-1] {
		for x := range row[1 : len(row)-1] {
			pos := geom.Pos{X: x + 1, Y: y + 1}

			n := g.GameGraph.GetNode(pos)
			cost := float64(1)

			// if n doesn't exist, it should have been a wall
			if n == nil && chars[y+1][x+1] != '#' {
				panic(errors.New("where is this?  it's not a wall."))
			}

			if n != nil && n.GetProperty("type").(ObjectType) != Wall {
				// if nodes left, right, above and below exist in the graph, they are not walls, and we need to
				// add edges

				// check for right node
				o := g.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y})
				if o != nil && o.GetProperty("type").(ObjectType) != Wall {
					n.AddEdge(o, cost)
				}
				// check for left node
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y})
				if o != nil && o.GetProperty("type").(ObjectType) != Wall {
					n.AddEdge(o, cost)
				}
				// check for node above
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1})
				if o != nil && o.GetProperty("type").(ObjectType) != Wall {
					n.AddEdge(o, cost)
				}
				// check for node below
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1})
				if o != nil && o.GetProperty("type").(ObjectType) != Wall {
					n.AddEdge(o, cost)
				}
			}

		}
	}

	for i := 0; i < len(g.starts); i++ {
		location := byte(i+48)
		g.startShortestPaths[i] = djikstra.GenerateShortestPaths(g.GameGraph, g.starts[i])
		for char, node := range g.keys {
			p, d := g.startShortestPaths[i].GetShortestPath(node)
			if d > 0 {
				g.keyDistances.AddPath(location, char, p, d)
				g.keyAccessibleFrom[char] = i
			}
		}
	}

	for char, node := range g.keys {

		g.keyShortestPaths[char] = djikstra.GenerateShortestPaths(g.GameGraph, node)
		for ochar, onode := range g.keys {
			if ochar != char {
				p, d := g.keyShortestPaths[char].GetShortestPath(onode)
				if d > 0 {
					g.keyDistances.AddPath(char, ochar, p, d)
				}
			}
		}
	}

	return g
}

func (g *Game) Execute() {
	bestPerm, minDistance := g.CalculateShortestPath(make([]byte,0,26))

	fmt.Println("shortest path through [", bestPerm, "] is: ", minDistance)
}

func (g *Game) CloseDoors() {
	for _, node := range g.doors {
		node.SetTraversable(false)
	}
}

func (g *Game) CalculateShortestPath(keysTaken []byte) (string, int) {

	minDistance := math.MaxInt64
	bestPerm := ""

	locations := []byte{48,49,50,51}
	acquiredKeys := make(map[byte]bool)
	// set our locations from the last keys taken
	// and mark that we've acquired the key in a map

	for i := len(keysTaken) -1; i >= 0; i-- {
		accessibleFrom := g.keyAccessibleFrom[keysTaken[i]]
		if locations[accessibleFrom] < byte('a') {
			locations[accessibleFrom] = keysTaken[i]
		}
		acquiredKeys[keysTaken[i]] = true
	}

	//fmt.Println("calculating shortest path from keys taken: ", string(keysTaken))

	keysAvailable := make([]byte, 0, 26)
	distances := make([]int,0,26)

	for char := range g.keys {
		if _, visited := acquiredKeys[char]; !visited {
			accessibleFrom := g.keyAccessibleFrom[char]

			d := g.keyDistances.GetDistance(locations[accessibleFrom], char, acquiredKeys)

			if d > 0 {
				keysAvailable = append(keysAvailable, char)
				distances = append(distances,d)
			}
		}
	}

	//fmt.Println("keys taken (", string(keysTaken), ") keys available: ", string(keysAvailable))

	for i := 0; i < len(keysAvailable); i++ {
		targetKeyName := keysAvailable[i]

		distanceToNextPermutationStart := distances[i]

		l := len(keysTaken)

		newKeysTaken := make([]byte, l+1, l+1)
		copy(newKeysTaken, keysTaken)
		newKeysTaken[l] = targetKeyName

		var perm string
		var dist int

		if len(newKeysTaken) == len(g.keys) {
			perm = string(newKeysTaken)
		} else {
			results := g.resultsCache.GetResults(locations,newKeysTaken)
			if results == nil {
				perm, dist = g.CalculateShortestPath(newKeysTaken)
				g.resultsCache.CacheResults(locations,newKeysTaken, perm, dist)
			} else {
				perm = results.permutation
				dist = results.distance
			}
		}

		if dist+distanceToNextPermutationStart < minDistance {
			minDistance = dist + distanceToNextPermutationStart
			bestPerm = perm
		}
	}

	return bestPerm, minDistance

}
