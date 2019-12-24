package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	vRe = regexp.MustCompile(`[xyz]=(-?\d+)`)
)

type Axis int

const (
	X Axis = iota
	Y
	Z
)

type Vector struct {
	X int
	Y int
	Z int
}

func (v Vector) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", v.X, v.Y, v.Z)
}

func (v *Vector) Copy() *Vector {
	u := new(Vector)
	u.X = v.X
	u.Y = v.Y
	u.Z = v.Z
	return u
}

func (v *Vector) Add(u *Vector) {
	v.X += u.X
	v.Y += u.Y
	v.Z += u.Z
}

type Moon struct {
	pos  *Vector
	vel  *Vector
	step *Vector
}
type Moons []*Moon

func (m Moon) String() string {
	return fmt.Sprintf("[p:%s, v:%s, e:%d]", m.pos, m.vel, m.TotalEnergy())
}

func (m *Moon) CalculateVelocityStepChange(n *Moon) {
	if n.pos.X > m.pos.X {
		m.step.X++
	} else if n.pos.X < m.pos.X {
		m.step.X--
	}
	if n.pos.Y > m.pos.Y {
		m.step.Y++
	} else if n.pos.Y < m.pos.Y {
		m.step.Y--
	}
	if n.pos.Z > m.pos.Z {
		m.step.Z++
	} else if n.pos.Z < m.pos.Z {
		m.step.Z--
	}
}

func (m *Moon) ApplyVelocityStepChange() {
	m.vel.Add(m.step)
	m.step = new(Vector)
	m.pos.Add(m.vel)
}

func (m *Moon) PotentialEnergy() int {
	sum := 0
	sum += int(math.Abs(float64(m.pos.X)))
	sum += int(math.Abs(float64(m.pos.Y)))
	sum += int(math.Abs(float64(m.pos.Z)))
	return sum
}

func (m *Moon) KineticEnergy() int {
	sum := 0
	sum += int(math.Abs(float64(m.vel.X)))
	sum += int(math.Abs(float64(m.vel.Y)))
	sum += int(math.Abs(float64(m.vel.Z)))
	return sum
}

func (m *Moon) TotalEnergy() int {
	return m.PotentialEnergy() * m.KineticEnergy()
}

func NewMoon(config string) *Moon {
	m := new(Moon)
	m.pos = new(Vector)
	m.vel = new(Vector)
	m.step = new(Vector)

	matches := vRe.FindAllStringSubmatch(config, -1)
	if len(matches) != 3 {
		panic(fmt.Errorf("invalid vector configuration: %s", config))
	}
	for _, components := range matches {
		if len(components) != 2 {
			panic(fmt.Errorf("invalid vector component"))
		}
		axis := components[0][0]
		val, err := strconv.Atoi(components[1])
		if err != nil {
			panic(errors.New("invalid vector value"))
		}
		switch axis {
		case 'x':
			m.pos.X = val
		case 'y':
			m.pos.Y = val
		case 'z':
			m.pos.Z = val
		default:
			panic(errors.New("invalid vector component"))
		}
	}

	return m
}

func (ms Moons) TotalEnergy() int {
	te := 0
	for _, m := range ms {
		te += m.TotalEnergy()
	}
	return te
}

func (ms Moons) TotalKineticEnergy() int {
	e := 0
	for _, m := range ms {
		e += m.KineticEnergy()
	}
	return e
}

func (ms Moons) TotalPotentialEnergy() int {
	e := 0
	for _, m := range ms {
		e += m.PotentialEnergy()
	}
	return e
}

func (ms Moons) Step() {
	for i, m := range ms {
		for j, n := range ms {
			if j != i {
				m.CalculateVelocityStepChange(n)
			}
		}
	}

	for i := range ms {
		ms[i].ApplyVelocityStepChange()
	}
}

func (ms Moons) GetAxisValues(axis Axis) []int {
	vals := make([]int, len(ms), len(ms))
	for i := range ms {
		var val int
		if axis == X {
			val = ms[i].pos.X
		} else if axis == Y {
			val = ms[i].pos.Y
		} else if axis == Z {
			val = ms[i].pos.Z
		}
		vals[i] = val
	}
	return vals
}

func equal(orig, cur []int) bool {
	for i := 1; i < len(orig); i++ {
		if orig[i] != cur[i] {
			return false
		}
	}
	return true
}

func lcm1(a []uint64) uint64 {
	values := make(map[uint64]int)
	i := uint64(1)
	for {
		for j := 0; j < len(a); j++ {
			product := a[j] * i
			if _, ok := values[product]; !ok {
				values[product] = 0
			}
			values[product]++
			if values[product] == len(a) {
				return product
			}
		}
		i++
	}

}

func primeFactors(n uint64) []uint64 {
	factors := make([]uint64, 0, 25)
	// number of 2s
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}

	// n must be odd at this point. So we can skip
	// one element (Note i = i +2)
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		for n%uint64(i) == uint64(0) {
			factors = append(factors, uint64(i))
			n = n / uint64(i)
		}
	}

	// This condition is to handle the case when n
	// is a prime number greater than 2
	if n > 2 {
		factors = append(factors, n)
	}

	return factors
}

func lcm(a []uint64) uint64 {
	factors := make(map[uint64]int)
	for _, n := range a {
		f := primeFactors(n)
		counts := make(map[uint64]int)
		for _, t := range f {
			if _, ok := counts[t]; !ok {
				counts[t] = 0
			}
			counts[t]++
		}
		for k, v := range counts {
			if _, ok := factors[k]; !ok {
				factors[k] = 0
			}
			if v > factors[k] {
				factors[k] = v
			}
		}
	}
	product := uint64(1)
	for k := range factors {
		for factors[k] > 0 {
			product *= k
			factors[k]--
		}
	}
	return product
}

func main() {
	configuration := getMoonConfigurations()
	moons := make(Moons, len(configuration), len(configuration))
	for i, config := range configuration {
		moons[i] = NewMoon(config)
	}

	steps := uint64(1000)

	for steps > 0 {

		moons.Step()

		steps--
	}

	fmt.Println(moons)

	fmt.Println("total energy:", moons.TotalEnergy())

	fmt.Println("calculating steps it takes to return to original positions...")

	// i noticed that when total energy is 0, it is halfway through the cycle, but even finding this
	// was too slow

	/*

		steps = uint64(0)

			moons = make(Moons, len(configuration), len(configuration))
			for i, config := range configuration {
				moons[i] = NewMoon(config)
			}

		for {

			moons.Step()

			steps++

			if moons.TotalKineticEnergy() == 0 {
				break
			}

		}

		fmt.Println("it took ", steps, " to get to total kinetic energy of 0")
		fmt.Println("so it should take ", steps * 2, " steps to return to original position")

	*/

	// looked for hints, and people attack this problem with the axes independently, and look for cycles that way
	// then find least common multiplier

	steps = uint64(0)

	moons = make(Moons, len(configuration), len(configuration))
	for i, config := range configuration {
		moons[i] = NewMoon(config)
	}

	origX := moons.GetAxisValues(X)
	origY := moons.GetAxisValues(Y)
	origZ := moons.GetAxisValues(Z)

	x, y, z := uint64(0), uint64(0), uint64(0)

	steps++ // need to advance here, because the last step to get the total energy back to 0 is one after we find the first time the positions match

	for x == 0 || y == 0 || z == 0 {

		moons.Step()
		steps++

		if x == 0 {
			if equal(moons.GetAxisValues(X), origX) {
				x = steps
			}
		}

		if y == 0 {
			if equal(moons.GetAxisValues(Y), origY) {
				y = steps
			}
		}

		if z == 0 {
			if equal(moons.GetAxisValues(Z), origZ) {
				z = steps
			}
		}

	}

	stepsToConverge := lcm([]uint64{x, y, z})
	fmt.Println(stepsToConverge)

}

func getMoonConfigurationsTest1() []string {
	config := `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`

	return strings.Split(config, "\n")
}

func getMoonConfigurationsTest2() []string {
	config := `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`

	return strings.Split(config, "\n")
}

func getMoonConfigurations() []string {
	config := `<x=3, y=15, z=8>
<x=5, y=-1, z=-2>
<x=-10, y=8, z=2>
<x=8, y=4, z=-5>`

	return strings.Split(config, "\n")
}

/**
--- Day 12: The N-Body Problem ---
The space near Jupiter is not a very safe place; you need to be careful of a big distracting red spot, extreme radiation, and a whole lot of moons swirling around. You decide to start by tracking the four largest moons: Io, Europa, Ganymede, and Callisto.

After a brief scan, you calculate the position of each moon (your puzzle input). You just need to simulate their motion so you can avoid them.

Each moon has a 3-dimensional position (x, y, and z) and a 3-dimensional velocity. The position of each moon is given in your scan; the x, y, and z velocity of each moon starts at 0.

Simulate the motion of the moons in time steps. Within each time step, first update the velocity of every moon by applying gravity. Then, once all moons' velocities have been updated, update the position of every moon by applying velocity. Time progresses by one step once all of the positions are updated.

To apply gravity, consider every pair of moons. On each axis (x, y, and z), the velocity of each moon changes by exactly +1 or -1 to pull the moons together. For example, if Ganymede has an x position of 3, and Callisto has a x position of 5, then Ganymede's x velocity changes by +1 (because 5 > 3) and Callisto's x velocity changes by -1 (because 3 < 5). However, if the positions on a given axis are the same, the velocity on that axis does not change for that pair of moons.

Once all gravity has been applied, apply velocity: simply add the velocity of each moon to its own position. For example, if Europa has a position of x=1, y=2, z=3 and a velocity of x=-2, y=0,z=3, then its new position would be x=-1, y=2, z=6. This process does not modify the velocity of any moon.

For example, suppose your scan reveals the following positions:

<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>
Simulating the motion of these moons would produce the following:

After 0 steps:
pos=<x=-1, y=  0, z= 2>, vel=<x= 0, y= 0, z= 0>
pos=<x= 2, y=-10, z=-7>, vel=<x= 0, y= 0, z= 0>
pos=<x= 4, y= -8, z= 8>, vel=<x= 0, y= 0, z= 0>
pos=<x= 3, y=  5, z=-1>, vel=<x= 0, y= 0, z= 0>

After 1 step:
pos=<x= 2, y=-1, z= 1>, vel=<x= 3, y=-1, z=-1>
pos=<x= 3, y=-7, z=-4>, vel=<x= 1, y= 3, z= 3>
pos=<x= 1, y=-7, z= 5>, vel=<x=-3, y= 1, z=-3>
pos=<x= 2, y= 2, z= 0>, vel=<x=-1, y=-3, z= 1>

After 2 steps:
pos=<x= 5, y=-3, z=-1>, vel=<x= 3, y=-2, z=-2>
pos=<x= 1, y=-2, z= 2>, vel=<x=-2, y= 5, z= 6>
pos=<x= 1, y=-4, z=-1>, vel=<x= 0, y= 3, z=-6>
pos=<x= 1, y=-4, z= 2>, vel=<x=-1, y=-6, z= 2>

After 3 steps:
pos=<x= 5, y=-6, z=-1>, vel=<x= 0, y=-3, z= 0>
pos=<x= 0, y= 0, z= 6>, vel=<x=-1, y= 2, z= 4>
pos=<x= 2, y= 1, z=-5>, vel=<x= 1, y= 5, z=-4>
pos=<x= 1, y=-8, z= 2>, vel=<x= 0, y=-4, z= 0>

After 4 steps:
pos=<x= 2, y=-8, z= 0>, vel=<x=-3, y=-2, z= 1>
pos=<x= 2, y= 1, z= 7>, vel=<x= 2, y= 1, z= 1>
pos=<x= 2, y= 3, z=-6>, vel=<x= 0, y= 2, z=-1>
pos=<x= 2, y=-9, z= 1>, vel=<x= 1, y=-1, z=-1>

After 5 steps:
pos=<x=-1, y=-9, z= 2>, vel=<x=-3, y=-1, z= 2>
pos=<x= 4, y= 1, z= 5>, vel=<x= 2, y= 0, z=-2>
pos=<x= 2, y= 2, z=-4>, vel=<x= 0, y=-1, z= 2>
pos=<x= 3, y=-7, z=-1>, vel=<x= 1, y= 2, z=-2>

After 6 steps:
pos=<x=-1, y=-7, z= 3>, vel=<x= 0, y= 2, z= 1>
pos=<x= 3, y= 0, z= 0>, vel=<x=-1, y=-1, z=-5>
pos=<x= 3, y=-2, z= 1>, vel=<x= 1, y=-4, z= 5>
pos=<x= 3, y=-4, z=-2>, vel=<x= 0, y= 3, z=-1>

After 7 steps:
pos=<x= 2, y=-2, z= 1>, vel=<x= 3, y= 5, z=-2>
pos=<x= 1, y=-4, z=-4>, vel=<x=-2, y=-4, z=-4>
pos=<x= 3, y=-7, z= 5>, vel=<x= 0, y=-5, z= 4>
pos=<x= 2, y= 0, z= 0>, vel=<x=-1, y= 4, z= 2>

After 8 steps:
pos=<x= 5, y= 2, z=-2>, vel=<x= 3, y= 4, z=-3>
pos=<x= 2, y=-7, z=-5>, vel=<x= 1, y=-3, z=-1>
pos=<x= 0, y=-9, z= 6>, vel=<x=-3, y=-2, z= 1>
pos=<x= 1, y= 1, z= 3>, vel=<x=-1, y= 1, z= 3>

After 9 steps:
pos=<x= 5, y= 3, z=-4>, vel=<x= 0, y= 1, z=-2>
pos=<x= 2, y=-9, z=-3>, vel=<x= 0, y=-2, z= 2>
pos=<x= 0, y=-8, z= 4>, vel=<x= 0, y= 1, z=-2>
pos=<x= 1, y= 1, z= 5>, vel=<x= 0, y= 0, z= 2>

After 10 steps:
pos=<x= 2, y= 1, z=-3>, vel=<x=-3, y=-2, z= 1>
pos=<x= 1, y=-8, z= 0>, vel=<x=-1, y= 1, z= 3>
pos=<x= 3, y=-6, z= 1>, vel=<x= 3, y= 2, z=-3>
pos=<x= 2, y= 0, z= 4>, vel=<x= 1, y=-1, z=-1>
Then, it might help to calculate the total energy in the system. The total energy for a single moon is its potential energy multiplied by its kinetic energy. A moon's potential energy is the sum of the absolute values of its x, y, and z position coordinates. A moon's kinetic energy is the sum of the absolute values of its velocity coordinates. Below, each line shows the calculations for a moon's potential energy (pot), kinetic energy (kin), and total energy:

Energy after 10 steps:
pot: 2 + 1 + 3 =  6;   kin: 3 + 2 + 1 = 6;   total:  6 * 6 = 36
pot: 1 + 8 + 0 =  9;   kin: 1 + 1 + 3 = 5;   total:  9 * 5 = 45
pot: 3 + 6 + 1 = 10;   kin: 3 + 2 + 3 = 8;   total: 10 * 8 = 80
pot: 2 + 0 + 4 =  6;   kin: 1 + 1 + 1 = 3;   total:  6 * 3 = 18
Sum of total energy: 36 + 45 + 80 + 18 = 179
In the above example, adding together the total energy for all moons after 10 steps produces the total energy in the system, 179.

Here's a second example:

<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>
Every ten steps of simulation for 100 steps produces:

After 0 steps:
pos=<x= -8, y=-10, z=  0>, vel=<x=  0, y=  0, z=  0>
pos=<x=  5, y=  5, z= 10>, vel=<x=  0, y=  0, z=  0>
pos=<x=  2, y= -7, z=  3>, vel=<x=  0, y=  0, z=  0>
pos=<x=  9, y= -8, z= -3>, vel=<x=  0, y=  0, z=  0>

After 10 steps:
pos=<x= -9, y=-10, z=  1>, vel=<x= -2, y= -2, z= -1>
pos=<x=  4, y= 10, z=  9>, vel=<x= -3, y=  7, z= -2>
pos=<x=  8, y=-10, z= -3>, vel=<x=  5, y= -1, z= -2>
pos=<x=  5, y=-10, z=  3>, vel=<x=  0, y= -4, z=  5>

After 20 steps:
pos=<x=-10, y=  3, z= -4>, vel=<x= -5, y=  2, z=  0>
pos=<x=  5, y=-25, z=  6>, vel=<x=  1, y=  1, z= -4>
pos=<x= 13, y=  1, z=  1>, vel=<x=  5, y= -2, z=  2>
pos=<x=  0, y=  1, z=  7>, vel=<x= -1, y= -1, z=  2>

After 30 steps:
pos=<x= 15, y= -6, z= -9>, vel=<x= -5, y=  4, z=  0>
pos=<x= -4, y=-11, z=  3>, vel=<x= -3, y=-10, z=  0>
pos=<x=  0, y= -1, z= 11>, vel=<x=  7, y=  4, z=  3>
pos=<x= -3, y= -2, z=  5>, vel=<x=  1, y=  2, z= -3>

After 40 steps:
pos=<x= 14, y=-12, z= -4>, vel=<x= 11, y=  3, z=  0>
pos=<x= -1, y= 18, z=  8>, vel=<x= -5, y=  2, z=  3>
pos=<x= -5, y=-14, z=  8>, vel=<x=  1, y= -2, z=  0>
pos=<x=  0, y=-12, z= -2>, vel=<x= -7, y= -3, z= -3>

After 50 steps:
pos=<x=-23, y=  4, z=  1>, vel=<x= -7, y= -1, z=  2>
pos=<x= 20, y=-31, z= 13>, vel=<x=  5, y=  3, z=  4>
pos=<x= -4, y=  6, z=  1>, vel=<x= -1, y=  1, z= -3>
pos=<x= 15, y=  1, z= -5>, vel=<x=  3, y= -3, z= -3>

After 60 steps:
pos=<x= 36, y=-10, z=  6>, vel=<x=  5, y=  0, z=  3>
pos=<x=-18, y= 10, z=  9>, vel=<x= -3, y= -7, z=  5>
pos=<x=  8, y=-12, z= -3>, vel=<x= -2, y=  1, z= -7>
pos=<x=-18, y= -8, z= -2>, vel=<x=  0, y=  6, z= -1>

After 70 steps:
pos=<x=-33, y= -6, z=  5>, vel=<x= -5, y= -4, z=  7>
pos=<x= 13, y= -9, z=  2>, vel=<x= -2, y= 11, z=  3>
pos=<x= 11, y= -8, z=  2>, vel=<x=  8, y= -6, z= -7>
pos=<x= 17, y=  3, z=  1>, vel=<x= -1, y= -1, z= -3>

After 80 steps:
pos=<x= 30, y= -8, z=  3>, vel=<x=  3, y=  3, z=  0>
pos=<x= -2, y= -4, z=  0>, vel=<x=  4, y=-13, z=  2>
pos=<x=-18, y= -7, z= 15>, vel=<x= -8, y=  2, z= -2>
pos=<x= -2, y= -1, z= -8>, vel=<x=  1, y=  8, z=  0>

After 90 steps:
pos=<x=-25, y= -1, z=  4>, vel=<x=  1, y= -3, z=  4>
pos=<x=  2, y= -9, z=  0>, vel=<x= -3, y= 13, z= -1>
pos=<x= 32, y= -8, z= 14>, vel=<x=  5, y= -4, z=  6>
pos=<x= -1, y= -2, z= -8>, vel=<x= -3, y= -6, z= -9>

After 100 steps:
pos=<x=  8, y=-12, z= -9>, vel=<x= -7, y=  3, z=  0>
pos=<x= 13, y= 16, z= -3>, vel=<x=  3, y=-11, z= -5>
pos=<x=-29, y=-11, z= -1>, vel=<x= -3, y=  7, z=  4>
pos=<x= 16, y=-13, z= 23>, vel=<x=  7, y=  1, z=  1>

Energy after 100 steps:
pot:  8 + 12 +  9 = 29;   kin: 7 +  3 + 0 = 10;   total: 29 * 10 = 290
pot: 13 + 16 +  3 = 32;   kin: 3 + 11 + 5 = 19;   total: 32 * 19 = 608
pot: 29 + 11 +  1 = 41;   kin: 3 +  7 + 4 = 14;   total: 41 * 14 = 574
pot: 16 + 13 + 23 = 52;   kin: 7 +  1 + 1 =  9;   total: 52 *  9 = 468
Sum of total energy: 290 + 608 + 574 + 468 = 1940
What is the total energy in the system after simulating the moons given in your scan for 1000 steps?

Your puzzle answer was 7179.

--- Part Two ---
All this drifting around in space makes you wonder about the nature of the universe. Does history really repeat itself? You're curious whether the moons will ever return to a previous state.

Determine the number of steps that must occur before all of the moons' positions and velocities exactly match a previous point in time.

For example, the first example above takes 2772 steps before they exactly match a previous point in time; it eventually returns to the initial state:

After 0 steps:
pos=<x= -1, y=  0, z=  2>, vel=<x=  0, y=  0, z=  0>
pos=<x=  2, y=-10, z= -7>, vel=<x=  0, y=  0, z=  0>
pos=<x=  4, y= -8, z=  8>, vel=<x=  0, y=  0, z=  0>
pos=<x=  3, y=  5, z= -1>, vel=<x=  0, y=  0, z=  0>

After 2770 steps:
pos=<x=  2, y= -1, z=  1>, vel=<x= -3, y=  2, z=  2>
pos=<x=  3, y= -7, z= -4>, vel=<x=  2, y= -5, z= -6>
pos=<x=  1, y= -7, z=  5>, vel=<x=  0, y= -3, z=  6>
pos=<x=  2, y=  2, z=  0>, vel=<x=  1, y=  6, z= -2>

After 2771 steps:
pos=<x= -1, y=  0, z=  2>, vel=<x= -3, y=  1, z=  1>
pos=<x=  2, y=-10, z= -7>, vel=<x= -1, y= -3, z= -3>
pos=<x=  4, y= -8, z=  8>, vel=<x=  3, y= -1, z=  3>
pos=<x=  3, y=  5, z= -1>, vel=<x=  1, y=  3, z= -1>

After 2772 steps:
pos=<x= -1, y=  0, z=  2>, vel=<x=  0, y=  0, z=  0>
pos=<x=  2, y=-10, z= -7>, vel=<x=  0, y=  0, z=  0>
pos=<x=  4, y= -8, z=  8>, vel=<x=  0, y=  0, z=  0>
pos=<x=  3, y=  5, z= -1>, vel=<x=  0, y=  0, z=  0>
Of course, the universe might last for a very long time before repeating. Here's a copy of the second example from above:

<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>
This set of initial positions takes 4686774924 steps before it repeats a previous state! Clearly, you might need to find a more efficient way to simulate the universe.

How many steps does it take to reach the first state that exactly matches a previous state?

Your puzzle answer was 428576638953552.

Both parts of this puzzle are complete! They provide two gold stars: **
 */