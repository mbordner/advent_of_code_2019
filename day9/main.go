package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	// op code to number of bytes including op code and parameters
	opCodes = map[int]int{
		1:  4, // add first 2 params, set 3 param location to value
		2:  4, // multiply first 2 params, set 3 param location to value
		3:  2, // read string
		4:  2, // output string
		5:  3, // if param 0 is not 0, jump to param 1 location
		6:  3, // if param 0 is 0, jump to param 1 location
		7:  4, // if param 0 is less than param 1, set location at param 3 to 1, else set it to 0
		8:  4, // if param 0 is equal to param 1, set location at param 3 to 1, else set it to 0
		9:  2, // set relative base
		99: 0, // halt
	}
)

func getProgramValue(program []string, pos int) string {
	if pos < len(program) {
		return program[pos]
	}
	return "0"
}

func setProgramValue(program []string, pos int, val string) []string {
	if pos < len(program) {
		program[pos] = val
		return program
	}
	newProgram := make([]string, pos+1, pos+1)
	copy(newProgram, program)
	newProgram[pos] = val
	return newProgram
}

func execute(program []string, in chan string, out chan<- string, quit chan<- string) []string {

	var lastOut string

	var relativeBase int

	i := 0
instructions:
	for i < len(program) {
		opCode, err := strconv.Atoi(program[i])
		if err != nil {
			panic(err)
		}
		if length, ok := opCodes[opCode%100]; ok {
			tmp := opCode
			opCode = tmp % 100
			tmp /= 100

			if opCode == 99 {
				quit <- lastOut
				break
			}

			l := length - 1

			// parameter address modes
			modes := make([]int, l, l)
			// op code paramPositions
			paramPositions := make([]int, l, l)

			for j := 1; j < length; j++ {
				modes[j-1] = tmp % 10

				switch modes[j-1] {
				case 0:
					// position mode
					pos, err := strconv.Atoi(program[i+j])
					if err != nil {
						panic(err)
					}
					paramPositions[j-1] = pos
				default:
					fallthrough
				case 1:
					// immediate mode
					paramPositions[j-1] = i + j
				case 2:
					// relative mode
					pos, err := strconv.Atoi(program[i+j])
					if err != nil {
						panic(err)
					}
					paramPositions[j-1] = pos + relativeBase
				}

				tmp /= 10
			}

			switch opCode {
			case 1:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[1])))
				}
				val3 := new(big.Int)
				val3.Add(val1, val2)
				program = setProgramValue(program, paramPositions[2], val3.Text(10))
			case 2:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[1])))
				}
				val3 := new(big.Int)
				val3.Mul(val1, val2)
				program = setProgramValue(program, paramPositions[2], val3.Text(10))
			case 3:
				in <- "input:"
				value := <-in
				program = setProgramValue(program, paramPositions[0], value)
			case 4:
				lastOut = getProgramValue(program,paramPositions[0])
				out <- lastOut
			case 5:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				if val1.Int64() != int64(0) {
					pos, err := strconv.Atoi(getProgramValue(program,paramPositions[1]))
					if err != nil {
						panic(err)
					}
					i = pos
					continue instructions
				}
			case 6:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				if val1.Int64() == 0 {
					pos, err := strconv.Atoi(getProgramValue(program,paramPositions[1]))
					if err != nil {
						panic(err)
					}
					i = pos
					continue instructions
				}
			case 7:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[1])))
				}
				if val1.Cmp(val2) < 0 {
					program = setProgramValue(program, paramPositions[2], "1")
				} else {
					program = setProgramValue(program, paramPositions[2], "0")
				}
			case 8:
				val1, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(getProgramValue(program,paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", getProgramValue(program,paramPositions[1])))
				}
				if val1.Cmp(val2) == 0 {
					program = setProgramValue(program, paramPositions[2], "1")
				} else {
					program = setProgramValue(program, paramPositions[2], "0")
				}
			case 9:
				val, err := strconv.Atoi(getProgramValue(program,paramPositions[0]))
				if err != nil {
					panic(err)
				}
				relativeBase += val
			}

			i += length

		} else {
			panic(fmt.Errorf("invalid opcode %d at pos %d", program[i], i))
		}
	}
	return program
}

func main() {

	in := make(chan string, 1)
	out := make(chan string, 1)
	quit := make(chan string, 1)

	go execute(getPart1Program(), in, out, quit)

programLoop:
	for {
		select {
		case prompt := <-in:
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(prompt)
			text, _ := reader.ReadString('\n')
			in <- text[0 : len(text)-1]
		case val := <-out:
			fmt.Println(val)
		case lastOut := <-quit:
			fmt.Println("program exited with last output: ", lastOut)
			break programLoop
		}
	}

}

func getProgramTest1() []string {
	program := `109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99`
	return strings.Split(program, ",")
}

func getProgramTest2() []string {
	program := `1102,34915192,34915192,7,4,7,99,0`
	return strings.Split(program, ",")
}

func getProgramTest3() []string {
	program := `104,1125899906842624,99`
	return strings.Split(program, ",")
}

func getPart1Program() []string {
	program := `1102,34463338,34463338,63,1007,63,34463338,63,1005,63,53,1102,1,3,1000,109,988,209,12,9,1000,209,6,209,3,203,0,1008,1000,1,63,1005,63,65,1008,1000,2,63,1005,63,904,1008,1000,0,63,1005,63,58,4,25,104,0,99,4,0,104,0,99,4,17,104,0,99,0,0,1101,234,0,1027,1101,0,568,1023,1102,844,1,1025,1101,0,23,1008,1102,1,1,1021,1102,27,1,1011,1101,0,26,1004,1102,1,586,1029,1102,29,1,1014,1101,0,22,1015,1102,36,1,1016,1101,35,0,1013,1102,20,1,1003,1102,1,37,1019,1101,30,0,1006,1102,34,1,1000,1101,571,0,1022,1102,1,28,1005,1101,39,0,1009,1102,38,1,1017,1102,591,1,1028,1102,1,31,1007,1102,24,1,1010,1101,0,33,1001,1101,0,21,1018,1101,0,0,1020,1101,25,0,1002,1102,32,1,1012,1101,0,237,1026,1101,0,853,1024,109,29,1206,-9,195,4,187,1106,0,199,1001,64,1,64,1002,64,2,64,109,-26,2102,1,0,63,1008,63,23,63,1005,63,223,1001,64,1,64,1105,1,225,4,205,1002,64,2,64,109,16,2106,0,8,1106,0,243,4,231,1001,64,1,64,1002,64,2,64,109,-19,21101,40,0,10,1008,1010,40,63,1005,63,265,4,249,1106,0,269,1001,64,1,64,1002,64,2,64,109,-2,2107,31,8,63,1005,63,289,1001,64,1,64,1105,1,291,4,275,1002,64,2,64,109,2,1208,7,28,63,1005,63,307,1106,0,313,4,297,1001,64,1,64,1002,64,2,64,109,-1,1207,9,24,63,1005,63,335,4,319,1001,64,1,64,1105,1,335,1002,64,2,64,109,5,1201,0,0,63,1008,63,25,63,1005,63,355,1105,1,361,4,341,1001,64,1,64,1002,64,2,64,109,-13,1202,9,1,63,1008,63,34,63,1005,63,383,4,367,1105,1,387,1001,64,1,64,1002,64,2,64,109,32,1205,-3,403,1001,64,1,64,1106,0,405,4,393,1002,64,2,64,109,-14,2108,31,-2,63,1005,63,423,4,411,1105,1,427,1001,64,1,64,1002,64,2,64,109,11,1206,1,439,1105,1,445,4,433,1001,64,1,64,1002,64,2,64,109,-21,1208,4,20,63,1005,63,467,4,451,1001,64,1,64,1105,1,467,1002,64,2,64,109,6,1207,-5,33,63,1005,63,487,1001,64,1,64,1106,0,489,4,473,1002,64,2,64,109,-12,1202,8,1,63,1008,63,34,63,1005,63,509,1106,0,515,4,495,1001,64,1,64,1002,64,2,64,109,28,1205,0,529,4,521,1106,0,533,1001,64,1,64,1002,64,2,64,109,3,21101,41,0,-9,1008,1015,38,63,1005,63,557,1001,64,1,64,1106,0,559,4,539,1002,64,2,64,109,-11,2105,1,10,1105,1,577,4,565,1001,64,1,64,1002,64,2,64,109,23,2106,0,-8,4,583,1105,1,595,1001,64,1,64,1002,64,2,64,109,-15,21108,42,42,-6,1005,1015,613,4,601,1106,0,617,1001,64,1,64,1002,64,2,64,109,-14,21107,43,44,8,1005,1015,639,4,623,1001,64,1,64,1106,0,639,1002,64,2,64,109,11,2107,38,-9,63,1005,63,661,4,645,1001,64,1,64,1106,0,661,1002,64,2,64,109,-2,21107,44,43,3,1005,1019,677,1105,1,683,4,667,1001,64,1,64,1002,64,2,64,109,-7,21108,45,42,1,1005,1010,703,1001,64,1,64,1106,0,705,4,689,1002,64,2,64,109,-5,2102,1,1,63,1008,63,28,63,1005,63,727,4,711,1106,0,731,1001,64,1,64,1002,64,2,64,109,13,21102,46,1,0,1008,1017,46,63,1005,63,753,4,737,1106,0,757,1001,64,1,64,1002,64,2,64,109,-4,2101,0,-5,63,1008,63,20,63,1005,63,781,1001,64,1,64,1105,1,783,4,763,1002,64,2,64,109,1,21102,47,1,0,1008,1014,48,63,1005,63,803,1105,1,809,4,789,1001,64,1,64,1002,64,2,64,109,-3,2101,0,-4,63,1008,63,31,63,1005,63,835,4,815,1001,64,1,64,1105,1,835,1002,64,2,64,109,6,2105,1,7,4,841,1001,64,1,64,1105,1,853,1002,64,2,64,109,-21,2108,33,10,63,1005,63,873,1001,64,1,64,1105,1,875,4,859,1002,64,2,64,109,6,1201,4,0,63,1008,63,30,63,1005,63,901,4,881,1001,64,1,64,1105,1,901,4,64,99,21102,27,1,1,21102,1,915,0,1106,0,922,21201,1,64720,1,204,1,99,109,3,1207,-2,3,63,1005,63,964,21201,-2,-1,1,21102,1,942,0,1105,1,922,21202,1,1,-1,21201,-2,-3,1,21101,957,0,0,1105,1,922,22201,1,-1,-2,1105,1,968,21202,-2,1,-2,109,-3,2106,0,0`
	return strings.Split(program, ",")
}

/**
--- Day 9: Sensor Boost ---
You've just said goodbye to the rebooted rover and left Mars when you receive a faint distress signal coming from the asteroid belt. It must be the Ceres monitoring station!

In order to lock on to the signal, you'll need to boost your sensors. The Elves send up the latest BOOST program - Basic Operation Of System Test.

While BOOST (your puzzle input) is capable of boosting your sensors, for tenuous safety reasons, it refuses to do so until the computer it runs on passes some checks to demonstrate it is a complete Intcode computer.

Your existing Intcode computer is missing one key feature: it needs support for parameters in relative mode.

Parameters in mode 2, relative mode, behave very similarly to parameters in position mode: the parameter is interpreted as a position. Like position mode, parameters in relative mode can be read from or written to.

The important difference is that relative mode parameters don't count from address 0. Instead, they count from a value called the relative base. The relative base starts at 0.

The address a relative mode parameter refers to is itself plus the current relative base. When the relative base is 0, relative mode parameters and position mode parameters with the same value refer to the same address.

For example, given a relative base of 50, a relative mode parameter of -7 refers to memory address 50 + -7 = 43.

The relative base is modified with the relative base offset instruction:

Opcode 9 adjusts the relative base by the value of its only parameter. The relative base increases (or decreases, if the value is negative) by the value of the parameter.
For example, if the relative base is 2000, then after the instruction 109,19, the relative base would be 2019. If the next instruction were 204,-34, then the value at address 1985 would be output.

Your Intcode computer will also need a few other capabilities:

The computer's available memory should be much larger than the initial program. Memory beyond the initial program starts with the value 0 and can be read or written like any other memory. (It is invalid to try to access memory at a negative address, though.)
The computer should have support for large numbers. Some instructions near the beginning of the BOOST program will verify this capability.
Here are some example programs that use these features:

109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99 takes no input and produces a copy of itself as output.
1102,34915192,34915192,7,4,7,99,0 should output a 16-digit number.
104,1125899906842624,99 should output the large number in the middle.
The BOOST program will ask for a single input; run it in test mode by providing it the value 1. It will perform a series of checks on each opcode, output any opcodes (and the associated parameter modes) that seem to be functioning incorrectly, and finally output a BOOST keycode.

Once your Intcode computer is fully functional, the BOOST program should report no malfunctioning opcodes when run in test mode; it should only output a single value, the BOOST keycode. What BOOST keycode does it produce?

Your puzzle answer was 3638931938.

--- Part Two ---
You now have a complete Intcode computer.

Finally, you can lock on to the Ceres distress signal! You just need to boost your sensors using the BOOST program.

The program runs in sensor boost mode by providing the input instruction the value 2. Once run, it will boost the sensors automatically, but it might take a few seconds to complete the operation on slower hardware. In sensor boost mode, the program will output a single value: the coordinates of the distress signal.

Run the BOOST program in sensor boost mode. What are the coordinates of the distress signal?

Your puzzle answer was 86025.

Both parts of this puzzle are complete! They provide two gold stars: **
 */