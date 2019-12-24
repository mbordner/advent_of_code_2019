package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	asteroidChar = '#'
)

type Pos struct {
	X int
	Y int
}

type Vector struct {
	Distance float64
	Angle    float64
}

type Asteroid struct {
	pos       Pos
	asteroids map[float64][]*Asteroid
	distances map[*Asteroid]float64
}

// https://i.stack.imgur.com/b6C3y.png
func (a *Asteroid) Project(b *Asteroid) *Vector {
	dy := float64(a.pos.Y) - float64(b.pos.Y)
	// swapping positions here because our coordinate space has origin in top left corner of grid
	// and we're trying to get the angles as if asteroid A is the origin here
	dx := float64(b.pos.X) - float64(a.pos.X)

	two := float64(2)

	v := new(Vector)
	v.Distance = math.Sqrt(math.Pow(dy, two) + math.Pow(dx, two))
	v.Angle = math.Atan2(dy, dx)

	return v
}

func (a *Asteroid) GetRotations() [][]*Asteroid {
	angles := make([]float64, 0, len(a.asteroids))
	for angle := range a.asteroids {
		angles = append(angles, angle)
	}
	halfPi := math.Pi / 2

	// angle up is π/2 right is 0, down is -π/2 and left is π,
	// we want to rotate clockwise pointing up, so we need to sort the asteroids by angles
	// π/2 -> -π, and add on π -> π/2 at the end

	// sort from π -> -π
	sort.Sort(sort.Reverse(sort.Float64Slice(angles)))
	for i, v := range angles {
		if v <= halfPi {
			// adjust the angles so starting from top and going clockwise
			angles = append(angles[i:], angles[:i]...)
			break
		}
	}
	asteroids := make([][]*Asteroid, len(angles), len(angles))
	for i := 0; i < len(angles); i++ {
		asteroids[i] = a.asteroids[angles[i]]
	}
	return asteroids
}

func (a *Asteroid) Add(b *Asteroid) {
	v := a.Project(b)
	a.distances[b] = v.Distance
	if _, ok := a.asteroids[v.Angle]; !ok {
		a.asteroids[v.Angle] = make([]*Asteroid, 0, 10)
	}
	index := 0
	for index < len(a.asteroids[v.Angle]) {
		t := a.asteroids[v.Angle][index]
		if v.Distance < a.distances[t] {
			break
		}
		index++
	}
	if index == len(a.asteroids[v.Angle]) {
		a.asteroids[v.Angle] = append(a.asteroids[v.Angle], b)
	} else {
		a.asteroids[v.Angle] = append(a.asteroids[v.Angle][:index], append([]*Asteroid{b}, a.asteroids[v.Angle][index:]...)...)
	}
}

func (a *Asteroid) GetVisible() []*Asteroid {
	visible := make([]*Asteroid, 0, len(a.asteroids))
	for _, asteroids := range a.asteroids {
		visible = append(visible, asteroids[0])
	}
	return visible
}

func (a *Asteroid) GetVisibleCount() int {
	return len(a.asteroids)
}

func (a *Asteroid) GetClosest() *Asteroid {
	var closest *Asteroid
	var distance *float64
	visible := a.GetVisible()
	for _, t := range visible {
		d := a.distances[t]
		if distance == nil {
			distance = &d
			closest = t
		} else if d < *distance {
			distance = &d
			closest = t
		}
	}
	return closest
}

func (a Asteroid) String() string {
	return fmt.Sprintf("%d,%d", a.pos.X, a.pos.Y)
}

func NewAsteroid(x, y int) *Asteroid {
	a := new(Asteroid)
	a.pos = Pos{X: x, Y: y}
	a.asteroids = make(map[float64][]*Asteroid)
	a.distances = make(map[*Asteroid]float64)
	return a
}

func main() {
	asteroidMap := getMap()

	asteroids := make([]*Asteroid, 0, strings.Count(asteroidMap, "#"))

	lines := strings.Split(asteroidMap, "\n")

	for y, line := range lines {
		for x, char := range line {
			if char == asteroidChar {
				asteroids = append(asteroids, NewAsteroid(x, y))
			}
		}
	}

	fmt.Println("number of asteroids: ", len(asteroids))

	var bestLocation *Asteroid

	for i := 0; i < len(asteroids); i++ {
		for j := 0; j < len(asteroids); j++ {
			if j != i {
				asteroids[i].Add(asteroids[j])
			}
		}
		if bestLocation == nil {
			bestLocation = asteroids[i]
		} else if asteroids[i].GetVisibleCount() > bestLocation.GetVisibleCount() {
			bestLocation = asteroids[i]
		}
	}

	fmt.Println("best location is at: ", bestLocation)
	fmt.Println("it can see ", bestLocation.GetVisibleCount(), " other asteroids")

	rotations := bestLocation.GetRotations()

	destroyed := 0
	var lastDestroyed *Asteroid
	for i := 0; i < len(rotations); i++ {
		if len(rotations[i]) > 0 {
			lastDestroyed = rotations[i][0]
			rotations[i] = rotations[i][1:]
			destroyed++
			if destroyed == 200 {
				break
			}
		}

		if i == len(rotations)-1 {
			i = -1
		}
	}

	fmt.Println("200th destroyed is: ", lastDestroyed)

}

func getMap() string {
	return `.............#..#.#......##........#..#
.#...##....#........##.#......#......#.
..#.#.#...#...#...##.#...#.............
.....##.................#.....##..#.#.#
......##...#.##......#..#.......#......
......#.....#....#.#..#..##....#.......
...................##.#..#.....#.....#.
#.....#.##.....#...##....#####....#.#..
..#.#..........#..##.......#.#...#....#
...#.#..#...#......#..........###.#....
##..##...#.#.......##....#.#..#...##...
..........#.#....#.#.#......#.....#....
....#.........#..#..##..#.##........#..
........#......###..............#.#....
...##.#...#.#.#......#........#........
......##.#.....#.#.....#..#.....#.#....
..#....#.###..#...##.#..##............#
...##..#...#.##.#.#....#.#.....#...#..#
......#............#.##..#..#....##....
.#.#.......#..#...###...........#.#.##.
........##........#.#...#.#......##....
.#.#........#......#..........#....#...
...............#...#........##..#.#....
.#......#....#.......#..#......#.......
.....#...#.#...#...#..###......#.##....
.#...#..##................##.#.........
..###...#.......#.##.#....#....#....#.#
...#..#.......###.............##.#.....
#..##....###.......##........#..#...#.#
.#......#...#...#.##......#..#.........
#...#.....#......#..##.............#...
...###.........###.###.#.....###.#.#...
#......#......#.#..#....#..#.....##.#..
.##....#.....#...#.##..#.#..##.......#.
..#........#.......##.##....#......#...
##............#....#.#.....#...........
........###.............##...#........#
#.........#.....#..##.#.#.#..#....#....
..............##.#.#.#...........#.....`
}

func getTestMap() string {
	return `.#..#
.....
#####
....#
...##`
}

func getTestMap2() string {
	return `......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`
}

func getTestMap3() string {
	return `#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`
}

func getTestMap4() string {
	return `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`
}

func getTestMap5() string {
	return `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`
}

/**
--- Day 10: Monitoring Station ---
You fly into the asteroid belt and reach the Ceres monitoring station. The Elves here have an emergency: they're having trouble tracking all of the asteroids and can't be sure they're safe.

The Elves would like to build a new monitoring station in a nearby area of space; they hand you a map of all of the asteroids in that region (your puzzle input).

The map indicates whether each position is empty (.) or contains an asteroid (#). The asteroids are much smaller than they appear on the map, and every asteroid is exactly in the center of its marked position. The asteroids can be described with X,Y coordinates where X is the distance from the left edge and Y is the distance from the top edge (so the top-left corner is 0,0 and the position immediately to its right is 1,0).

Your job is to figure out which asteroid would be the best place to build a new monitoring station. A monitoring station can detect any asteroid to which it has direct line of sight - that is, there cannot be another asteroid exactly between them. This line of sight can be at any angle, not just lines aligned to the grid or diagonally. The best location is the asteroid that can detect the largest number of other asteroids.

For example, consider the following map:

.#..#
.....
#####
....#
...##
The best location for a new monitoring station on this map is the highlighted asteroid at 3,4 because it can detect 8 asteroids, more than any other location. (The only asteroid it cannot detect is the one at 1,0; its view of this asteroid is blocked by the asteroid at 2,2.) All other asteroids are worse locations; they can detect 7 or fewer other asteroids. Here is the number of other asteroids a monitoring station on each asteroid could detect:

.7..7
.....
67775
....7
...87
Here is an asteroid (#) and some examples of the ways its line of sight might be blocked. If there were another asteroid at the location of a capital letter, the locations marked with the corresponding lowercase letter would be blocked and could not be detected:

#.........
...A......
...B..a...
.EDCG....a
..F.c.b...
.....c....
..efd.c.gb
.......c..
....f...c.
...e..d..c
Here are some larger examples:

Best is 5,8 with 33 other asteroids detected:

......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####
Best is 1,2 with 35 other asteroids detected:

#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
Best is 6,3 with 41 other asteroids detected:

.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..
Best is 11,13 with 210 other asteroids detected:

.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
Find the best location for a new monitoring station. How many other asteroids can be detected from that location?

Your puzzle answer was 299.

--- Part Two ---
Once you give them the coordinates, the Elves quickly deploy an Instant Monitoring Station to the location and discover the worst: there are simply too many asteroids.

The only solution is complete vaporization by giant laser.

Fortunately, in addition to an asteroid scanner, the new monitoring station also comes equipped with a giant rotating laser perfect for vaporizing asteroids. The laser starts by pointing up and always rotates clockwise, vaporizing any asteroid it hits.

If multiple asteroids are exactly in line with the station, the laser only has enough power to vaporize one of them before continuing its rotation. In other words, the same asteroids that can be detected can be vaporized, but if vaporizing one asteroid makes another one detectable, the newly-detected asteroid won't be vaporized until the laser has returned to the same position by rotating a full 360 degrees.

For example, consider the following map, where the asteroid with the new monitoring station (and laser) is marked X:

.#....#####...#..
##...##.#####..##
##...#...#.#####.
..#.....X...###..
..#.#.....#....##
The first nine asteroids to get vaporized, in order, would be:

.#....###24...#..
##...##.13#67..9#
##...#...5.8####.
..#.....X...###..
..#.#.....#....##
Note that some asteroids (the ones behind the asteroids marked 1, 5, and 7) won't have a chance to be vaporized until the next full rotation. The laser continues rotating; the next nine to be vaporized are:

.#....###.....#..
##...##...#.....#
##...#......1234.
..#.....X...5##..
..#.9.....8....76
The next nine to be vaporized are then:

.8....###.....#..
56...9#...#.....#
34...7...........
..2.....X....##..
..1..............
Finally, the laser completes its first full rotation (1 through 3), a second rotation (4 through 8), and vaporizes the last asteroid (9) partway through its third rotation:

......234.....6..
......1...5.....7
.................
........X....89..
.................
In the large example above (the one with the best monitoring station location at 11,13):

The 1st asteroid to be vaporized is at 11,12.
The 2nd asteroid to be vaporized is at 12,1.
The 3rd asteroid to be vaporized is at 12,2.
The 10th asteroid to be vaporized is at 12,8.
The 20th asteroid to be vaporized is at 16,0.
The 50th asteroid to be vaporized is at 16,9.
The 100th asteroid to be vaporized is at 10,16.
The 199th asteroid to be vaporized is at 9,6.
The 200th asteroid to be vaporized is at 8,2.
The 201st asteroid to be vaporized is at 10,9.
The 299th and final asteroid to be vaporized is at 11,1.
The Elves are placing bets on which will be the 200th asteroid to be vaporized. Win the bet by determining which asteroid that will be; what do you get if you multiply its X coordinate by 100 and then add its Y coordinate? (For example, 8,2 becomes 802.)

Your puzzle answer was 1419.

Both parts of this puzzle are complete! They provide two gold stars: **
 */