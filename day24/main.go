package main

import (
	"errors"
	"github.com/mbordner/advent_of_code_2019/day24/part1"
	"github.com/mbordner/advent_of_code_2019/day24/part2"
	"os"
	"strings"
)

func doPart1() {
	g := part1.NewGame(getGame("game.txt"))
	g.Run(false)
}

func doPart2() {
	g := part2.NewGame(getGame("game.txt"))
	g.Run(200, false)
}

func main() {
	doPart2()
}

func getGame(filename string) []string {
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

	return strings.Split(string(buffer), "\n")
}
