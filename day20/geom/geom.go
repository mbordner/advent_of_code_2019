package geom

import "fmt"

type Direction int

const (
	Unknown Direction = 0
	North   Direction = 1
	South   Direction = 2
	West    Direction = 3
	East    Direction = 4
)

type BoundingBox struct {
	xMin int
	xMax int
	yMin int
	yMax int
}

func (bb BoundingBox) String() string {
	p1 := Pos{X: bb.xMin, Y: bb.yMin}
	p2 := Pos{X: bb.xMax, Y: bb.yMax}
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

type Pos struct {
	X int
	Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("{x:%d, y:%d}", p.X, p.Y)
}
