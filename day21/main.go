package main

import (
	"github.com/mbordner/advent_of_code_2019/day21/intcode"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := make(chan string, 1)
	out := make(chan string, 1)
	quit := make(chan string, 1)

	program := getProgram("program.txt")

	intCodeComputer := intcode.NewIntCodeComputer(program, in, out, quit, true)

	go intCodeComputer.Execute()

	// part 1 commands
	/*
	commands := strings.Split(`NOT C T
AND D T
OR T J
NOT A T
OR T J
WALK`,"\n")
	 */

	//  (D && !(A && B && (C || !H)))

	commands := strings.Split(`NOT C T
AND D T
OR T J
AND I J
NOT A T
OR T J
RUN`,"\n")

	l := 0
	for i := range commands {
		l += len(commands[i])
	}
	bytes := make([]byte,0,l)
	for _, c := range commands {
		for _, r := range c {
			bytes = append(bytes,byte(r))
		}
		bytes = append(bytes,byte(10))
	}

programLoop:
	for {
		select {
		case <-in:

			b := bytes[0]
			bytes = bytes[1:]
			str := fmt.Sprintf("%d",b)
			in <- str

		case response := <-out:

			b, e := strconv.Atoi(response)
			if e != nil {
				panic(e)
			}
			if b > 127 {
				fmt.Println(response)
			} else {
				fmt.Print(string(byte(b)))
			}

			intCodeComputer.OutputProcessed()

		case <-quit:
			fmt.Println("\nprogram exited")
			break programLoop

		}
	}
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