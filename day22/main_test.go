package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Create_Deck(t *testing.T) {
	assert.Equal(t, []int64{0, 1, 2}, createDeck(3))
}

func Test_Deal_Into_New_Stack(t *testing.T) {
	assert.Equal(t, []int64{3, 2, 1}, dealIntoNewStack([]int64{1, 2, 3}))
	assert.Equal(t, []int64{4, 3, 2, 1}, dealIntoNewStack([]int64{1, 2, 3, 4}))
}

func Test_Cut(t *testing.T) {
	assert.Equal(t, []int64{3, 4, 5, 6, 7, 8, 9, 0, 1, 2}, cut([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, int64(3)))
	assert.Equal(t, []int64{6, 7, 8, 9, 0, 1, 2, 3, 4, 5}, cut([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, int64(-4)))
}

func Test_Deal_With_Increment(t *testing.T) {
	assert.Equal(t, []int64{0, 7, 4, 1, 8, 5, 2, 9, 6, 3}, dealWithIncrement([]int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, int64(3)))
}

func Test_Shuffles(t *testing.T) {
	table := []struct {
		shuffles string
		expected []int64
	}{
		{`deal with increment 7
deal into new stack
deal into new stack`, []int64{0, 3, 6, 9, 2, 5, 8, 1, 4, 7},},

		{`cut 6
deal with increment 7
deal into new stack`, []int64{3, 0, 7, 4, 1, 8, 5, 2, 9, 6},},

		{`deal with increment 7
deal with increment 9
cut -2`, []int64{6, 3, 0, 7, 4, 1, 8, 5, 2, 9},},

		{`deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1`, []int64{9, 2, 5, 8, 1, 4, 7, 0, 3, 6},},
	}

	for _, test := range table {

		deck := shuffle(createDeck(int64(len(test.expected))), strings.Split(test.shuffles, "\n"))

		assert.Equal(t, test.expected, deck)

	}
}
