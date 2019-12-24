package geom

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

type Pos struct {
	X int
	Y int
}
