package part2

import (
	"fmt"
	"math"
)

type Side int

const (
	Right Side = iota
	Bottom
	Left
	Top
)

var (
	states = NewStates()
)

type States struct {
	min    int
	max    int
	states map[int]*State
}

func NewStates() *States {
	s := new(States)
	s.states = make(map[int]*State)
	return s
}

func (ss *States) GetMin() int {
	return ss.min
}

func (ss *States) GetMax() int {
	return ss.max
}

func (ss *States) Add(id int, s *State) {
	ss.states[id] = s
	if id < ss.min {
		ss.min = id
	}
	if id > ss.max {
		ss.max = id
	}
	if prev, ok := ss.states[id-1]; ok {
		prev.next = s
		s.prev = prev
	}
	if next, ok := ss.states[id+1]; ok {
		next.prev = s
		s.next = next
	}
}

func (ss *States) GetState(id int) *State {
	if s, ok := ss.states[id]; ok {
		return s
	}
	return nil
}

func (ss *States) GetStates() []*State {
	l := ss.max + 1 - ss.min
	arr := make([]*State, l, l)
	for i, j := ss.min, 0; j < l; i, j = i+1, j+1 {
		arr[j] = ss.states[i]
	}
	return arr
}

type State struct {
	val  int32
	prev *State
	next *State
	id   int
}

func NewState(id int, prev *State, next *State, val int32) *State {
	s := new(State)
	s.id = id
	s.prev = prev
	s.next = next
	s.val = val
	return s
}

func (s State) Clone() *State {
	t := new(State)
	t.id = s.id
	t.next = s.next
	t.prev = s.prev
	t.val = s.val
	return t
}

func (s State) GetCount() int {
	n := int(s.val)
	sum := 0
	for n > 0 {
		sum += n & 1
		n >>= 1
	}
	return sum
}

func (s State) Print(cols int) {
	bits := cols * cols
	fmt.Print("\ndepth: ", s.id)
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
		s.val |= int32(1 << n)
	} else {
		mask := ^(1 << n)
		s.val &= int32(mask)
	}
}

func (s State) IsSet(n int) bool {
	val := s.val & (1 << n)
	return val > 0
}

func (s State) GetSideCount(side Side, cols int, counts []int) int {
	count := 0
	switch side {
	case Top:
		for i := 0; i < cols; i++ {
			if s.IsSet(i) {
				count++
			}
		}
	case Bottom:
		for i := len(counts) - cols; i < len(counts); i++ {
			if s.IsSet(i) {
				count++
			}
		}
	case Left:
		for i := 0; i < len(counts); i += cols {
			if s.IsSet(i) {
				count++
			}
		}
	case Right:
		for i := len(counts) - 1; i > 0; i -= cols {
			if s.IsSet(i) {
				count++
			}
		}
	}

	return count
}

func (s State) GetCounts(cols int, counts []int) {
	center := len(counts) / 2
	for i := range counts {

		if i == center {
			continue // counts for center will come from the next level
		}

		counts[i] = 0

		if i >= cols {
			// check pos above if not the top row
			if i-cols == center {
				// if below center
				if s.next != nil {
					// if has next level, get bottom count of next level
					counts[i] += s.next.GetSideCount(Bottom, cols, counts)
				}

			} else {
				if s.IsSet(i - cols) {
					counts[i]++
				}
			}
		} else if s.prev != nil {
			// if top row, and has previous level, add count of previous level
			// square above center
			if s.prev.IsSet(center - cols) {
				counts[i]++
			}
		}

		if i%cols != 0 {
			// check to the left if not the first column
			if i-1 == center {
				// if right of center
				if s.next != nil {
					// if has next level, get right side count of next level
					counts[i] += s.next.GetSideCount(Right, cols, counts)
				}
			} else {
				if s.IsSet(i - 1) {
					counts[i]++
				}
			}
		} else if s.prev != nil {
			// if first column and has previous level, add count of previous level
			// square to left of center
			if s.prev.IsSet(center - 1) {
				counts[i]++
			}
		}

		if (i+1)%cols != 0 {
			// check to the right if not the last column
			if i+1 == center {
				// if left of center
				if s.next != nil {
					// if has next level, get left side count of next level
					counts[i] += s.next.GetSideCount(Left, cols, counts)
				}

			} else {
				if s.IsSet(i + 1) {
					counts[i]++
				}
			}
		} else if s.prev != nil {
			// if last column and has previous level, add count of previous level
			// square to right of center
			if s.prev.IsSet(center + 1) {
				counts[i]++
			}
		}

		if i < len(counts)-cols {
			// check below if not the last row
			if i+cols == center {
				// if above center
				if s.next != nil {
					counts[i] += s.next.GetSideCount(Top, cols, counts)
				}
			} else {
				if s.IsSet(i + cols) {
					counts[i]++
				}
			}
		} else if s.prev != nil {
			// if last row and has previous level, add count of previous level
			// square below center
			if s.prev.IsSet(center + cols) {
				counts[i]++
			}
		}

	}
}

type Game struct {
	minutes int
	initial *State
	cols    int
	n       int
}

func (g *Game) GetMinutes() int {
	return g.minutes
}

func NewGame(area []string) *Game {
	g := new(Game)
	g.initial = NewState(0, nil, nil, 0)

	n := 0
	for _, r := range area {
		for _, c := range r {
			if c == '#' {
				g.initial.Set(n, true)
			}
			n++
		}
	}
	g.n = n

	g.cols = int(math.Sqrt(float64(n)))
	states.Add(0, g.initial.Clone())

	return g
}

func (g *Game) Run(mins int, print bool) {
	cols := g.cols
	center := g.n / 2

	if print {
		fmt.Printf("\nInitial State:")
		g.initial.Print(cols)
	}

	for {
		min := states.GetMin()
		max := states.GetMax()
		states.Add(min-1, NewState(min-1, nil, states.GetState(min), 0))
		states.Add(max+1, NewState(max+1, states.GetState(max), nil, 0))

		current := states.GetStates()
		counts := make([][]int, len(current), len(current))
		for i, s := range current {
			counts[i] = make([]int, g.n, g.n)
			s.GetCounts(cols, counts[i])
		}

		for i, s := range current {
			for n, c := range counts[i] {
				if n != center {
					if s.IsSet(n) && c != 1 {
						// a bug dies, unless exactly 1 bug is adjacent
						s.Set(n, false) // if already dead, no biggie
					} else if !s.IsSet(n) && (c == 1 || c == 2) {
						// empty spawns a bug if exactly 1 or 2 bugs are adjacent
						s.Set(n, true)
					}
				}
			}

		}

		g.minutes++

		if print {
			fmt.Printf("\n\nAfter %d minutes:\n", g.minutes)
			for _, s := range states.GetStates() {
				s.Print(cols)
			}
		}

		if g.minutes == mins {
			break
		}

	}

	sum := 0
	for _, s := range states.GetStates() {
		sum += s.GetCount()
	}
	fmt.Println("after ", mins, " minutes there are ", sum, " bugs")
}
