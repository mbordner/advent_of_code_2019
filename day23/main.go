package main

import (
	"github.com/mbordner/advent_of_code_2019/day23/geom"
	"github.com/mbordner/advent_of_code_2019/day23/intcode"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Packet struct {
	To   int
	From int
	geom.Pos
}

func (p Packet) String() string {
	return fmt.Sprintf("{to:%d, from:%d, data:%s }", p.To, p.From, p.Pos)
}

func NewPacket(to, from, x, y int) *Packet {
	p := Packet{
		To:   to,
		From: from,
		Pos:  geom.Pos{X: x, Y: y},
	}
	return &p
}

type Computer struct {
	ID              int
	program         []string
	intCodeComputer *intcode.IntCodeComputer
	inputQueue      []int
	outputQueue     []int
	networkOut      chan *Packet
	in              chan string
	out             chan string
	quit            chan string
	wg              *sync.WaitGroup
	mux             sync.Mutex
	idSent          bool
}

func NewComputer(id int) *Computer {
	c := new(Computer)

	c.ID = id

	c.in = make(chan string, 1)
	c.out = make(chan string, 1)
	c.quit = make(chan string, 1)

	c.networkOut = make(chan *Packet, 25)

	c.inputQueue = make([]int, 0, 25)
	c.outputQueue = make([]int, 0, 3)

	return c
}

func (c *Computer) Execute(p []string, wg *sync.WaitGroup) {
	c.wg = wg
	c.program = make([]string, len(p), len(p))
	copy(c.program, p)
	prompt := "ready"
	c.intCodeComputer = intcode.NewIntCodeComputer(c.program, c.in, c.out, c.quit, true, &prompt)
	go c.loop()
	go c.intCodeComputer.Execute()
}

func (c *Computer) Receive(p *Packet) {
	c.mux.Lock()
	c.inputQueue = append(c.inputQueue, p.X)
	c.inputQueue = append(c.inputQueue, p.Y)
	c.mux.Unlock()
}

func (c *Computer) GetNetworkOutputChannel() <-chan *Packet {
	return c.networkOut
}

func (c *Computer) loop() {
programLoop:
	for {
		select {
		case prompt := <-c.intCodeComputer.GetPromptChannel():
			if prompt == "ready" {
				c.mux.Lock()

				if c.idSent {

					if len(c.inputQueue) > 0 {

						p := c.inputQueue[0]
						c.inputQueue = c.inputQueue[1:]
						c.mux.Unlock()
						c.in <- fmt.Sprintf("%d", p)

					} else {
						c.mux.Unlock()
						c.in <- "-1" // no packets available
					}
				} else {
					// first input must be the network id
					c.idSent = true
					c.in <- fmt.Sprintf("%d", c.ID)
					c.mux.Unlock()

				}
			}
		case p := <-c.out:

			c.mux.Lock()

			i, e := strconv.Atoi(p)
			if e != nil {
				panic(e)
			}

			c.outputQueue = append(c.outputQueue,i)

			if len(c.outputQueue) == 3 {
				to := c.outputQueue[0]
				x := c.outputQueue[1]
				y := c.outputQueue[2]
				c.outputQueue = make([]int, 0, 3)
				packet := NewPacket(to, c.ID, x, y)
				c.networkOut <- packet
			}

			c.mux.Unlock()

			c.intCodeComputer.OutputProcessed()


		case <-c.quit:
			fmt.Printf("computer %d shutting down\n", c.ID)
			break programLoop

		}
	}

	close(c.in)
	close(c.out)
	close(c.quit)

	close(c.networkOut)
}

func main() {
	program := getProgram("program.txt")

	var wg sync.WaitGroup

	const NUM_COMPUTERS = 50

	computers := make([]*Computer, NUM_COMPUTERS, NUM_COMPUTERS)
	netout := make([]<-chan *Packet, NUM_COMPUTERS, NUM_COMPUTERS)

	for i := 0; i < NUM_COMPUTERS; i++ {
		computers[i] = NewComputer(i)
		wg.Add(1)
	}

	for i := range computers {
		netout[i] = computers[i].GetNetworkOutputChannel()
		computers[i].Execute(program, &wg)
	}

	for _, no := range netout {
		go func(no <-chan *Packet) {
			for p := range no {
				fmt.Printf("incoming packet: %s\n", p)
				if p.To >= 0 && p.To < len(computers) {
					computers[p.To].Receive(p)
				}
			}
		}(no)
	}

	wg.Wait()

	fmt.Println("all computers shutdown.")
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
