package main

// https://adventofcode.com/2019/day/1
func puzzle(input int) int {
	return input/3 - 2
}

// https://adventofcode.com/2019/day/2
func puzzle2(input int) int {
	fuel := input/3 - 2
	more := fuel
	for more/3 > 2 {
		more = more/3 - 2
		fuel += more
	}
	return fuel
}