package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	EMPTY      = iota
	HORIZONTAL = iota
	VERTICAL   = iota
	CORNER     = iota
	CROSSROAD  = iota
)

const (
	UP    = iota
	DOWN  = iota
	LEFT  = iota
	RIGHT = iota
)

type Track struct {
	direction int
	c         rune
}

const (
	NONE            = iota
	FIRST_LEFT      = iota
	SECOND_STRAIGHT = iota
	THIRD_RIGHT     = iota
)

type Cart struct {
	direction int
	x         int
	y         int
	c         rune
	decision  int
	collide   bool
}

type Index struct {
	trackDirection  int
	cartDirection   int
	trackCharacter  rune
	currentDecision int
}

type Change struct {
	cartDirection int
	cartCharacter rune
	nextDecision  int
}

func main() {
	var world [][]Track
	var carts []Cart

	lookup := map[Index]Change{
		Index{CORNER, LEFT, '/', NONE}:                Change{DOWN, 'v', NONE},
		Index{CORNER, LEFT, '/', NONE}:                Change{DOWN, 'v', NONE},
		Index{CORNER, LEFT, '\\', NONE}:               Change{UP, '^', NONE},
		Index{CORNER, RIGHT, '/', NONE}:               Change{UP, '^', NONE},
		Index{CORNER, RIGHT, '\\', NONE}:              Change{DOWN, 'v', NONE},
		Index{CORNER, UP, '/', NONE}:                  Change{RIGHT, '>', NONE},
		Index{CORNER, UP, '\\', NONE}:                 Change{LEFT, '<', NONE},
		Index{CORNER, DOWN, '/', NONE}:                Change{LEFT, '<', NONE},
		Index{CORNER, DOWN, '\\', NONE}:               Change{RIGHT, '>', NONE},
		Index{CROSSROAD, LEFT, '+', FIRST_LEFT}:       Change{DOWN, 'v', SECOND_STRAIGHT},
		Index{CROSSROAD, LEFT, '+', SECOND_STRAIGHT}:  Change{LEFT, '<', THIRD_RIGHT},
		Index{CROSSROAD, LEFT, '+', THIRD_RIGHT}:      Change{UP, '^', FIRST_LEFT},
		Index{CROSSROAD, RIGHT, '+', FIRST_LEFT}:      Change{UP, '^', SECOND_STRAIGHT},
		Index{CROSSROAD, RIGHT, '+', SECOND_STRAIGHT}: Change{RIGHT, '>', THIRD_RIGHT},
		Index{CROSSROAD, RIGHT, '+', THIRD_RIGHT}:     Change{DOWN, 'v', FIRST_LEFT},
		Index{CROSSROAD, UP, '+', FIRST_LEFT}:         Change{LEFT, '<', SECOND_STRAIGHT},
		Index{CROSSROAD, UP, '+', SECOND_STRAIGHT}:    Change{UP, '^', THIRD_RIGHT},
		Index{CROSSROAD, UP, '+', THIRD_RIGHT}:        Change{RIGHT, '>', FIRST_LEFT},
		Index{CROSSROAD, DOWN, '+', FIRST_LEFT}:       Change{RIGHT, '>', SECOND_STRAIGHT},
		Index{CROSSROAD, DOWN, '+', SECOND_STRAIGHT}:  Change{DOWN, 'v', THIRD_RIGHT},
		Index{CROSSROAD, DOWN, '+', THIRD_RIGHT}:      Change{LEFT, '<', FIRST_LEFT},
	}

	y := 0
	forEachLineInFile("input", func(line string) {
		world = append(world, []Track{})
		x := 0
		for _, c := range line {
			var track Track
			switch c {
			case '|':
				track = Track{VERTICAL, c}
			case '-':
				track = Track{HORIZONTAL, c}
			case '/':
				track = Track{CORNER, c}
			case '\\':
				track = Track{CORNER, c}
			case '+':
				track = Track{CROSSROAD, c}
			case '^':
				track = Track{VERTICAL, '|'}
				carts = append(carts, Cart{UP, x, y, '^', FIRST_LEFT, false})
			case 'v':
				track = Track{VERTICAL, '|'}
				carts = append(carts, Cart{DOWN, x, y, 'v', FIRST_LEFT, false})
			case '<':
				track = Track{HORIZONTAL, '-'}
				carts = append(carts, Cart{LEFT, x, y, '<', FIRST_LEFT, false})
			case '>':
				track = Track{HORIZONTAL, '-'}
				carts = append(carts, Cart{RIGHT, x, y, '>', FIRST_LEFT, false})
			default:
				track = Track{EMPTY, 0}
			}
			world[y] = append(world[y], track)
			x++
		}
		y++
	})

	partOneFound := false
	carsOnMap := len(carts)
	for iter := 0; ; iter++ {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y == carts[j].y {
				return carts[i].x < carts[j].x
			}
			return carts[i].y < carts[j].y
		})

		for cartIndex := 0; cartIndex < len(carts); cartIndex++ {
			cart := carts[cartIndex]
			if cart.collide {
				continue
			}
			track := world[cart.y][cart.x]

			index := Index{track.direction, cart.direction, track.c, cart.decision}
			if track.direction == CORNER {
				index.currentDecision = NONE
			}
			if _, exists := lookup[index]; exists {
				cart.direction = lookup[index].cartDirection
				cart.c = lookup[index].cartCharacter
				if lookup[index].nextDecision != NONE {
					cart.decision = lookup[index].nextDecision
				}
			}

			switch cart.direction {
			case LEFT:
				cart.x--
			case RIGHT:
				cart.x++
			case UP:
				cart.y--
			case DOWN:
				cart.y++
			default:
			}

			for j, cart2 := range carts {
				if cart2.collide {
					continue
				}
				if cart.x == cart2.x && cart.y == cart2.y {
					fmt.Println("Collision!")
					cart.collide, carts[j].collide = true, true
					carsOnMap -= 2
					if !partOneFound {
						fmt.Println("Part one, first collision:")
						fmt.Println("Found at iteration", iter)
						fmt.Printf("Answer: %d,%d\n", cart.x, cart.y)
						partOneFound = true
					}
				}
			}
			carts[cartIndex] = cart
		}

		if carsOnMap == 1 {
			fmt.Println("Part two, last car:")
			fmt.Println("Found in iteration", iter)
			for _, cart := range carts {
				if !cart.collide {
					fmt.Printf("Answer: %d,%d\n", cart.x, cart.y)
					return
				}
			}
		}
	}
}

func forEachLineInFile(filename string, callback func(string)) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		callback(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// removed visualization code
//var cart *Cart
//for y = 0; y < len(world); y++ {
//	for x := 0; x < len(world[y]); x++ {
//		track := world[y][x]
//		for _, c := range carts {
//			if c.collide {
//				continue
//			}
//			if x == c.x && y == c.y {
//				cart = &c
//				break
//			}
//		}
//		if cart != nil {
//			fmt.Print(string(cart.c))
//		} else if track.c == 0 {
//			fmt.Print(" ")
//		} else {
//			fmt.Print(string(track.c))
//		}
//	}
//	fmt.Println()
//}
