package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	input = `143845
86139
53043
124340
73213
108435
126874
131397
85618
107774
66872
94293
51015
51903
147655
112891
100993
143374
83737
145868
144768
89793
124127
135366
94017
81678
102325
75394
103852
81896
148050
142780
50503
110691
117851
137382
92841
138222
128414
146834
59968
136456
122397
147157
83595
59916
75690
125025
147797
112494
76247
100221
63389
59070
97466
91905
126234
76561
128170
102778
82342
131097
51609
148204
74812
64925
127927
79056
73307
78431
88770
97688
103564
76001
105232
145361
77845
87518
117293
110054
135599
85005
85983
118255
103031
142440
140505
99614
69593
69161
78795
54808
115582
117976
148858
84193
147285
89038
92677
106574`
)

func main() {
	total1, total2 := 0, 0
	inputs := strings.Split(input, "\n")
	for _, s := range inputs {
		i, e := strconv.ParseInt(s, 10, 32)
		if e != nil {
			panic(e)
		}
		total1 += puzzle(int(i))
		total2 += puzzle2(int(i))
	}
	fmt.Println("puzzle 1 total:",total1)
	fmt.Println("puzzle 2 total:",total2)
}

/**
--- Day 1: The Tyranny of the Rocket Equation ---
Santa has become stranded at the edge of the Solar System while delivering presents to other planets! To accurately calculate his position in space, safely align his warp drive, and return to Earth in time to save Christmas, he needs you to bring him measurements from fifty stars.

Collect stars by solving puzzles. Two puzzles will be made available on each day in the Advent calendar; the second puzzle is unlocked when you complete the first. Each puzzle grants one star. Good luck!

The Elves quickly load you into a spacecraft and prepare to launch.

At the first Go / No Go poll, every Elf is Go until the Fuel Counter-Upper. They haven't determined the amount of fuel required yet.

Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module, take its mass, divide by three, round down, and subtract 2.

For example:

For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
For a mass of 1969, the fuel required is 654.
For a mass of 100756, the fuel required is 33583.
The Fuel Counter-Upper needs to know the total fuel requirement. To find it, individually calculate the fuel needed for the mass of each module (your puzzle input), then add together all the fuel values.

What is the sum of the fuel requirements for all of the modules on your spacecraft?

Your puzzle answer was 3427972.

--- Part Two ---
During the second Go / No Go poll, the Elf in charge of the Rocket Equation Double-Checker stops the launch sequence. Apparently, you forgot to include additional fuel for the fuel you just added.

Fuel itself requires fuel just like a module - take its mass, divide by three, round down, and subtract 2. However, that fuel also requires fuel, and that fuel requires fuel, and so on. Any mass that would require negative fuel should instead be treated as if it requires zero fuel; the remaining mass, if any, is instead handled by wishing really hard, which has no mass and is outside the scope of this calculation.

So, for each module mass, calculate its fuel and add it to the total. Then, treat the fuel amount you just calculated as the input mass and repeat the process, continuing until a fuel requirement is zero or negative. For example:

A module of mass 14 requires 2 fuel. This fuel requires no further fuel (2 divided by 3 and rounded down is 0, which would call for a negative fuel), so the total fuel required is still just 2.
At first, a module of mass 1969 requires 654 fuel. Then, this fuel requires 216 more fuel (654 / 3 - 2). 216 then requires 70 more fuel, which requires 21 fuel, which requires 5 fuel, which requires no further fuel. So, the total fuel required for a module of mass 1969 is 654 + 216 + 70 + 21 + 5 = 966.
The fuel required by a module of mass 100756 and its fuel is: 33583 + 11192 + 3728 + 1240 + 411 + 135 + 43 + 12 + 2 = 50346.
What is the sum of the fuel requirements for all of the modules on your spacecraft when also taking into account the mass of the added fuel? (Calculate the fuel requirements for each module separately, then add them all up at the end.)

Your puzzle answer was 5139078.

Both parts of this puzzle are complete! They provide two gold stars: **
 */