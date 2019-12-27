package intcode

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
)

type Memory struct {
	Ptr int `json:"ptr"`
	Program []string `json:"program"`
	RelativeBase int `json:"relativeBase"`
	LastOut string `json:"lastOut"`
}

type IntCodeComputer struct {
	opCodes       map[int]int
	origProgram   []string
	program       []string
	ptr           int
	relativeBase  int
	lastOut       string
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

func (c *IntCodeComputer) Save(filename string) {
	mem := Memory{
		Ptr: c.ptr,
		Program: c.program,
		RelativeBase: c.relativeBase,
		LastOut: c.lastOut,
	}

	file, _ := os.OpenFile(filename, os.O_RDWR | os.O_TRUNC, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)
	err := encoder.Encode(mem)
	if err != nil {
		panic(err)
	}
}

func (c *IntCodeComputer) Load(filename string) {
	file, _ := os.Open(filename)
	defer file.Close()

	decoder := json.NewDecoder(file)

	mem := Memory{}

	decoder.Decode(&mem)

	c.ptr = mem.Ptr
	c.program = mem.Program
	c.relativeBase = mem.RelativeBase
	c.lastOut = mem.LastOut

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


instructions:
	for c.ptr < len(c.program) {
		opCode, err := strconv.Atoi(c.program[c.ptr])
		if err != nil {
			panic(err)
		}
		if length, ok := c.opCodes[opCode%100]; ok {
			tmp := opCode
			opCode = tmp % 100
			tmp /= 100

			if opCode == 99 {
				c.quit <- c.lastOut
				close(c.prompt)
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
					pos, err := strconv.Atoi(c.program[c.ptr+j])
					if err != nil {
						panic(err)
					}
					paramPositions[j-1] = pos
				default:
					fallthrough
				case 1:
					// immediate mode
					paramPositions[j-1] = c.ptr + j
				case 2:
					// relative mode
					pos, err := strconv.Atoi(c.program[c.ptr+j])
					if err != nil {
						panic(err)
					}
					paramPositions[j-1] = pos + c.relativeBase
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
				c.lastOut = c.getProgramValue(paramPositions[0])
				c.outputWG.Add(1)
				c.out <- c.lastOut
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
					c.ptr = pos
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
					c.ptr = pos
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
				c.relativeBase += val
			}

			c.ptr += length

		} else {
			panic(fmt.Errorf("invalid opcode %s at pos %d", c.program[c.ptr], c.ptr))
		}
	}
}
