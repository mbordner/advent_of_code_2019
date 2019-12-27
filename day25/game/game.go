package game

import (
	"bytes"
	"fmt"
	"github.com/mbordner/advent_of_code_2019/day25/graph"
	"github.com/mbordner/advent_of_code_2019/day25/graph/djikstra"
	"regexp"
	"sort"
	"strings"
)

type Empty struct{}

var empty Empty

var (
	reRoomId   = regexp.MustCompile(`==\s(.*)\s==`)
	reListItem = regexp.MustCompile(`-\s(.*)`)
)

type Room struct {
	ID    string
	Doors []string
	Items []string
}

func NewRoom(id string) *Room {
	r := new(Room)
	r.ID = id
	r.Doors = make([]string, 0, 4)
	r.Items = make([]string, 0, 10)
	return r
}

type Game struct {
	GameGraph          *graph.Graph
	DontWant           map[string]Empty
	bytes              []byte
	buf                *bytes.Buffer
	commands           []string
	ptr                int
	start              *graph.Node
	current            *graph.Node
	lastTraveled       *graph.Edge
	inventory          []string
	securityCheckpoint string
	invAttempt         int
	possibleItems      []string
}

func (g *Game) sortDoorsToExploreAll(doors []string) {
	order := make(map[string]int)
	dirCame := "north"
	if g.lastTraveled != nil {
		dirCame = g.lastTraveled.GetProperty("direction").(string)
	}
	switch dirCame {
	case "west":
		order["south"] = 0
		order["west"] = 1
		order["north"] = 2
		order["east"] = 3
	case "north":
		order["west"] = 0
		order["north"] = 1
		order["east"] = 2
		order["south"] = 3
	case "east":
		order["north"] = 0
		order["east"] = 1
		order["south"] = 2
		order["west"] = 3
	case "south":
		order["east"] = 0
		order["south"] = 1
		order["west"] = 2
		order["north"] = 3
	}

	sort.Slice(doors, func(i, j int) bool {
		return order[doors[i]] < order[doors[j]]
	})
}

func (g *Game) OutputByte(b byte) {
	g.buf.WriteByte(b)
	roomInfo := g.analyzeBuffer()
	if roomInfo != nil {

		// add any take commands for items available, excluding those we explicitly do not want
		for _, item := range roomInfo.Items {
			if _, ok := g.DontWant[item]; !ok {
				g.commands = append(g.commands, fmt.Sprintf("take %s", item))
				g.inventory = append(g.inventory, item)
			}
		}

		// get node, if it's nil, we have not been here, so we'll add it to the graph
		n := g.GameGraph.GetNode(roomInfo.ID)
		if n == nil {
			n = g.GameGraph.CreateNode(roomInfo.ID)
			for _, oRoom := range roomInfo.Doors {
				e := n.AddEdge(nil, float64(1)) // adding nil destination, because we have not traveled this route
				e.AddProperty("direction", oRoom)
			}
		}

		g.current = n
		if g.lastTraveled != nil {
			g.lastTraveled.SetDestination(n) // link the previous node to this one
		}

		if g.start == nil {
			g.start = g.current
		}

		allVisited := true

		edges := n.GetEdges()
		for _, e := range edges {
			if e.GetDestination() == nil {
				g.lastTraveled = e
				g.commands = append(g.commands, e.GetProperty("direction").(string))
				allVisited = false
				break
			}
		}

		if allVisited || g.invAttempt > 0 {
			if n.GetID().(string) == g.securityCheckpoint {
				// we've picked up all the things, and went to the checkpoint
				// at this point, we have all of the items that don't ruin us,
				// and one of the edges is pointing to the checkpoint..
				// this edge is the direction we want to go
				edges := n.GetEdges()
				var dir string
				for _, e := range edges {
					if e.GetDestination() == n {
						dir = e.GetProperty("direction").(string)
						g.lastTraveled = e
						break
					}
				}

				if g.invAttempt == 0 {
					g.possibleItems = make([]string,len(g.inventory),len(g.inventory))
					copy(g.possibleItems,g.inventory)
				}
				// on the first try, or every time we enter the room after the first try, we pick up everything
				// so we drop everything
				for _, i := range g.possibleItems {
					g.commands = append(g.commands, fmt.Sprintf("drop %s", i))
				}

				g.inventory = make([]string, 0, len(g.possibleItems))

				g.invAttempt += 1 // increment the attempt
				for i := range g.possibleItems {
					n := 1 << i
					if n&g.invAttempt > 0 {
						g.commands = append(g.commands, fmt.Sprintf("take %s", g.possibleItems[i]))
						g.inventory = append(g.inventory, g.possibleItems[i])
					}
				}

				// attempt
				g.commands = append(g.commands, dir)

			} else if n.GetID().(string) == g.start.GetID().(string) {
				// all nodes have been visited, we should be at the start since
				// we were going around the world turning left all the time
				// now, just issue the commands to get to the security checkpoint
				sps := djikstra.GenerateShortestPaths(g.GameGraph, n)
				cp := g.GameGraph.GetNode(g.securityCheckpoint)
				sp, distance := sps.GetShortestPath(cp)

				cmds := make([]string, 0, int(distance))
				for len(sp) > 0 {
					edges := n.GetEdges()
					for _, e := range edges {
						if e.GetDestination() == sp[0] {
							cmds = append(cmds, e.GetProperty("direction").(string))
							n = sp[0]
							sp = sp[1:]
							break
						}
					}
				}

				g.commands = append(g.commands, cmds...)
			} else if g.lastTraveled.GetSource().GetID().(string) == g.securityCheckpoint {
				// we should have passed checkpoint
				fmt.Println(n.GetID().(string))
			}
		}
	}

}

func (g *Game) analyzeBuffer() *Room {
	var r *Room
	lines := strings.Split(g.buf.String(), "\n")
	i := len(lines) - 1
	if lines[i] == "Command?" {

		// work backwards through the Command? to find the room we're in, because if we get thrown out
		// past the security checkpoint, we get automatically returned without a command
		for i--; i >= 0; i-- {
			if matches := reRoomId.FindStringSubmatch(lines[i]); len(matches) > 0 {
				r = NewRoom(matches[1])
				var door, item bool
				for _, line := range lines[i:] {
					if line == "Doors here lead:" {
						door = true
					} else if line == "Items here:" {
						item = true
					} else if line == "" {
						door = false
						item = false
					} else if matches = reListItem.FindStringSubmatch(line); len(matches) > 0 {
						if door {
							r.Doors = append(r.Doors, matches[1])
						} else if item {
							r.Items = append(r.Items, matches[1])
						}
					}
				}
				g.sortDoorsToExploreAll(r.Doors)
				break
			}
		}

		// we got to the Command? prompt, so reset this buffer, and analyze the data about the room
		g.buf.Reset()
	}
	return r
}

func (g *Game) GetAllCommands() []string {
	return g.commands
}

func (g *Game) GetCurrentCommands() []string {
	cmds := g.commands[g.ptr:]
	g.ptr = len(g.commands)
	return cmds
}

func NewGame(dw []string, checkpoint string) *Game {
	g := new(Game)
	g.GameGraph = graph.NewGraph()
	g.DontWant = make(map[string]Empty)
	for i := range dw {
		g.DontWant[dw[i]] = empty
	}
	g.bytes = make([]byte, 0, 8092)
	g.buf = bytes.NewBuffer(g.bytes)
	g.commands = make([]string, 0, 80)
	g.inventory = make([]string, 0, 20)
	g.securityCheckpoint = checkpoint
	return g
}
