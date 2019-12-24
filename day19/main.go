package main

import (
	"github.com/mbordner/advent_of_code_2019/day19/intcode"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

// for each new row y, from 0, find the first x in the row (from 0) where the tractor beam is affecting
// save that x spot, count to the right, if we are no longer affecting, repeat on next row, otherwise if we reach our limit, count down from saved spot, if we reach our limit, we found the spot, if we stop being affected, advance right and repeat from save that x

func part2() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	quit := make(chan string, 1)

	program := getProgram("program1.txt")

	intCodeComputer := intcode.NewIntCodeComputer(program, in, out, quit, true)

	compute := func(x, y int) (result bool) {
		intCodeComputer.Reset()
		go intCodeComputer.Execute()
		in <- fmt.Sprintf("%d", x)
		in <- fmt.Sprintf("%d", y)

		output := <-out
		if output == "1" {
			result = true
		}
		intCodeComputer.OutputProcessed()
		<-quit
		return
	}

	row := 1000

	const SHIP = 100

outter:
	for {
		xStart := 0
		for !compute(xStart, row) {
			xStart++
		}
		xStop := xStart + 1
		for compute(xStop, row) {
			xStop++
		}
		fmt.Println("row ", row, " x starts at ", xStart, " and ends at ", xStop)
		if xStop-xStart >= SHIP {
			for x := xStart; xStop-x >= SHIP; x++ {
				y := row + SHIP - 1
				fmt.Println("row ", row, " x starts at ", xStart, " and ends at ", xStop, " - checking ", x, ",", y)
				if compute(x, y) {
					fmt.Println("answer is: ", x*10000+row)
					break outter
				}
			}
		}
		row++
	}

}

func part1() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	quit := make(chan string, 1)

	program := getProgram("program1.txt")

	intCodeComputer := intcode.NewIntCodeComputer(program, in, out, quit, true)

	gameMap := getGameMap(50, 50)

	for y, row := range gameMap {
		for x := range row {
			intCodeComputer.Reset()
			go intCodeComputer.Execute()
			in <- fmt.Sprintf("%d", x)
			in <- fmt.Sprintf("%d", y)

			output := <-out
			if output == "1" {
				gameMap[y][x] = byte('#')
			}
			intCodeComputer.OutputProcessed()
			<-quit
			fmt.Print(string(gameMap[y][x]))
		}
		fmt.Print("\n")
	}

	count := 0
	for y := range gameMap {
		count += strings.Count(string(gameMap[y]), "#")
	}

	fmt.Println("number of positions affected: ", count)

}

func main() {
	//part1()
	part2()
}

func getGameMap(x, y int) [][]byte {
	gameMap := make([][]byte, y, y)
	for j := range gameMap {
		gameMap[j] = bytes.Repeat([]byte{byte('.')}, x)
	}
	return gameMap
}

func getProgram(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		panic(err)
	}

	if bytesread != int(filesize) {
		panic(errors.New("didn't read all of the file"))
	}

	return strings.Split(string(buffer), ",")
}

/**
--- Day 19: Tractor Beam ---
Unsure of the state of Santa's ship, you borrowed the tractor beam technology from Triton. Time to test it out.

When you're safely away from anything else, you activate the tractor beam, but nothing happens. It's hard to tell whether it's working if there's nothing to use it on. Fortunately, your ship's drone system can be configured to deploy a drone to specific coordinates and then check whether it's being pulled. There's even an Intcode program (your puzzle input) that gives you access to the drone system.

The program uses two input instructions to request the X and Y position to which the drone should be deployed. Negative numbers are invalid and will confuse the drone; all numbers should be zero or positive.

Then, the program will output whether the drone is stationary (0) or being pulled by something (1). For example, the coordinate X=0, Y=0 is directly in front of the tractor beam emitter, so the drone control program will always report 1 at that location.

To better understand the tractor beam, it is important to get a good picture of the beam itself. For example, suppose you scan the 10x10 grid of points closest to the emitter:

       X
  0->      9
 0#.........
 |.#........
 v..##......
  ...###....
  ....###...
Y .....####.
  ......####
  ......####
  .......###
 9........##
In this example, the number of points affected by the tractor beam in the 10x10 area closest to the emitter is 27.

However, you'll need to scan a larger area to understand the shape of the beam. How many points are affected by the tractor beam in the 50x50 area closest to the emitter? (For each of X and Y, this will be 0 through 49.)

Your puzzle answer was 162.

--- Part Two ---
You aren't sure how large Santa's ship is. You aren't even sure if you'll need to use this thing on Santa's ship, but it doesn't hurt to be prepared. You figure Santa's ship might fit in a 100x100 square.

The beam gets wider as it travels away from the emitter; you'll need to be a minimum distance away to fit a square of that size into the beam fully. (Don't rotate the square; it should be aligned to the same axes as the drone grid.)

For example, suppose you have the following tractor beam readings:

#.......................................
.#......................................
..##....................................
...###..................................
....###.................................
.....####...............................
......#####.............................
......######............................
.......#######..........................
........########........................
.........#########......................
..........#########.....................
...........##########...................
...........############.................
............############................
.............#############..............
..............##############............
...............###############..........
................###############.........
................#################.......
.................########OOOOOOOOOO.....
..................#######OOOOOOOOOO#....
...................######OOOOOOOOOO###..
....................#####OOOOOOOOOO#####
.....................####OOOOOOOOOO#####
.....................####OOOOOOOOOO#####
......................###OOOOOOOOOO#####
.......................##OOOOOOOOOO#####
........................#OOOOOOOOOO#####
.........................OOOOOOOOOO#####
..........................##############
..........................##############
...........................#############
............................############
.............................###########
In this example, the 10x10 square closest to the emitter that fits entirely within the tractor beam has been marked O. Within it, the point closest to the emitter (the only highlighted O) is at X=25, Y=20.

Find the 100x100 square closest to the emitter that fits entirely within the tractor beam; within that square, find the point closest to the emitter. What value do you get if you take that point's X coordinate, multiply it by 10000, then add the point's Y coordinate? (In the example above, this would be 250020.)

Your puzzle answer was 13021056.

Both parts of this puzzle are complete! They provide two gold stars: **
 */
