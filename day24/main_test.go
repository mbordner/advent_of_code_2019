package main

import (
	"github.com/mbordner/advent_of_code_2019/day24/part1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Part1_State(t *testing.T) {
	s := part1.NewState(nil)

	s.Set(0, true)
	assert.Equal(t, 1, int(s))
	s.Set(1, true)
	assert.Equal(t, 3, int(s))

	assert.True(t,s.IsSet(1))
	assert.False(t,s.IsSet(2))

	s.Set(2,true)
	assert.Equal(t,7,int(s))

	var counts [4]int
	assert.Equal(t,0,counts[0])
	s.GetCounts(2,counts[:])
	assert.Equal(t,2,counts[0])
	assert.Equal(t,1,counts[1])
	assert.Equal(t,1,counts[2])
	assert.Equal(t,2,counts[3])

	s.Set(3,true)

	var counts2 [9]int
	s.GetCounts(3,counts2[:])
	assert.Equal(t,2,counts2[4])
}

func Test_Part1_Test1(t *testing.T) {
	g := part1.NewGame(getGame("test1.txt"))
	g.Run(false)
	assert.Equal(t,86,g.GetMinutes())
	assert.Equal(t,2129920,int(g.GetState()))
}
