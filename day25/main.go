package main

import (
	"errors"
	"fmt"
	tty "github.com/mattn/go-tty"
	"github.com/mbordner/advent_of_code_2019/day25/game"
	"github.com/mbordner/advent_of_code_2019/day25/intcode"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
== Security Checkpoint ==
In the next room, a pressure-sensitive floor will verify your identity.

Doors here lead:
- north
- east

Items here:
- planetoid
- festive hat
- space heater
- loom
- space law space brochure
- sand
- pointer
- wreath

*/

func main() {
	inChan := make(chan string, 40)
	outChan := make(chan string, 1)
	quitChan := make(chan string, 1)

	prompt := ""

	computer := intcode.NewIntCodeComputer(getProgram("program.txt"), inChan, outChan, quitChan, true, &prompt)

	promptChan := computer.GetPromptChannel()

	//computer.Load("./game.json")

	go computer.Execute()

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	outTTY := tty.Output()

	var bot *game.Game
	bot = game.NewGame([]string{"infinite loop","giant electromagnet","molten lava","escape pod","photons"}, "Security Checkpoint")

	inputQueue := make([]byte, 0, 1024)

programLoop:
	for {
		select {
		case <-promptChan:
			var r rune
			var err error
			if len(inputQueue) > 0 {
				r = rune(inputQueue[0])
				inputQueue = inputQueue[1:]
			} else {
				r, err = tty.ReadRune()
			}

			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(outTTY, "%c", r)

			b := byte(r)
			if b == 13 {
				b = 10
			}
			inChan <- fmt.Sprintf("%d", b)

		case output := <-outChan:

			v, e := strconv.Atoi(output)
			if e != nil {
				panic(e)
			}

			r := rune(v)

			if r == 10 {
				//computer.Save("./game.json")
			}

			fmt.Fprintf(outTTY, "%c", r)

			if bot != nil {
				bot.OutputByte(byte(r))
				cmds := bot.GetCurrentCommands()
				if len(cmds) > 0 {
					for _, cmd := range cmds {
						inputQueue = append(inputQueue, []byte(cmd)...)
						inputQueue = append(inputQueue, byte(10))
					}
				}
			}

			computer.OutputProcessed()

		case <-quitChan:
			fmt.Println("game exited")
			break programLoop

		}
	}

	close(inChan)
	close(outChan)
	close(quitChan)

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
--- Day 25: Cryostasis ---
As you approach Santa's ship, your sensors report two important details:

First, that you might be too late: the internal temperature is -40 degrees.

Second, that one faint life signature is somewhere on the ship.

The airlock door is locked with a code; your best option is to send in a small droid to investigate the situation. You attach your ship to Santa's, break a small hole in the hull, and let the droid run in before you seal it up again. Before your ship starts freezing, you detach your ship and set it to automatically stay within range of Santa's ship.

This droid can follow basic instructions and report on its surroundings; you can communicate with it through an Intcode program (your puzzle input) running on an ASCII-capable computer.

As the droid moves through its environment, it will describe what it encounters. When it says Command?, you can give it a single instruction terminated with a newline (ASCII code 10). Possible instructions are:

Movement via north, south, east, or west.
To take an item the droid sees in the environment, use the command take <name of item>. For example, if the droid reports seeing a red ball, you can pick it up with take red ball.
To drop an item the droid is carrying, use the command drop <name of item>. For example, if the droid is carrying a green ball, you can drop it with drop green ball.
To get a list of all of the items the droid is currently carrying, use the command inv (for "inventory").
Extra spaces or other characters aren't allowed - instructions must be provided precisely.

Santa's ship is a Reindeer-class starship; these ships use pressure-sensitive floors to determine the identity of droids and crew members. The standard configuration for these starships is for all droids to weigh exactly the same amount to make them easier to detect. If you need to get past such a sensor, you might be able to reach the correct weight by carrying items from the environment.

Look around the ship and see if you can find the password for the main airlock.

Your puzzle answer was 529920.

--- Part Two ---
As you move through the main airlock, the air inside the ship is already heating up to reasonable levels. Santa explains that he didn't notice you coming because he was just taking a quick nap. The ship wasn't frozen; he just had the thermostat set to "North Pole".

You make your way over to the navigation console. It beeps. "Status: Stranded. Please supply measurements from 49 stars to recalibrate."

"49 stars? But the Elves told me you needed fifty--"

Santa just smiles and nods his head toward the window. There, in the distance, you can see the center of the Solar System: the Sun!

The navigation console beeps again.

If you like, you can .

Both parts of this puzzle are complete! They provide two gold stars: **
 */