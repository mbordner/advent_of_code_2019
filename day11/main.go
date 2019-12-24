package main

import (
	"errors"
	"fmt"
	"strings"
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

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d",p.X,p.Y)
}

type Direction int
type Color int

func (c Color) String() string {
	switch c {
	case White:
		return "white"
	case Black:
		return "black"
	}
	panic(errors.New("unknown color"))
}

func (d Direction) String() string {
	switch d {
	case Left:
		return "left"
	case Right:
		return "right"
	case Down:
		return "down"
	case Up:
		return "up"
	}
	panic(errors.New("unknown direction"))
}

const (
	Left Direction = iota
	Right
	Up
	Down
)

const (
	Black Color = iota
	White
)

var (
	panels = make(map[Pos][]Color)
)

func getDirectionAfterTurn(current Direction, turn Direction) Direction {
	switch current {
	case Left:
		if turn == Left {
			return Down
		} else if turn == Right {
			return Up
		}
	case Right:
		if turn == Left {
			return Up
		} else if turn == Right {
			return Down
		}
	case Down:
		if turn == Left {
			return Right
		} else if turn == Right {
			return Left
		}
	case Up:
		if turn == Left {
			return Left
		} else if turn == Right {
			return Right
		}
	}
	panic(errors.New("undefined current direction"))
}

func getPanelColor(pos Pos) Color {
	if colors, ok := panels[pos]; ok {
		return colors[len(colors)-1]
	}
	return Black
}

type Robot struct {
	pos Pos
	dir Direction
}

func (r *Robot) GetPosition() Pos {
	return r.pos
}

func (r *Robot) GetDirection() Direction {
	return r.dir
}

func (r *Robot) TurnAndAdvance(dir Direction) {
	newDir := getDirectionAfterTurn(r.dir, dir)
	switch newDir {
	case Left:
		r.pos.X--
	case Right:
		r.pos.X++
	case Down:
		r.pos.Y--
	case Up:
		r.pos.Y++
	}
	r.dir = newDir
}

func (r *Robot) Paint(color Color) {
	if _, ok := panels[r.pos]; !ok {
		panels[r.pos] = make([]Color, 0, 10)
	}
	panels[r.pos] = append(panels[r.pos], color)
}

func NewRobot() *Robot {
	r := new(Robot)
	r.dir = Up
	return r
}


func main() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	quit := make(chan string, 1)

	c := NewIntCodeComputer(getProgram(), in, out, quit, true)
	r := NewRobot()

	r.Paint(White)

	// run the int code computer against the program,
	// on input it will want the color of the tile the robot is currently sitting
	// after processing the input it will spit out two numbers in succession,
	// the color of the panel to paint for the panel the robot is currently on,
	// and then a direction to turn, .. after turning right or left, the robot
	// needs to advance forward one panel.
	// -- we assume the robot always starts at position 0,0 and will calculate
	// a path from this origin
	// -- after the path is constructed we can get a bounding box for the entire
	// path, generate a [][]character array and paint the panels of the path
	go c.Execute()

programLoop:
	for {
		select {
		case <-in:

			color := getPanelColor(r.GetPosition())
			//fmt.Printf("current color: %s\n",color)

			input := "1"
			if color == Black {
				input = "0"
			}

			in <- input

		case out1 := <-out:
			c.OutputProcessed()
			out2 := <-out

			color := Black
			if out1 == "1" {
				color = White
			}

			direction := Left
			if out2 == "1" {
				direction = Right
			}

			//fmt.Printf("painting %s\n",color)
			r.Paint(color)
			r.TurnAndAdvance(direction)
			//fmt.Printf("turning %s, new direction: %s, new position: %s\n",direction,r.dir,r.pos)

			c.OutputProcessed()

		case <-quit:
			fmt.Println("program exited")
			break programLoop
		}
	}

	fmt.Println("number of panels painted: ", len(panels))

	// calculate the size of the matrix of characters we'll need to print
	bb := new(BoundingBox)
	for pos := range panels {
		if pos.X < bb.xMin {
			bb.xMin = pos.X
		}
		if pos.X > bb.xMax {
			bb.xMax = pos.X
		}
		if pos.Y < bb.yMin {
			bb.yMin = pos.Y
		}
		if pos.Y > bb.yMax {
			bb.yMax = pos.Y
		}
	}

	fmt.Println(*bb)

	cols := bb.xMax - bb.xMin + 1
	rows := bb.yMax - bb.yMin + 1

	chars := make([][]byte, rows, rows)
	for i := range chars {
		chars[i] = make([]byte, cols, cols)
		for j := range chars[i] {
			chars[i][j] = ' ' // initialize to black
		}
	}

	// write out the white panels
	for pos, colors := range panels {
		// adjust by offsets, we have some positions left and below origin with negative numbers
		r := pos.Y - bb.yMin
		c := pos.X - bb.xMin

		switch colors[len(colors)-1] {
		case White:
			chars[r][c] = '#'
		}
	}

	// flip rows
	for i := 0; i < len(chars)/2; i++ {
		chars[i], chars[len(chars)-1-i] = chars[len(chars)-1-i], chars[i]
	}

	// print data
	for _, row := range chars {
		fmt.Println(string(row))
	}

}

func getProgram() []string {
	return strings.Split(`3,8,1005,8,320,1106,0,11,0,0,0,104,1,104,0,3,8,102,-1,8,10,101,1,10,10,4,10,1008,8,1,10,4,10,1001,8,0,29,2,101,10,10,3,8,102,-1,8,10,1001,10,1,10,4,10,108,1,8,10,4,10,101,0,8,54,2,3,16,10,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,0,10,4,10,102,1,8,81,1006,0,75,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,0,8,10,4,10,101,0,8,105,3,8,102,-1,8,10,1001,10,1,10,4,10,1008,8,1,10,4,10,1001,8,0,128,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,0,8,10,4,10,102,1,8,149,1,105,5,10,1,105,20,10,3,8,102,-1,8,10,101,1,10,10,4,10,108,0,8,10,4,10,1002,8,1,179,1,101,1,10,2,109,8,10,1006,0,74,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,1,10,4,10,1001,8,0,213,1006,0,60,2,1105,9,10,1,1005,11,10,3,8,1002,8,-1,10,101,1,10,10,4,10,108,1,8,10,4,10,1002,8,1,245,1,6,20,10,1,1103,11,10,2,6,11,10,2,1103,0,10,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,0,10,4,10,1002,8,1,284,2,1103,12,10,2,1104,14,10,2,1004,12,10,2,1009,4,10,101,1,9,9,1007,9,968,10,1005,10,15,99,109,642,104,0,104,1,21102,1,48063419288,1,21102,1,337,0,1105,1,441,21101,0,846927340300,1,21101,0,348,0,1105,1,441,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,21102,1,235245104151,1,21102,395,1,0,1105,1,441,21102,29032123584,1,1,21101,0,406,0,1105,1,441,3,10,104,0,104,0,3,10,104,0,104,0,21101,0,709047878500,1,21101,429,0,0,1106,0,441,21101,868402070284,0,1,21102,1,440,0,1105,1,441,99,109,2,22102,1,-1,1,21101,40,0,2,21101,0,472,3,21102,462,1,0,1105,1,505,109,-2,2106,0,0,0,1,0,0,1,109,2,3,10,204,-1,1001,467,468,483,4,0,1001,467,1,467,108,4,467,10,1006,10,499,1102,1,0,467,109,-2,2106,0,0,0,109,4,2101,0,-1,504,1207,-3,0,10,1006,10,522,21101,0,0,-3,22101,0,-3,1,21202,-2,1,2,21101,1,0,3,21102,541,1,0,1106,0,546,109,-4,2106,0,0,109,5,1207,-3,1,10,1006,10,569,2207,-4,-2,10,1006,10,569,21202,-4,1,-4,1105,1,637,22102,1,-4,1,21201,-3,-1,2,21202,-2,2,3,21101,588,0,0,1105,1,546,22102,1,1,-4,21101,0,1,-1,2207,-4,-2,10,1006,10,607,21101,0,0,-1,22202,-2,-1,-2,2107,0,-3,10,1006,10,629,21201,-1,0,1,21102,629,1,0,106,0,504,21202,-2,-1,-2,22201,-4,-2,-4,109,-5,2105,1,0`, ",")
}

/**
--- Day 11: Space Police ---
On the way to Jupiter, you're pulled over by the Space Police.

"Attention, unmarked spacecraft! You are in violation of Space Law! All spacecraft must have a clearly visible registration identifier! You have 24 hours to comply or be sent to Space Jail!"

Not wanting to be sent to Space Jail, you radio back to the Elves on Earth for help. Although it takes almost three hours for their reply signal to reach you, they send instructions for how to power up the emergency hull painting robot and even provide a small Intcode program (your puzzle input) that will cause it to paint your ship appropriately.

There's just one problem: you don't have an emergency hull painting robot.

You'll need to build a new emergency hull painting robot. The robot needs to be able to move around on the grid of square panels on the side of your ship, detect the color of its current panel, and paint its current panel black or white. (All of the panels are currently black.)

The Intcode program will serve as the brain of the robot. The program uses input instructions to access the robot's camera: provide 0 if the robot is over a black panel or 1 if the robot is over a white panel. Then, the program will output two values:

First, it will output a value indicating the color to paint the panel the robot is over: 0 means to paint the panel black, and 1 means to paint the panel white.
Second, it will output a value indicating the direction the robot should turn: 0 means it should turn left 90 degrees, and 1 means it should turn right 90 degrees.
After the robot turns, it should always move forward exactly one panel. The robot starts facing up.

The robot will continue running for a while like this and halt when it is finished drawing. Do not restart the Intcode computer inside the robot during this process.

For example, suppose the robot is about to start running. Drawing black panels as ., white panels as #, and the robot pointing the direction it is facing (< ^ > v), the initial state and region near the robot looks like this:

.....
.....
..^..
.....
.....
The panel under the robot (not visible here because a ^ is shown instead) is also black, and so any input instructions at this point should be provided 0. Suppose the robot eventually outputs 1 (paint white) and then 0 (turn left). After taking these actions and moving forward one panel, the region now looks like this:

.....
.....
.<#..
.....
.....
Input instructions should still be provided 0. Next, the robot might output 0 (paint black) and then 0 (turn left):

.....
.....
..#..
.v...
.....
After more outputs (1,0, 1,0):

.....
.....
..^..
.##..
.....
The robot is now back where it started, but because it is now on a white panel, input instructions should be provided 1. After several more outputs (0,1, 1,0, 1,0), the area looks like this:

.....
..<#.
...#.
.##..
.....
Before you deploy the robot, you should probably have an estimate of the area it will cover: specifically, you need to know the number of panels it paints at least once, regardless of color. In the example above, the robot painted 6 panels at least once. (It painted its starting panel twice, but that panel is still only counted once; it also never painted the panel it ended on.)

Build a new emergency hull painting robot and run the Intcode program on it. How many panels does it paint at least once?

Your puzzle answer was 1909.

--- Part Two ---
You're not sure what it's trying to paint, but it's definitely not a registration identifier. The Space Police are getting impatient.

Checking your external ship cameras again, you notice a white panel marked "emergency hull painting robot starting panel". The rest of the panels are still black, but it looks like the robot was expecting to start on a white panel, not a black one.

Based on the Space Law Space Brochure that the Space Police attached to one of your windows, a valid registration identifier is always eight capital letters. After starting the robot on a single white panel instead, what registration identifier does it paint on your hull?

Your puzzle answer was JUFEKHPH.

Both parts of this puzzle are complete! They provide two gold stars: **
 */