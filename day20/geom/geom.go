package geom

import (
	"fmt"
	"math"
)

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

func (bb *BoundingBox) Extend(p Pos) {
	if p.X < bb.xMin {
		bb.xMin = p.X
	}
	if p.X > bb.xMax {
		bb.xMax = p.X
	}
	if p.Y > bb.yMax {
		bb.yMax = p.Y
	}
	if p.Y < bb.yMin {
		bb.yMin = p.Y
	}
}

func (bb *BoundingBox) DistanceFromEdge(p Pos) int {
	d := math.MaxInt64

	t := bb.xMax - p.X
	if t < d {
		d = t
	}

	t = p.X - bb.xMin
	if t < d {
		d = t
	}

	t = p.Y - bb.yMin
	if t < d {
		d = t
	}

	t = bb.yMax - p.Y
	if t < d {
		d = t
	}

	return d
}

type Pos struct {
	X int
	Y int
	Z int
}

func (p Pos) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", p.X, p.Y, p.Z)
}
