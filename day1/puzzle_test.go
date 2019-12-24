package main

import (
	"fmt"
	"testing"
)

type Input struct {
	value int
}

func (i Input) String() string {
	return fmt.Sprintf("{[%d]}", i.value)
}

func Test_puzzle(t *testing.T) {

	table := []struct {
		inputs    Input
		expected int
	}{
		{Input{1969},654},
		{Input{ 100756}, 33583},
	}

	for _, test := range table {

		val := puzzle(test.inputs.value)

		if val != test.expected {
			t.Log("Expected: ", test.expected)
			t.Log("Actual: ", val)
		}

	}
}

func Test_puzzle2(t *testing.T) {

	table := []struct {
		inputs    Input
		expected int
	}{
		{Input{14},2},
		{Input{1969},966},
		{Input{ 100756}, 50346},
	}

	for _, test := range table {

		val := puzzle2(test.inputs.value)

		if val != test.expected {
			t.Log("Expected: ", test.expected)
			t.Log("Actual: ", val)
		}

	}
}