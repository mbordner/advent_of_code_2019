package main

import (
	"github.com/mbordner/advent_of_code_2019/day20/part1"
	"errors"
	"fmt"
	"os"
	"strings"
)

func doPart1() {
	maze := getMaze("input.txt")
	game := part1.NewGame(maze)
	path, distance := game.ShortestPath("AA", "ZZ")
	if path == nil {
		panic(errors.New("invalid path"))
	}
	fmt.Println("distance from AA to ZZ for part1: ", distance)
}

func main() {
	doPart1()
}

func getMaze(filename string) [][]byte {
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

	rows := strings.Split(string(buffer), "\n")

	chars := make([][]byte, len(rows), len(rows))
	for y := range rows {
		chars[y] = []byte(rows[y])
	}

	return chars
}
