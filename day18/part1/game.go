package part1

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

func (dc *DistanceCache) GetResults(keysTaken []byte) *DistanceCacheResults {
	dc.SortAcquiredKeys(keysTaken)

	if results, ok := dc.cache[string(keysTaken)]; ok {
		return &results
	}

	return nil
}

func (dc *DistanceCache) CacheResults(keysTaken []byte, permutation string, distance int) {
	dc.SortAcquiredKeys(keysTaken)
	dc.cache[string(keysTaken)] = DistanceCacheResults{permutation: permutation, distance: distance}
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
	for _, n := range path {
		if n.GetProperty("type").(ObjectType) == Door {
			requiredKey := n.GetProperty("value").(byte)
			kd.keyRequires[a][b][requiredKey+32] = true
		}
	}
}

func (kd *KeyDistances) GetDistance(a byte, b byte, acquiredKeys map[byte]bool) int {
	for rk, v := range kd.keyRequires[a][b] {
		if v {
			if acquired, ok := acquiredKeys[rk]; !ok || !acquired {
				return 0
			}
		}
	}

	return kd.keyDistance[a][b]
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
	start              *graph.Node
	startShortestPaths djikstra.ShortestPaths
	keyShortestPaths   map[byte]djikstra.ShortestPaths
	keyDistances       *KeyDistances
	resultsCache       *DistanceCache
}

func NewGame(chars [][]byte) *Game {
	g := new(Game)
	g.keys = make(map[byte]*graph.Node)
	g.doors = make(map[byte]*graph.Node)
	g.keyShortestPaths = make(map[byte]djikstra.ShortestPaths)
	g.keyDistances = NewKeyDistances()
	g.GameGraph = graph.NewGraph()
	g.resultsCache = NewDistanceCache()

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
					g.start = n
				}
			}
		}
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

			if n != nil {
				// if nodes left, right, above and below exist in the graph, they are not walls, and we need to
				// add edges

				// check for right node
				o := g.GameGraph.GetNode(geom.Pos{X: pos.X + 1, Y: pos.Y})
				if o != nil {
					n.AddEdge(o, cost)
				}
				// check for left node
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X - 1, Y: pos.Y})
				if o != nil {
					n.AddEdge(o, cost)
				}
				// check for node above
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y - 1})
				if o != nil {
					n.AddEdge(o, cost)
				}
				// check for node below
				o = g.GameGraph.GetNode(geom.Pos{X: pos.X, Y: pos.Y + 1})
				if o != nil {
					n.AddEdge(o, cost)
				}
			}

		}
	}

	g.startShortestPaths = djikstra.GenerateShortestPaths(g.GameGraph, g.start)

	for char, node := range g.keys {
		p, d := g.startShortestPaths.GetShortestPath(node)
		g.keyDistances.AddPath(byte('0'), char, p, d)

		g.keyShortestPaths[char] = djikstra.GenerateShortestPaths(g.GameGraph, node)
		for ochar, onode := range g.keys {
			if ochar != char {
				p, d := g.keyShortestPaths[char].GetShortestPath(onode)
				g.keyDistances.AddPath(char, ochar, p, d)
			}
		}
	}

	/*
		g.startShortestPaths = djikstra.GenerateShortestPaths(g.GameGraph, g.start)

		ph := newPathHeap(len(g.keys))

		for char, node := range g.keys {
			nodes, distance := g.startShortestPaths.GetShortestPath(node)

			if distance > 0 {
				path := Path{
					Value:    char,
					Distance: int(distance),
					Nodes:    nodes,
				}

				ph.push(&path)
			}


		}

		for ph.paths.Len() > 0 {
			path := ph.pop()
			fmt.Println("----")
			fmt.Println("distance to key ", string(path.Value), " is ", path.Distance)
			for _, n := range path.Nodes {
				value := n.GetProperty("value").(byte)
				pos := n.GetID().(geom.Pos)

				if value >= 'A' && value <= 'Z' {
					fmt.Println("blocked by ", string(value), " at pos: ", pos)
				}

				//fmt.Printf("-> [%s]%s ", string(value), pos)
			}
			fmt.Println("\n")
		}

		g.startShortestPaths = djikstra.GenerateShortestPaths(g.GameGraph, g.start)

		keysAvailableAtStart := make([]byte, 0, 26)

		for char, node := range g.keys {
			_, distance := g.startShortestPaths.GetShortestPath(node)

			if distance > 0 {
				keysAvailableAtStart = append(keysAvailableAtStart, char)
			}
		}

		fmt.Println("keys available at start: ", string(keysAvailableAtStart))

	*/

	return g
}

func (g *Game) Execute() {
	bestPerm, minDistance := g.CalculateShortestPath(make([]byte, 0, 26))

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

	acquiredKeys := make(map[byte]bool)
	for _, b := range keysTaken {
		acquiredKeys[b] = true
	}

	//fmt.Println("calculating shortest path from keys taken: ", string(keysTaken))

	keysAvailable := make([]byte, 0, 26)

	location := byte('0')

	if len(keysTaken) > 0 {
		location = keysTaken[len(keysTaken)-1]
	}

	for char := range g.keys {
		if _, visited := acquiredKeys[char]; !visited {

			d := g.keyDistances.GetDistance(location, char, acquiredKeys)

			if d > 0 {
				keysAvailable = append(keysAvailable, char)
			}
		}
	}

	//fmt.Println("keys taken (", string(keysTaken), ") keys available: ", string(keysAvailable))

	for i := 0; i < len(keysAvailable); i++ {
		targetKeyName := keysAvailable[i]

		distanceToNextPermutationStart := g.keyDistances.GetDistance(location, keysAvailable[i], acquiredKeys)

		l := len(keysTaken)

		newKeysTaken := make([]byte, l+1, l+1)
		copy(newKeysTaken, keysTaken)
		newKeysTaken[l] = targetKeyName

		var perm string
		var dist int

		if len(newKeysTaken) == len(g.keys) {
			perm = string(newKeysTaken)
		} else {
			results := g.resultsCache.GetResults(newKeysTaken)
			if results == nil {
				perm, dist = g.CalculateShortestPath(newKeysTaken)
				g.resultsCache.CacheResults(newKeysTaken,perm,dist)
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
