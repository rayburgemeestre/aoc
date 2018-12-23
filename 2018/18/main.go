package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

}

const (
	OpenGround rune = '.'
	Tree            = '|'
	Lumberyard      = '#'
)

var (
	world     [][]rune
	nextworld [][]rune
)

func solve(filename string, minutes int, visualize bool) int {
	world = [][]rune{}
	nextworld = [][]rune{}
	y := 0
	forEachLineInFile(filename, func(s string) bool {
		var (
			line  []rune
			line2 []rune
		)
		x := 0
		for _, c := range s {
			line = append(line, c)
			line2 = append(line2, c)
			x++
		}
		y++
		world = append(world, line)
		nextworld = append(nextworld, line2)
		return true
	})

	forEachAdjecentAcre := func(y int, x int, callback func(x int, y int, c rune)) {
		f := func(y int, x int) {
			if y < 0 || y >= len(world) || x < 0 || x >= len(world[y]) {
				return
			}
			callback(x, y, world[y][x])
		}
		f(y-1, x-1) // top left
		f(y-1, x)   // top
		f(y-1, x+1) // top right
		f(y, x-1)   // left
		f(y, x+1)   // right
		f(y+1, x-1) // bottom left
		f(y+1, x)   // bottom
		f(y+1, x+1) // bottom right
	}

	if visualize {
		fmt.Println("Initially")
		for y := 0; y < len(world); y++ {
			for x := 0; x < len(world[y]); x++ {
				fmt.Print(string(world[y][x]))
			}
			fmt.Println()
		}
	}

	wooded := 0
	lumberyard := 0
	previousAnswer := 0
	for round := 1; round <= minutes; round++ {
		for y := 0; y < len(world); y++ {
			for x := 0; x < len(world[y]); x++ {
				acre := world[y][x]
				rule1 := 0
				rule2 := 0
				rule3a := 0
				rule3b := 0
				becomeTree := false
				becomeLumberyard := false
				becomeOpen := false
				forEachAdjecentAcre(y, x, func(aX int, aY int, c rune) {
					// An open acre will become filled with trees if three or more adjacent acres contained trees. Otherwise, nothing happens.
					if acre == OpenGround && c == Tree {
						rule1++
						if rule1 >= 3 {
							becomeTree = true
						}
					}
					// An acre filled with trees will become a lumberyard if three or more adjacent acres were lumberyards. Otherwise, nothing happens.
					if acre == Tree && c == Lumberyard {
						rule2++
						if rule2 >= 3 {
							becomeLumberyard = true
						}
					}
					// An acre containing a lumberyard will remain a lumberyard if it was adjacent to at least one other lumberyard and at least one acre containing trees. Otherwise, it becomes open.
					if acre == Lumberyard {
						if c == Lumberyard {
							rule3a++
						} else if c == Tree {
							rule3b++
						}
						if rule3a > 0 && rule3b > 0 {
							becomeLumberyard = true
							becomeOpen = false
						} else {
							becomeLumberyard = false
							becomeOpen = true
						}
					}
				})
				if becomeTree {
					nextworld[y][x] = Tree
				} else if becomeLumberyard {
					nextworld[y][x] = Lumberyard
				} else if becomeOpen {
					nextworld[y][x] = OpenGround
				}
			}
		}

		wooded = 0
		lumberyard = 0
		for y := 0; y < len(world); y++ {
			for x := 0; x < len(world[y]); x++ {
				world[y][x] = nextworld[y][x]
				if world[y][x] == Tree {
					wooded++
				} else if world[y][x] == Lumberyard {
					lumberyard++
				}
			}
		}

		if visualize {
			fmt.Println("After minute #", round)
			for y := 0; y < len(world); y++ {
				for x := 0; x < len(world[y]); x++ {
					fmt.Print(string(world[y][x]))
				}
				fmt.Println()
			}
		}

		// Used for hand solving part two
		fmt.Println(round, "Answer", wooded*lumberyard-previousAnswer, wooded, lumberyard)

		// Above prints show there is a repeating pattern every 28 steps. For example 1009 it adds 28 compared to
		// the previous step:
		//1009 Answer 47 557 321
		// and 28 steps later, and later,...
		//1037 Answer 47 557 321
		//1065 Answer 47 557 321
		//1093 Answer 47 557 321

		// I picked step 1009:
		//1009 Answer 47 557 321

		// Figured between 1000000000 - 1009 there are 999998972 steps (rounded)
		// 999998972 * 28 = 999998972, and 1000000000 - 999998972 = 1028 (remaining),
		// Between 1028 and 1009 are 19 extra steps. We have to check 1009 + 19 (1028)
		//
		// 1028 Answer -5652 556 345

		// So now we can calculate it, 556 * 345 = 191820

		previousAnswer = wooded * lumberyard
	}
	fmt.Println("After all rounds:", wooded, lumberyard)
	fmt.Println("Answer:", wooded*lumberyard)

	return wooded * lumberyard
}

func forEachLineInFile(filename string, callback func(string) bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !callback(scanner.Text()) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
