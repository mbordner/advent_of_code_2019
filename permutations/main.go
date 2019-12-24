package main

import "fmt"

// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func factorial(n int64) int64 {
	if n > 1 {
		return n * factorial(n-1)
	}
	return n
}

func main() {
	array := []int{0,1,2,3,4}

	permutations := make([][]int,0,factorial(int64(len(array))))

	for p := make([]int, len(array)); p[0] < len(p); nextPerm(p) {
		permutation := getPerm(array, p)
		permutations = append(permutations,permutation)
	}

	fmt.Println("number of permutations: ",len(permutations))
	fmt.Println(permutations)

}