package intcode

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"
)

type IntCodeComputer struct {
	opCodes       map[int]int
	origProgram   []string
	program       []string
	prompt        chan string
	in            chan string
	out           chan string
	quit          chan<- string
	outputWG      sync.WaitGroup
	pauseOnOutput bool
	inputPrompt   *string
}

func NewIntCodeComputer(program []string, in chan string, out chan string, quit chan<- string, pauseOnOutput bool, inputPrompt *string) *IntCodeComputer {
	c := new(IntCodeComputer)
	c.origProgram = make([]string, len(program), len(program))
	c.program = program
	copy(c.origProgram, c.program)
	c.prompt = make(chan string, 1)
	c.in = in
	c.out = out
	c.quit = quit
	c.pauseOnOutput = pauseOnOutput
	c.inputPrompt = inputPrompt
	// op code to number of bytes including op code and parameters
	c.opCodes = map[int]int{
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
	return c
}

func (c *IntCodeComputer) GetPromptChannel() <-chan string {
	return c.prompt
}

func (c *IntCodeComputer) OutputProcessed() {
	c.outputWG.Done()
}

func (c *IntCodeComputer) getProgramValue(pos int) string {
	if pos < len(c.program) {
		return c.program[pos]
	}
	return "0"
}

func (c *IntCodeComputer) setProgramValue(pos int, val string) string {
	if pos < len(c.program) {
		c.program[pos] = val
		return val
	}
	newProgram := make([]string, pos+1, pos+1)
	for i := range newProgram {
		newProgram[i] = "0"
	}
	copy(newProgram, c.program)
	newProgram[pos] = val
	c.program = newProgram
	return val
}

func (c *IntCodeComputer) GetProgram() []string {
	return c.program
}

func (c *IntCodeComputer) Reset() {
	c.program = make([]string, len(c.origProgram), len(c.origProgram))
	copy(c.program, c.origProgram)
}

func (c *IntCodeComputer) Execute() {

	var lastOut string

	var relativeBase int

	i := 0
instructions:
	for i < len(c.program) {
		opCode, err := strconv.Atoi(c.program[i])
		if err != nil {
			panic(err)
		}
		if length, ok := c.opCodes[opCode%100]; ok {
			tmp := opCode
			opCode = tmp % 100
			tmp /= 100

			if opCode == 99 {
				c.quit <- lastOut
				break instructions
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
					pos, err := strconv.Atoi(c.program[i+j])
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
					pos, err := strconv.Atoi(c.program[i+j])
					if err != nil {
						panic(err)
					}
					paramPositions[j-1] = pos + relativeBase
				}

				tmp /= 10
			}

			switch opCode {
			case 1:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[1])))
				}
				val3 := new(big.Int)
				val3.Add(val1, val2)
				c.setProgramValue(paramPositions[2], val3.Text(10))
			case 2:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[1])))
				}
				val3 := new(big.Int)
				val3.Mul(val1, val2)
				c.setProgramValue(paramPositions[2], val3.Text(10))
			case 3:
				if c.inputPrompt != nil {
					c.prompt <- *(c.inputPrompt)
				}
				value := <-c.in // wait for input to be received
				c.setProgramValue(paramPositions[0], value)
			case 4:
				lastOut = c.getProgramValue(paramPositions[0])
				c.outputWG.Add(1)
				c.out <- lastOut
				if !c.pauseOnOutput {
					c.OutputProcessed()
				}
				c.outputWG.Wait()
			case 5:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				if val1.Int64() != int64(0) {
					pos, err := strconv.Atoi(c.getProgramValue(paramPositions[1]))
					if err != nil {
						panic(err)
					}
					i = pos
					continue instructions
				}
			case 6:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				if val1.Int64() == 0 {
					pos, err := strconv.Atoi(c.getProgramValue(paramPositions[1]))
					if err != nil {
						panic(err)
					}
					i = pos
					continue instructions
				}
			case 7:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[1])))
				}
				if val1.Cmp(val2) < 0 {
					c.setProgramValue(paramPositions[2], "1")
				} else {
					c.setProgramValue(paramPositions[2], "0")
				}
			case 8:
				val1, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[0]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[0])))
				}
				val2, ok := new(big.Int).SetString(c.getProgramValue(paramPositions[1]), 10)
				if !ok {
					panic(fmt.Errorf("invalid int %s", c.getProgramValue(paramPositions[1])))
				}
				if val1.Cmp(val2) == 0 {
					c.setProgramValue(paramPositions[2], "1")
				} else {
					c.setProgramValue(paramPositions[2], "0")
				}
			case 9:
				val, err := strconv.Atoi(c.getProgramValue(paramPositions[0]))
				if err != nil {
					panic(err)
				}
				relativeBase += val
			}

			i += length

		} else {
			panic(fmt.Errorf("invalid opcode %s at pos %d", c.program[i], i))
		}
	}
}
