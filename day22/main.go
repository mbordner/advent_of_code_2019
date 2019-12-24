package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reCut      = regexp.MustCompile(`cut\s(-?\d+)`)
	reDealNew  = regexp.MustCompile(`(deal into new stack)`)
	reDealWInc = regexp.MustCompile(`deal with increment (\d+)`)
)

func dealIntoNewStack(deck []int64) []int64 {
	// this is just reversing the deck
	for i, j, h := 0, len(deck)-1, len(deck)/2; i < h; i, j = i+1, j-1 {
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func createDeck(n int64) []int64 {
	deck := make([]int64, n, n)
	for i := int64(0); i < n; i++ {
		deck[i] = i
	}
	return deck
}

func cut(deck []int64, n int64) []int64 {
	if n > 0 {
		return append(deck[n:], deck[:n]...)
	}
	return append(deck[int64(len(deck))+n:], deck[:int64(len(deck))+n]...)
}

func dealWithIncrement(deck []int64, n int64) []int64 {
	newDeck := make([]int64, len(deck), len(deck))
	pos := int64(0)
	for i := int64(0); i < int64(len(deck)); i++ {
		newDeck[pos] = deck[i]
		pos += n
		if pos >= int64(len(deck)) {
			pos -= int64(len(deck))
		}
	}
	return newDeck
}

func shuffle(deck []int64, shuffles []string) []int64 {
	for i := range shuffles {
		if matches := reCut.FindStringSubmatch(shuffles[i]); len(matches) > 0 {
			n, e := strconv.ParseInt(matches[1], 10, 64)
			if e != nil {
				panic(e)
			}
			deck = cut(deck, n)
		} else if matches := reDealNew.FindStringSubmatch(shuffles[i]); len(matches) > 0 {
			deck = dealIntoNewStack(deck)
		} else if matches := reDealWInc.FindStringSubmatch(shuffles[i]); len(matches) > 0 {
			n, e := strconv.ParseInt(matches[1], 10, 64)
			if e != nil {
				panic(e)
			}
			deck = dealWithIncrement(deck, n)
		}
	}
	return deck
}

func part1() {
	deck := createDeck(10007)
	shuffles := getShuffles("input.txt")

	deck = shuffle(deck, shuffles)

	for i := 0; i < len(deck); i++ {
		if deck[i] == 2019 {
			fmt.Println("part 1 answer: ", i)
			break
		}
	}
}

func equal(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func main() {
	shuffles := getShuffles("input.txt")
	nums := []int64{10007, 10009, 10037, 10039, 10061, 10067, 10069, 10079, 10091, 10093, 10099, 10103, 10111, 10133}

	for _, v := range nums {
		deck := createDeck(v)
		deck = shuffle(deck, shuffles)
		fmt.Println("prime: ",v," with differing sequence of: ",deck[1]-deck[0],deck[2]-deck[1],deck[3]-deck[2],deck[4]-deck[3],deck[5]-deck[4],deck[6]-deck[5])
	}

	deck := createDeck(10007)
	deck = shuffle(deck, shuffles)






	for i:=len(deck)-1; i > 0; i-- {
		fmt.Print(deck[i]-deck[i-1]," ")
	}



	//deck = createDeck(int64(119315717514047))
	//deck = shuffle(deck,shuffles)

}

func getShuffles(filename string) []string {
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

	return strings.Split(string(buffer), "\n")
}
