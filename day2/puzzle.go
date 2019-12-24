package main

import "fmt"

type empty struct{}

var (
	opCodes = map[int]empty{
		1:  empty{},
		2:  empty{},
		99: empty{},
	}
)

// https://adventofcode.com/2019/day/2
func puzzle(program []int) []int {
	i := 0
	for i < len(program) {
		if _, ok := opCodes[program[i]]; ok {
			opCode := program[i]

			if opCode == 99 {
				break
			}

			a := program[program[i+1]]
			b := program[program[i+2]]

			var value int
			if opCode == 1 {
				value = a + b
			} else if opCode == 2 {
				value = a * b
			}

			program[program[i+3]] = value

			i += 4

		} else {
			panic(fmt.Errorf("invalid opcode %d at pos %d", program[i], i))
		}
	}
	return program
}
