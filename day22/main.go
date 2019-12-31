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

/**
--- Day 22: Slam Shuffle ---
There isn't much to do while you wait for the droids to repair your ship. At least you're drifting in the right direction. You decide to practice a new card shuffle you've been working on.

Digging through the ship's storage, you find a deck of space cards! Just like any deck of space cards, there are 10007 cards in the deck numbered 0 through 10006. The deck must be new - they're still in factory order, with 0 on the top, then 1, then 2, and so on, all the way through to 10006 on the bottom.

You've been practicing three different techniques that you use while shuffling. Suppose you have a deck of only 10 cards (numbered 0 through 9):

To deal into new stack, create a new stack of cards by dealing the top card of the deck onto the top of the new stack repeatedly until you run out of cards:

Top          Bottom
0 1 2 3 4 5 6 7 8 9   Your deck
                      New stack

  1 2 3 4 5 6 7 8 9   Your deck
                  0   New stack

    2 3 4 5 6 7 8 9   Your deck
                1 0   New stack

      3 4 5 6 7 8 9   Your deck
              2 1 0   New stack

Several steps later...

                  9   Your deck
  8 7 6 5 4 3 2 1 0   New stack

                      Your deck
9 8 7 6 5 4 3 2 1 0   New stack
Finally, pick up the new stack you've just created and use it as the deck for the next technique.

To cut N cards, take the top N cards off the top of the deck and move them as a single unit to the bottom of the deck, retaining their order. For example, to cut 3:

Top          Bottom
0 1 2 3 4 5 6 7 8 9   Your deck

      3 4 5 6 7 8 9   Your deck
0 1 2                 Cut cards

3 4 5 6 7 8 9         Your deck
              0 1 2   Cut cards

3 4 5 6 7 8 9 0 1 2   Your deck
You've also been getting pretty good at a version of this technique where N is negative! In that case, cut (the absolute value of) N cards from the bottom of the deck onto the top. For example, to cut -4:

Top          Bottom
0 1 2 3 4 5 6 7 8 9   Your deck

0 1 2 3 4 5           Your deck
            6 7 8 9   Cut cards

        0 1 2 3 4 5   Your deck
6 7 8 9               Cut cards

6 7 8 9 0 1 2 3 4 5   Your deck
To deal with increment N, start by clearing enough space on your table to lay out all of the cards individually in a long line. Deal the top card into the leftmost position. Then, move N positions to the right and deal the next card there. If you would move into a position past the end of the space on your table, wrap around and keep counting from the leftmost card again. Continue this process until you run out of cards.

For example, to deal with increment 3:


0 1 2 3 4 5 6 7 8 9   Your deck
. . . . . . . . . .   Space on table
^                     Current position

Deal the top card to the current position:

  1 2 3 4 5 6 7 8 9   Your deck
0 . . . . . . . . .   Space on table
^                     Current position

Move the current position right 3:

  1 2 3 4 5 6 7 8 9   Your deck
0 . . . . . . . . .   Space on table
      ^               Current position

Deal the top card:

    2 3 4 5 6 7 8 9   Your deck
0 . . 1 . . . . . .   Space on table
      ^               Current position

Move right 3 and deal:

      3 4 5 6 7 8 9   Your deck
0 . . 1 . . 2 . . .   Space on table
            ^         Current position

Move right 3 and deal:

        4 5 6 7 8 9   Your deck
0 . . 1 . . 2 . . 3   Space on table
                  ^   Current position

Move right 3, wrapping around, and deal:

          5 6 7 8 9   Your deck
0 . 4 1 . . 2 . . 3   Space on table
    ^                 Current position

And so on:

0 7 4 1 8 5 2 9 6 3   Space on table
Positions on the table which already contain cards are still counted; they're not skipped. Of course, this technique is carefully designed so it will never put two cards in the same position or leave a position empty.

Finally, collect the cards on the table so that the leftmost card ends up at the top of your deck, the card to its right ends up just below the top card, and so on, until the rightmost card ends up at the bottom of the deck.

The complete shuffle process (your puzzle input) consists of applying many of these techniques. Here are some examples that combine techniques; they all start with a factory order deck of 10 cards:

deal with increment 7
deal into new stack
deal into new stack
Result: 0 3 6 9 2 5 8 1 4 7
cut 6
deal with increment 7
deal into new stack
Result: 3 0 7 4 1 8 5 2 9 6
deal with increment 7
deal with increment 9
cut -2
Result: 6 3 0 7 4 1 8 5 2 9
deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1
Result: 9 2 5 8 1 4 7 0 3 6
Positions within the deck count from 0 at the top, then 1 for the card immediately below the top card, and so on to the bottom. (That is, cards start in the position matching their number.)

After shuffling your factory order deck of 10007 cards, what is the position of card 2019?

Your puzzle answer was 2604.

--- Part Two ---
After a while, you realize your shuffling skill won't improve much more with merely a single deck of cards. You ask every 3D printer on the ship to make you some more cards while you check on the ship repairs. While reviewing the work the droids have finished so far, you think you see Halley's Comet fly past!

When you get back, you discover that the 3D printers have combined their power to create for you a single, giant, brand new, factory order deck of 119315717514047 space cards.

Finally, a deck of cards worthy of shuffling!

You decide to apply your complete shuffle process (your puzzle input) to the deck 101741582076661 times in a row.

You'll need to be careful, though - one wrong move with this many cards and you might overflow your entire ship!

After shuffling your new, giant, factory order deck that many times, what number is on the card that ends up in position 2020?

Your puzzle answer was 79608410258462.

Both parts of this puzzle are complete! They provide two gold stars: **
 */