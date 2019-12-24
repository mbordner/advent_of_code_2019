package game

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"image"
	"log"
)

type ObjectType int

const (
	Empty ObjectType = iota
	Wall
	Block
	HorizontalPaddle
	Ball
)

type JoyStickPosition int

const (
	Neutral JoyStickPosition = 0
	Left    JoyStickPosition = -1
	Right   JoyStickPosition = 1
)

type BoundingBox struct {
	xMin int
	xMax int
	yMin int
	yMax int
}

type Pos struct {
	X int
	Y int
}

type Object struct {
	Type ObjectType
	Pos  Pos
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
	boundingBox BoundingBox
	objects     map[Pos]*Object
	score       int
	joystick    JoyStickPosition
	in          chan<- string
	compquit    <-chan string
	quit        chan<- string
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
				if g.joystick > Left {
					g.joystick--
				}
				g.in <- fmt.Sprintf("%d", g.joystick)
				g.joystick = Neutral
			case "<Right>":
				if g.joystick < Right {
					g.joystick++
				}
				g.in <- fmt.Sprintf("%d", g.joystick)
				g.joystick = Neutral
			case "<Space>":
				g.joystick = Neutral
				g.in <- fmt.Sprintf("%d", g.joystick)
			}
		case <-g.compquit:
			break programLoop
		}

	}

	g.Shutdown()

	g.quit <- "user exited."

}

func (g *Game) Shutdown() {
	ui.Close()
}

func (g *Game) Draw(buf *ui.Buffer) {
	g.Block.Draw(buf)

	for p, o := range g.objects {

		char := ui.BARS[0]
		color := ui.ColorClear
		switch o.Type {
		case Wall:
			char = ui.SHADED_BLOCKS[2]
			color = ui.ColorWhite
		case Block:
			char = ui.BARS[5]
			color = ui.ColorYellow
		case HorizontalPaddle:
			char = ui.BARS[1]
			color = ui.ColorYellow
		case Ball:
			char = ui.BARS[2]
			color = ui.ColorWhite
		}
		buf.SetCell(
			ui.NewCell(char, ui.NewStyle(color)),
			image.Pt(p.X, p.Y),
		)
	}

	buf.SetString(fmt.Sprintf("Score: %d",g.score),ui.NewStyle(ui.ColorGreen),image.Pt(5,g.boundingBox.yMax+2))
}

func (g *Game) GetScore() int {
	return g.score
}

func (g *Game) SetScore(s int) {
	g.score = s
}

func (g *Game) GetJoystick() JoyStickPosition {
	return g.joystick
}

func (g *Game) SetJoystick(p JoyStickPosition) {
	g.joystick = p
}

func (g *Game) GetObjects() []*Object {
	objects := make([]*Object, 0, len(g.objects))
	for _, o := range g.objects {
		objects = append(objects, o)
	}
	return objects
}

func (g *Game) SetTile(objType ObjectType, x int, y int) {
	o := NewObject(objType, x, y)

	if x < g.boundingBox.xMin {
		g.boundingBox.xMin = x
	}
	if x > g.boundingBox.xMax {
		g.boundingBox.xMax = x
	}
	if y < g.boundingBox.yMin {
		g.boundingBox.yMin = y
	}
	if y > g.boundingBox.yMax {
		g.boundingBox.yMax = y
	}

	g.objects[o.Pos] = o

	g.SetRect(g.boundingBox.xMin, g.boundingBox.yMin, g.boundingBox.xMax+1, g.boundingBox.yMax+5)

	ui.Render(g)
}

func NewGame(in chan<- string, compquit <-chan string, quit chan<- string) *Game {
	g := new(Game)
	g.in = in
	g.compquit = compquit
	g.quit = quit
	g.objects = make(map[Pos]*Object)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	go g.loop()

	return g
}
