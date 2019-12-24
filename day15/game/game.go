package game

import (
	"github.com/mbordner/advent_of_code_2019/day15/geom"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"image"
	"log"
)

const FIELD_OF_VISION = 25

type ObjectType int

const (
	Empty ObjectType = iota
	Wall
	OxygenSystem
	Droid
	Start
)

type Object struct {
	Type         ObjectType
	Pos          geom.Pos
	HasOxygen    bool
	ShortestPath bool
}

func (o *Object) FillWithOxygen() {
	o.HasOxygen = true
}

func NewObject(objType ObjectType, x int, y int) *Object {
	o := new(Object)
	o.Type = objType
	o.Pos.X = x
	o.Pos.Y = y
	return o
}

type Game struct {
	ui.Block
	boundingBox  geom.BoundingBox
	objects      map[geom.Pos]*Object
	window       map[geom.Pos]*Object
	user         *Object
	lastDirReq   geom.Direction
	in           chan<- string
	out          <-chan string
	movecomplete chan<- string
	compquit     <-chan string
	quit         chan<- string
}

func (g *Game) SetLastDir(dir geom.Direction) {
	g.lastDirReq = dir
}

func (g *Game) loop() {

	uiEvents := ui.PollEvents()
programLoop:
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				break programLoop
			case "<Left>":
				g.lastDirReq = geom.West
				g.in <- fmt.Sprintf("%d", g.lastDirReq)

			case "<Right>":
				g.lastDirReq = geom.East
				g.in <- fmt.Sprintf("%d", g.lastDirReq)

			case "<Up>":
				g.lastDirReq = geom.North
				g.in <- fmt.Sprintf("%d", g.lastDirReq)

			case "<Down>":
				g.lastDirReq = geom.South
				g.in <- fmt.Sprintf("%d", g.lastDirReq)

			}
		case response := <-g.out:

			switch response {
			case "0": // The repair droid hit a wall. Its position has not changed.
				var wall *Object
				x, y := g.user.Pos.X, g.user.Pos.Y

				switch g.lastDirReq {
				case geom.North:
					y--
				case geom.South:
					y++
				case geom.West:
					x--
				case geom.East:
					x++
				}

				wall = NewObject(Wall, x, y)
				g.objects[wall.Pos] = wall

			case "1": // The repair droid has moved one step in the requested direction.

				if _, ok := g.objects[geom.Pos{X: g.user.Pos.X, Y: g.user.Pos.Y}]; !ok {
					empty := NewObject(Empty, g.user.Pos.X, g.user.Pos.Y)
					g.objects[empty.Pos] = empty
				}

				switch g.lastDirReq {
				case geom.North:
					g.user.Pos.Y--
				case geom.South:
					g.user.Pos.Y++
				case geom.West:
					g.user.Pos.X--
				case geom.East:
					g.user.Pos.X++
				}

			case "2": // The repair droid has moved one step in the requested direction; its new position is the location of the oxygen system.

				if _, ok := g.objects[geom.Pos{X: g.user.Pos.X, Y: g.user.Pos.Y}]; !ok {
					empty := NewObject(Empty, g.user.Pos.X, g.user.Pos.Y)
					g.objects[empty.Pos] = empty
				}

				switch g.lastDirReq {
				case geom.North:
					g.user.Pos.Y--
				case geom.South:
					g.user.Pos.Y++
				case geom.West:
					g.user.Pos.X--
				case geom.East:
					g.user.Pos.X++
				}

				oxygen := NewObject(OxygenSystem, g.user.Pos.X, g.user.Pos.Y)
				g.objects[oxygen.Pos] = oxygen
			}

			g.Refresh()
			g.movecomplete <- "next"

		case <-g.compquit:
			break programLoop
		}

	}

	g.quit <- "user exited."

}

func (g *Game) Shutdown() {
	ui.Close()
}

func (g *Game) Draw(buf *ui.Buffer) {

	droidBGColor := ui.ColorYellow

	if o, ok := g.window[geom.Pos{X: FIELD_OF_VISION, Y: FIELD_OF_VISION}]; ok {
		if o.Type == OxygenSystem {
			droidBGColor = ui.ColorRed
		} else if o.Type == Start {
			droidBGColor = ui.ColorGreen
		}
		if o.HasOxygen {
			droidBGColor = ui.ColorWhite
		}
	}

	g.window[geom.Pos{X: FIELD_OF_VISION, Y: FIELD_OF_VISION}] = g.user

	for p, o := range g.window {

		char := ui.BARS[0]
		var style ui.Style

		switch o.Type {
		case Wall:
			char = ui.BARS[8]
			style = ui.NewStyle(ui.ColorBlue, ui.ColorBlue, ui.ModifierBold)
		case Empty:
			char = ui.BARS[8]
			if o.HasOxygen {
				style = ui.NewStyle(ui.ColorWhite, ui.ColorWhite, ui.ModifierBold)
			} else if o.ShortestPath {
				style = ui.NewStyle(ui.ColorGreen, ui.ColorGreen, ui.ModifierBold)
			} else {
				style = ui.NewStyle(ui.ColorYellow, ui.ColorYellow, ui.ModifierBold)
			}

		case Start:
			char = 'S'
			if o.HasOxygen {
				style = ui.NewStyle(ui.ColorClear, ui.ColorWhite, ui.ModifierBold)
			} else {
				style = ui.NewStyle(ui.ColorClear, ui.ColorGreen, ui.ModifierBold)
			}

		case OxygenSystem:
			char = 'G'
			if o.HasOxygen {
				style = ui.NewStyle(ui.ColorClear, ui.ColorWhite, ui.ModifierBold)
			} else {
				style = ui.NewStyle(ui.ColorClear, ui.ColorRed, ui.ModifierBold)
			}
		case Droid:
			char = ui.IRREGULAR_BLOCKS[13]
			style = ui.NewStyle(ui.ColorMagenta, droidBGColor, ui.ModifierBold)
		}
		buf.SetCell(
			ui.NewCell(char, style),
			image.Pt(p.X, p.Y),
		)
	}

	g.Block.Draw(buf)
}

func (g *Game) GetObjects() []*Object {
	objects := make([]*Object, 0, len(g.objects))
	for _, o := range g.objects {
		objects = append(objects, o)
	}
	return objects
}

func (g *Game) GetObject(pos geom.Pos) *Object {
	if o, ok := g.objects[pos]; ok {
		return o
	}
	return nil
}

func (g *Game) Refresh() {
	g.window = make(map[geom.Pos]*Object)

	for p, o := range g.objects {
		dx := g.user.Pos.X - p.X
		dy := g.user.Pos.Y - p.Y
		if dx <= FIELD_OF_VISION && dy <= FIELD_OF_VISION {
			wp := geom.Pos{X: FIELD_OF_VISION - dx, Y: FIELD_OF_VISION - dy}
			g.window[wp] = o
		}

	}

	g.SetRect(0, 0, FIELD_OF_VISION+FIELD_OF_VISION+1, FIELD_OF_VISION+FIELD_OF_VISION+1)

	ui.Render(g)
}

func NewGame(in chan<- string, out <-chan string, movecomplete chan<- string, compquit <-chan string, quit chan<- string) *Game {
	g := new(Game)
	g.in = in
	g.out = out
	g.movecomplete = movecomplete
	g.compquit = compquit
	g.quit = quit
	g.objects = make(map[geom.Pos]*Object)

	g.user = NewObject(Droid, 0, 0)
	start := NewObject(Start, 0, 0)
	g.objects[start.Pos] = start

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	go g.loop()
	g.Refresh()

	return g
}
