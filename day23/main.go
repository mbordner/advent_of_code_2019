package main

import (
	"errors"
	"fmt"
	"github.com/mbordner/advent_of_code_2019/day23/geom"
	"github.com/mbordner/advent_of_code_2019/day23/intcode"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

type Nat struct {
	lastPacket *Packet
	quit       chan bool
	computers []*Computer
}

func NewNat(computers []*Computer) *Nat {
	n := new(Nat)
	n.computers = computers
	go n.loop()
	return n
}

func (n *Nat) Halt() {
	n.quit <- true
}

func (n *Nat) loop() {
	tick := time.Tick(time.Duration(1) * time.Second)

	const IDLE_THRESHOLD = 5

programLoop:
	for {
		select {
		case <-tick:
			if n.lastPacket != nil {
				idle := true
				for i:= range n.computers {
					if n.computers[i].GetIdleCount() < IDLE_THRESHOLD {
						idle = false
						break
					}
				}
				if idle {
					fmt.Printf(">>>>> sending idle packet %s\n",n.lastPacket)
					n.computers[0].Receive(n.lastPacket)
				}
			}
		case <-n.quit:
			break programLoop
		}
	}
}

func (n *Nat) Receive(p *Packet) {
	n.lastPacket = p
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
	idleCount       int
	running bool
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

func (c *Computer) GetIdleCount() int {
	return c.idleCount
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
	if c.running {
		c.mux.Lock()
		c.inputQueue = append(c.inputQueue, p.X)
		c.inputQueue = append(c.inputQueue, p.Y)
		c.mux.Unlock()
	}
}

func (c *Computer) GetNetworkOutputChannel() <-chan *Packet {
	return c.networkOut
}

func (c *Computer) loop() {
	c.running = true

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
						c.idleCount = 0
						c.mux.Unlock()
						c.in <- fmt.Sprintf("%d", p)

					} else {
						c.idleCount++
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

			c.outputQueue = append(c.outputQueue, i)

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

	c.running = false

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

	nat := NewNat(computers)

	for _, no := range netout {
		go func(no <-chan *Packet) {
			for p := range no {
				//fmt.Printf("incoming packet: %s\n", p)
				if p.To >= 0 && p.To < len(computers) {
					computers[p.To].Receive(p)
				} else if p.To == 255 {
					nat.Receive(p)
				}
			}
		}(no)
	}

	nat.Halt()

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

/**
--- Day 23: Category Six ---
The droids have finished repairing as much of the ship as they can. Their report indicates that this was a Category 6 disaster - not because it was that bad, but because it destroyed the stockpile of Category 6 network cables as well as most of the ship's network infrastructure.

You'll need to rebuild the network from scratch.

The computers on the network are standard Intcode computers that communicate by sending packets to each other. There are 50 of them in total, each running a copy of the same Network Interface Controller (NIC) software (your puzzle input). The computers have network addresses 0 through 49; when each computer boots up, it will request its network address via a single input instruction. Be sure to give each computer a unique network address.

Once a computer has received its network address, it will begin doing work and communicating over the network by sending and receiving packets. All packets contain two values named X and Y. Packets sent to a computer are queued by the recipient and read in the order they are received.

To send a packet to another computer, the NIC will use three output instructions that provide the destination address of the packet followed by its X and Y values. For example, three output instructions that provide the values 10, 20, 30 would send a packet with X=20 and Y=30 to the computer with address 10.

To receive a packet from another computer, the NIC will use an input instruction. If the incoming packet queue is empty, provide -1. Otherwise, provide the X value of the next packet; the computer will then use a second input instruction to receive the Y value for the same packet. Once both values of the packet are read in this way, the packet is removed from the queue.

Note that these input and output instructions never block. Specifically, output instructions do not wait for the sent packet to be received - the computer might send multiple packets before receiving any. Similarly, input instructions do not wait for a packet to arrive - if no packet is waiting, input instructions should receive -1.

Boot up all 50 computers and attach them to your network. What is the Y value of the first packet sent to address 255?

Your puzzle answer was 21089.

--- Part Two ---
Packets sent to address 255 are handled by a device called a NAT (Not Always Transmitting). The NAT is responsible for managing power consumption of the network by blocking certain packets and watching for idle periods in the computers.

If a packet would be sent to address 255, the NAT receives it instead. The NAT remembers only the last packet it receives; that is, the data in each packet it receives overwrites the NAT's packet memory with the new packet's X and Y values.

The NAT also monitors all computers on the network. If all computers have empty incoming packet queues and are continuously trying to receive packets without sending packets, the network is considered idle.

Once the network is idle, the NAT sends only the last packet it received to address 0; this will cause the computers on the network to resume activity. In this way, the NAT can throttle power consumption of the network when the ship needs power in other areas.

Monitor packets released to the computer at address 0 by the NAT. What is the first Y value delivered by the NAT to the computer at address 0 twice in a row?

Your puzzle answer was 16658.

Both parts of this puzzle are complete! They provide two gold stars: **
 */
