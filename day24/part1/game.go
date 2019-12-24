package part1

import "fmt"

type Empty struct{}

var empty Empty

type State int32

func NewState(p *State) State {
	s := new(State)
	if p != nil {
		*s = *p
	}
	return *s
}

func (s State) Print( cols int ) {
	bits := cols * cols
	for i := 0; i < bits; i++ {
		if i%cols == 0 {
			fmt.Print("\n")
		}
		if s.IsSet(i) {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
}

func (s *State) Set(n int, b bool) {
	if b {
		*s = State(int(*s) | 1 << n)
	} else {
		mask := ^(1 << n)
		*s = State(int(*s) & mask)
	}
}

func (s State) IsSet(n int) bool {
	val := int(s) & (1 << n)
	return val > 0
}

func (s State) GetCounts(cols int, counts []int) {
	for i := range counts {

		counts[i] = 0

		if i >= cols {
			// check pos above if not the top row
			if s.IsSet(i - cols) {
				counts[i]++
			}
		}

		if i%cols != 0 {
			// check to the left if not the first column
			if s.IsSet(i - 1) {
				counts[i]++
			}
		}

		if (i+1)%cols != 0 {
			// check to the right if not the last column
			if s.IsSet(i + 1) {
				counts[i]++
			}
		}

		if i < len(counts)-cols {
			// check below if not the last row
			if s.IsSet(i + cols) {
				counts[i]++
			}
		}

	}
}

type Game struct {
	counts     [25]int
	minutes int
	states     map[State]Empty
	initial    State
	current State
}

func (g *Game) GetState() State {
	return g.current
}

func (g *Game) GetMinutes() int {
	return g.minutes
}

func NewGame(area []string) *Game {
	g := new(Game)
	g.states = make(map[State]Empty)
	g.initial = NewState(nil)

	n := 0
	for _, r := range area {
		for _, c := range r {
			if c == '#' {
				g.initial.Set(n, true)
			}
			n++
		}
	}

	g.current = g.initial

	g.states[g.initial] = empty
	return g
}

func (g *Game) Run(print bool) {
	cols := 5

	if print {
		fmt.Printf("\nInitial State:")
		g.current.Print(cols)
	}

	for {
		g.current.GetCounts(cols,g.counts[:])
		state := NewState(&(g.current))
		for n, c := range g.counts {
			if state.IsSet(n) && c != 1 {
				// a bug dies, unless exactly 1 bug is adjacent
				state.Set(n,false) // if already dead, no biggie
			} else if !state.IsSet(n) && (c == 1 || c == 2) {
				// empty spawns a bug if exactly 1 or 2 bugs are adjacent
				state.Set(n,true)
			}
		}
		g.current = state
		g.minutes++

		if print {
			fmt.Printf("\nAfter %d minutes:",g.minutes)
			g.current.Print(cols)
		}

		if _, repeated := g.states[state]; repeated {
			fmt.Printf("seeing a repeated layout, with bio diversity of %d after %d minutes\n",state,g.minutes)
			break
		}
		g.states[state] = empty

	}
}
