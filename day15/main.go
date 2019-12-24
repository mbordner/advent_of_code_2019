package main

import (
	"github.com/mbordner/advent_of_code_2019/day15/game"
	"github.com/mbordner/advent_of_code_2019/day15/geom"
	"github.com/mbordner/advent_of_code_2019/day15/graph"
	"github.com/mbordner/advent_of_code_2019/day15/intcode"
	"fmt"
	"strings"
	"time"
)

func main() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	gameoutput := make(chan string, 1)
	gameinput := make(chan string, 1)
	movecomplete := make(chan string, 1)
	compquit := make(chan string, 1)
	quit := make(chan string, 1)

	program := getProgram()

	intCodeComputer := intcode.NewIntCodeComputer(program, in, out, compquit, true)
	gameUI := game.NewGame(gameinput, gameoutput, movecomplete, compquit, quit)
	gameGraph := graph.NewGraph()

	go intCodeComputer.Execute()

	// discover the map
	for {
		dir := gameGraph.GetNextDirection()
		if dir != geom.Unknown {
			gameUI.SetLastDir(dir)
			in <- fmt.Sprintf("%d", dir)
			response := <-out
			switch response {
			case "0":
				gameGraph.SetTraversable(false)
			case "1":
				gameGraph.SetTraversable(true)
			case "2":
				gameGraph.SetTraversable(true)
				gameGraph.SetGoal()
			}
			gameoutput <- response
			<-movecomplete
			intCodeComputer.OutputProcessed()
		} else {
			break
		}
	}
	gameGraph.RemoveImpassable()

	// get shortest path to goal
	path := gameGraph.GenerateShortestPath()
	for _, p := range path {
		o := gameUI.GetObject(p.Pos)
		o.ShortestPath = true
	}
	gameUI.Refresh()

	generations := make([][]*graph.Node,0,len(path))

	// fill area with oxygen
	go func() {
		goal := gameGraph.Goal

		generation := []*graph.Node{goal}
		for len(generation) > 0 {
			generations = append(generations,generation)
			for _, n := range generation {
				gameUI.GetObject(n.Pos).FillWithOxygen()
			}
			gameUI.Refresh()

			time.Sleep(time.Duration(20)*time.Millisecond)

			nextGeneration := make([]*graph.Node,0,len(generation)*4)
			for _,n := range generation {
				if n.North != nil && gameUI.GetObject(n.North.Pos).HasOxygen == false {
					nextGeneration = append(nextGeneration,n.North)
				}
				if n.South != nil && gameUI.GetObject(n.South.Pos).HasOxygen == false {
					nextGeneration = append(nextGeneration,n.South)
				}
				if n.East != nil && gameUI.GetObject(n.East.Pos).HasOxygen == false {
					nextGeneration = append(nextGeneration,n.East)
				}
				if n.West != nil && gameUI.GetObject(n.West.Pos).HasOxygen == false {
					nextGeneration = append(nextGeneration,n.West)
				}
			}

			generation = nextGeneration
		}

	}()


programLoop:
	for {
		select {
		case move := <-gameinput:
			in <- move

		case response := <-out:
			gameoutput <- response

		case <-movecomplete:
			intCodeComputer.OutputProcessed()

		case <-quit:
			fmt.Println("program exited")
			break programLoop

		}
	}

	gameUI.Shutdown()

	fmt.Println("fewest number of movement commands required to move the repair droid from its starting position to the location of the oxygen system: ",len(path)-1)
	fmt.Println("number of minutes it will take to fill the area with oxygen: ",len(generations)-1)
}

func getProgram() []string {
	program := `3,1033,1008,1033,1,1032,1005,1032,31,1008,1033,2,1032,1005,1032,58,1008,1033,3,1032,1005,1032,81,1008,1033,4,1032,1005,1032,104,99,101,0,1034,1039,1001,1036,0,1041,1001,1035,-1,1040,1008,1038,0,1043,102,-1,1043,1032,1,1037,1032,1042,1106,0,124,1001,1034,0,1039,1001,1036,0,1041,1001,1035,1,1040,1008,1038,0,1043,1,1037,1038,1042,1106,0,124,1001,1034,-1,1039,1008,1036,0,1041,102,1,1035,1040,101,0,1038,1043,102,1,1037,1042,1105,1,124,1001,1034,1,1039,1008,1036,0,1041,101,0,1035,1040,1001,1038,0,1043,101,0,1037,1042,1006,1039,217,1006,1040,217,1008,1039,40,1032,1005,1032,217,1008,1040,40,1032,1005,1032,217,1008,1039,1,1032,1006,1032,165,1008,1040,3,1032,1006,1032,165,1101,0,2,1044,1105,1,224,2,1041,1043,1032,1006,1032,179,1102,1,1,1044,1106,0,224,1,1041,1043,1032,1006,1032,217,1,1042,1043,1032,1001,1032,-1,1032,1002,1032,39,1032,1,1032,1039,1032,101,-1,1032,1032,101,252,1032,211,1007,0,45,1044,1105,1,224,1101,0,0,1044,1106,0,224,1006,1044,247,1002,1039,1,1034,1002,1040,1,1035,1001,1041,0,1036,1002,1043,1,1038,102,1,1042,1037,4,1044,1106,0,0,7,39,95,7,98,8,11,47,17,33,19,4,29,41,87,34,59,22,75,5,1,46,41,29,32,11,55,25,53,41,77,27,52,33,41,65,72,24,43,83,72,3,14,92,2,43,82,30,87,19,94,47,91,10,8,67,24,4,68,85,63,4,93,29,55,34,23,65,40,3,36,90,57,97,37,2,65,8,1,16,83,93,67,44,71,97,27,70,76,20,40,90,36,73,27,89,57,13,66,37,95,76,26,84,33,48,34,86,85,30,81,6,61,33,83,84,22,21,67,27,11,49,28,69,41,60,98,6,69,41,54,82,18,37,65,10,42,47,41,2,72,16,66,39,93,37,2,41,52,49,20,78,30,7,38,15,40,81,21,14,82,44,48,7,96,33,36,70,52,18,71,1,81,66,47,1,38,78,80,38,63,53,80,16,58,55,93,31,89,36,36,78,65,71,34,83,4,55,60,29,10,30,84,15,59,31,96,16,21,58,26,38,35,58,50,16,46,25,26,82,59,12,11,98,4,17,42,66,83,72,23,14,92,22,9,5,87,5,79,85,19,87,71,28,61,32,56,92,56,19,78,94,39,24,73,58,28,37,81,11,99,25,46,73,44,5,22,41,76,55,84,31,16,36,65,84,40,29,81,66,16,94,23,54,23,29,51,20,25,23,69,44,23,18,99,80,55,39,10,71,7,33,63,94,93,62,26,35,25,50,61,39,84,38,54,43,56,23,67,17,70,34,23,90,93,24,46,60,31,46,33,53,81,10,62,23,89,86,43,39,73,82,38,9,61,42,66,68,30,28,95,4,25,54,22,21,80,32,61,13,6,66,47,59,4,31,59,17,87,72,30,72,51,30,30,62,43,53,88,42,48,13,21,80,8,30,61,14,77,22,27,60,87,30,65,14,33,76,67,9,95,26,84,40,21,52,11,86,23,30,86,57,28,6,69,4,11,63,21,2,65,51,39,58,82,16,51,96,23,3,44,21,62,31,38,47,73,30,29,94,24,14,88,1,51,72,42,57,48,63,33,95,78,15,17,68,64,61,10,31,58,68,36,15,52,19,13,26,38,72,41,66,15,56,88,18,98,87,15,43,89,96,3,94,55,25,26,27,6,48,3,29,90,88,6,18,29,88,90,43,3,81,61,16,31,93,42,26,46,31,56,66,17,76,37,15,50,33,81,16,10,83,87,37,39,92,80,62,6,59,77,9,32,91,61,97,24,44,62,61,11,36,94,59,54,34,23,67,18,86,31,39,77,73,44,67,27,57,5,54,65,29,21,81,2,65,39,24,82,6,55,33,97,72,35,16,85,19,28,57,94,21,15,86,5,52,53,39,69,20,32,52,5,86,95,44,47,77,9,57,14,62,49,54,7,70,29,16,42,87,99,30,36,67,68,14,42,73,4,87,97,39,61,18,11,39,77,83,17,83,27,1,72,30,21,95,38,35,96,15,78,27,66,40,4,95,90,94,4,20,63,71,19,54,11,28,96,46,13,42,94,84,9,22,79,37,14,50,13,58,64,90,30,69,18,20,90,4,21,31,95,88,22,81,36,20,11,82,59,95,38,43,72,3,78,38,33,62,48,36,22,16,3,87,53,91,37,12,19,49,18,25,14,67,78,79,9,70,88,34,98,38,8,90,98,56,13,26,34,82,77,40,97,82,63,32,57,26,58,53,29,56,3,62,17,78,67,69,33,49,62,47,36,60,9,81,12,96,6,78,86,98,34,70,41,87,86,47,15,46,36,49,20,76,31,48,1,68,19,96,0,0,21,21,1,10,1,0,0,0,0,0,0`
	return strings.Split(program, ",")
}

/**
--- Day 15: Oxygen System ---
Out here in deep space, many things can go wrong. Fortunately, many of those things have indicator lights. Unfortunately, one of those lights is lit: the oxygen system for part of the ship has failed!

According to the readouts, the oxygen system must have failed days ago after a rupture in oxygen tank two; that section of the ship was automatically sealed once oxygen levels went dangerously low. A single remotely-operated repair droid is your only option for fixing the oxygen system.

The Elves' care package included an Intcode program (your puzzle input) that you can use to remotely control the repair droid. By running that program, you can direct the repair droid to the oxygen system and fix the problem.

The remote control program executes the following steps in a loop forever:

Accept a movement command via an input instruction.
Send the movement command to the repair droid.
Wait for the repair droid to finish the movement operation.
Report on the status of the repair droid via an output instruction.
Only four movement commands are understood: north (1), south (2), west (3), and east (4). Any other command is invalid. The movements differ in direction, but not in distance: in a long enough east-west hallway, a series of commands like 4,4,4,4,3,3,3,3 would leave the repair droid back where it started.

The repair droid can reply with any of the following status codes:

0: The repair droid hit a wall. Its position has not changed.
1: The repair droid has moved one step in the requested direction.
2: The repair droid has moved one step in the requested direction; its new position is the location of the oxygen system.
You don't know anything about the area around the repair droid, but you can figure it out by watching the status codes.

For example, we can draw the area using D for the droid, # for walls, . for locations the droid can traverse, and empty space for unexplored locations. Then, the initial state looks like this:



   D


To make the droid go north, send it 1. If it replies with 0, you know that location is a wall and that the droid didn't move:


   #
   D


To move east, send 4; a reply of 1 means the movement was successful:


   #
   .D


Then, perhaps attempts to move north (1), south (2), and east (4) are all met with replies of 0:


   ##
   .D#
    #

Now, you know the repair droid is in a dead end. Backtrack with 3 (which you already know will get a reply of 1 because you already know that location is open):


   ##
   D.#
    #

Then, perhaps west (3) gets a reply of 0, south (2) gets a reply of 1, south again (2) gets a reply of 0, and then west (3) gets a reply of 2:


   ##
  #..#
  D.#
   #
Now, because of the reply of 2, you know you've found the oxygen system! In this example, it was only 2 moves away from the repair droid's starting position.

What is the fewest number of movement commands required to move the repair droid from its starting position to the location of the oxygen system?

Your puzzle answer was 246.

--- Part Two ---
You quickly repair the oxygen system; oxygen gradually fills the area.

Oxygen starts in the location containing the repaired oxygen system. It takes one minute for oxygen to spread to all open locations that are adjacent to a location that already contains oxygen. Diagonal locations are not adjacent.

In the example above, suppose you've used the droid to explore the area fully and have the following map (where locations that currently contain oxygen are marked O):

 ##
#..##
#.#..#
#.O.#
 ###
Initially, the only location which contains oxygen is the location of the repaired oxygen system. However, after one minute, the oxygen spreads to all open (.) locations that are adjacent to a location containing oxygen:

 ##
#..##
#.#..#
#OOO#
 ###
After a total of two minutes, the map looks like this:

 ##
#..##
#O#O.#
#OOO#
 ###
After a total of three minutes:

 ##
#O.##
#O#OO#
#OOO#
 ###
And finally, the whole region is full of oxygen after a total of four minutes:

 ##
#OO##
#O#OO#
#OOO#
 ###
So, in this example, all locations contain oxygen after 4 minutes.

Use the repair droid to get a complete map of the area. How many minutes will it take to fill with oxygen?

Your puzzle answer was 376.

Both parts of this puzzle are complete! They provide two gold stars: **
 */