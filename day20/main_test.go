package main

import (
	"github.com/mbordner/advent_of_code_2019/day20/geom"
	"github.com/mbordner/advent_of_code_2019/day20/part1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Part1_Test_Case_1(t *testing.T) {
	maze := getMaze("test1.txt")
	game := part1.NewGame(maze)

	pd := game.GetPortalData("AA")
	assert.NotNil(t, pd)

	n := pd.GetPathNode()
	sid := n.GetID().(geom.Pos).String()

	ot := pd.Nodes[0].GetProperty("type").(part1.ObjectType)
	assert.Equal(t, part1.HalfPortal, ot)

	assert.Equal(t, "{x:9, y:2}", sid)

	pd = game.GetPortalData("BC")
	assert.NotNil(t, pd)

	n = pd.GetPathNode()
	sid = n.GetID().(geom.Pos).String()

	ot = pd.Nodes[0].GetProperty("type").(part1.ObjectType)
	assert.Equal(t, part1.Portal, ot)

	assert.Contains(t, []string{"{x:9, y:6}", "{x:2, y:8}"}, sid)

	pd = game.GetPortalData("DE")
	assert.NotNil(t, pd)

	n = pd.GetPathNode()
	sid = n.GetID().(geom.Pos).String()

	assert.Contains(t, []string{"{x:6, y:10}", "{x:2, y:13}"}, sid)

	pd = game.GetPortalData("ZZ")
	assert.NotNil(t, pd)

	n = pd.GetPathNode()
	sid = n.GetID().(geom.Pos).String()

	assert.Equal(t, "{x:13, y:16}", sid)

	assert.Equal(t, 63, game.GameGraph.Len())

	path, distance := game.ShortestPath("AA", "ZZ")

	assert.NotNil(t, path)
	assert.NotZero(t, distance)
	assert.Equal(t, 23, distance)
}

func Test_Part1_Test_Case_2(t *testing.T) {
	maze := getMaze("test2.txt")
	game := part1.NewGame(maze)

	path, distance := game.ShortestPath("AA", "ZZ")

	assert.NotNil(t, path)
	assert.Equal(t, 58, distance)
}
