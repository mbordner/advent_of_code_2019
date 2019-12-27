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
