package main

import (
	"errors"
	"fmt"
	"math/big"
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


func polypow(a,b,n,m int64) (int64,int64) {
	M := big.NewInt(m)

	if n == 0 {
		return 1,0
	} else if n%2 == 0 {

		X := big.NewInt(a)
		X = X.Mul(X,X).Mod(X,M)

		Y := big.NewInt(a)
		Y = Y.Mul(Y,big.NewInt(b)).Add(Y,big.NewInt(b)).Mod(Y,M)

		return polypow(X.Int64(),Y.Int64(),n/2,m)
	} else {
		c,d := polypow(a,b,n-1,m)

		X := big.NewInt(a)
		Y := big.NewInt(a)

		return X.Mul(X,big.NewInt(c)).Mod(X,M).Int64(),
			Y.Mul(Y,big.NewInt(d)).Add(Y,big.NewInt(b)).Mod(Y,M).Int64()
	}
}



func calculateInverseFunctionCoefficients(m int64, shuffles []string) (a,b int64) {
	// convert rules to linear polynomial.
	// (g∘f)(x) = g(f(x))
	a, b = int64(1),int64(0)

	M := big.NewInt(m)

	for i := len(shuffles) - 1; i >= 0; i-- {
		if matches := reCut.FindStringSubmatch(shuffles[i]); len(matches) > 0 {
			n, e := strconv.ParseInt(matches[1], 10, 64)
			if e != nil {
				panic(e)
			}

			b = (b+n)%m

		} else if matches := reDealNew.FindStringSubmatch(shuffles[i]); len(matches) > 0 {

			a = -a
			b = m-b-1

		} else if matches := reDealWInc.FindStringSubmatch(shuffles[i]); len(matches) > 0 {
			n, e := strconv.ParseInt(matches[1], 10, 64)
			if e != nil {
				panic(e)
			}

			N := big.NewInt(n)
			Z := N.ModInverse(N,M)
			A := big.NewInt(a)
			B := big.NewInt(b)

			/**
			func (z *Int) ModInverse(g, n *Int) *Int
			ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ and returns z. If g and n are not relatively prime, g has no multiplicative inverse in the ring ℤ/nℤ. In this case, z is unchanged and the return value is nil.
			 */

			if Z == nil {
				panic(errors.New("couldn't find an inverse"))
			}

			A = A.Mul(A,Z).Mod(A,M)
			B = B.Mul(B,Z).Mod(B,M)

			a = A.Int64()
			b = B.Int64()
		}

	}

	return a,b
}


func main() {
	decksize := int64(119315717514047)
	numShuffles := int64(101741582076661)

	a,b := calculateInverseFunctionCoefficients(decksize,getShuffles("input.txt"))
	a,b = polypow(a,b,numShuffles,decksize)

	A := big.NewInt(a)
	A = A.Mul(A,big.NewInt(int64(2020))).Add(A,big.NewInt(b)).Mod(A,big.NewInt(decksize))


	fmt.Println(A.Int64())

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
