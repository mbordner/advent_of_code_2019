package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Movement(t *testing.T) {

	r := NewRobot()

	r.TurnAndAdvance(Left)

	assert.Equal(t,Pos{-1,0},r.pos)
	assert.Equal(t,Left,r.dir)

	r.TurnAndAdvance(Left)

	assert.Equal(t,Pos{-1,-1},r.pos)
	assert.Equal(t,Down,r.dir)

	r.TurnAndAdvance(Left)

	assert.Equal(t,Pos{0,-1},r.pos)
	assert.Equal(t,Right,r.dir)

	r.TurnAndAdvance(Left)

	assert.Equal(t,Pos{0,0},r.pos)
	assert.Equal(t,Up,r.dir)



	r.TurnAndAdvance(Right)

	assert.Equal(t,Pos{1,0},r.pos)
	assert.Equal(t,Right,r.dir)

	r.TurnAndAdvance(Right)

	assert.Equal(t,Pos{1,-1},r.pos)
	assert.Equal(t,Down,r.dir)

	r.TurnAndAdvance(Right)

	assert.Equal(t,Pos{0,-1},r.pos)
	assert.Equal(t,Left,r.dir)

	r.TurnAndAdvance(Right)

	assert.Equal(t,Pos{0,0},r.pos)
	assert.Equal(t,Up,r.dir)
}