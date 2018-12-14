package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	UP    = iota
	DOWN  = iota
	LEFT  = iota
	RIGHT = iota
)

const (
	__           = iota
	FIRST_LEFT   = iota
	SECOND_AHEAD = iota
	THIRD_RIGHT  = iota
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
	var world [][]rune
	var carts []Cart

	// Lookup table for determining the changes to the cart (new direction, ascii art and next decision),
	// based on its current state (direction, track it's on and it's planned decision)
	lookup := map[Index]Change{
		Index{LEFT, '/', __}:            {DOWN, 'v', __},
		Index{LEFT, '/', __}:            {DOWN, 'v', __},
		Index{LEFT, '\\', __}:           {UP, '^', __},
		Index{RIGHT, '/', __}:           {UP, '^', __},
		Index{RIGHT, '\\', __}:          {DOWN, 'v', __},
		Index{UP, '/', __}:              {RIGHT, '>', __},
		Index{UP, '\\', __}:             {LEFT, '<', __},
		Index{DOWN, '/', __}:            {LEFT, '<', __},
		Index{DOWN, '\\', __}:           {RIGHT, '>', __},
		Index{LEFT, '+', FIRST_LEFT}:    {DOWN, 'v', SECOND_AHEAD},
		Index{LEFT, '+', SECOND_AHEAD}:  {LEFT, '<', THIRD_RIGHT},
		Index{LEFT, '+', THIRD_RIGHT}:   {UP, '^', FIRST_LEFT},
		Index{RIGHT, '+', FIRST_LEFT}:   {UP, '^', SECOND_AHEAD},
		Index{RIGHT, '+', SECOND_AHEAD}: {RIGHT, '>', THIRD_RIGHT},
		Index{RIGHT, '+', THIRD_RIGHT}:  {DOWN, 'v', FIRST_LEFT},
		Index{UP, '+', FIRST_LEFT}:      {LEFT, '<', SECOND_AHEAD},
		Index{UP, '+', SECOND_AHEAD}:    {UP, '^', THIRD_RIGHT},
		Index{UP, '+', THIRD_RIGHT}:     {RIGHT, '>', FIRST_LEFT},
		Index{DOWN, '+', FIRST_LEFT}:    {RIGHT, '>', SECOND_AHEAD},
		Index{DOWN, '+', SECOND_AHEAD}:  {DOWN, 'v', THIRD_RIGHT},
		Index{DOWN, '+', THIRD_RIGHT}:   {LEFT, '<', FIRST_LEFT},
	}

	// Read all "world" data
	y := 0
	forEachLineInFile("input", func(line string) {
		world = append(world, []rune{})
		x := 0
		for _, c := range line {
			switch c {
			case '^':
				c, carts = '|', append(carts, Cart{UP, x, y, '^', FIRST_LEFT, false})
			case 'v':
				c, carts = '|', append(carts, Cart{DOWN, x, y, 'v', FIRST_LEFT, false})
			case '<':
				c, carts = '-', append(carts, Cart{LEFT, x, y, '<', FIRST_LEFT, false})
			case '>':
				c, carts = '-', append(carts, Cart{RIGHT, x, y, '>', FIRST_LEFT, false})
			}
			world[y] = append(world[y], c)
			x++
		}
		y++
	})

	partOneFound := false
	carsOnMap := len(carts)
	for iter := 0; ; iter++ {
		// Make sure the carts slice is sorted by (row, col) because it's important to determine the order in which
		// the carts move
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y == carts[j].y {
				return carts[i].x < carts[j].x
			}
			return carts[i].y < carts[j].y
		})

		for cartIndex, cart := range carts {
			if cart.collide {
				continue
			}

			track := world[cart.y][cart.x]
			index := Index{cart.direction, track, cart.decision}
			if track == '/' || track == '\\' {
				index.currentDecision = __
			}

			// Update cart state
			if _, exists := lookup[index]; exists {
				next := lookup[index]
				cart.direction = next.cartDirection
				cart.c = next.cartCharacter
				if next.nextDecision != __ {
					cart.decision = next.nextDecision
				}
			}

			// Move cart
			switch cart.direction {
			case LEFT:
				cart.x--
			case RIGHT:
				cart.x++
			case UP:
				cart.y--
			case DOWN:
				cart.y++
			}

			// Check for collisions
			for i, cart2 := range carts {
				if !cart2.collide && cart.x == cart2.x && cart.y == cart2.y {
					fmt.Println("Collision at", cart.x, cart.y)
					cart.collide, carts[i].collide = true, true
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
			for _, cart := range carts {
				if !cart.collide {
					fmt.Println("Part two, last car:")
					fmt.Println("Found in iteration", iter)
					fmt.Printf("Answer: %d,%d\n", cart.x, cart.y)
					//visualize(y, world, carts)
					return
				}
			}
		}
	}
}

func visualize(y int, world [][]rune, carts []Cart) {
	for y = 0; y < len(world); y++ {
		for x := 0; x < len(world[y]); x++ {
			track := world[y][x]
			cart := func() *Cart {
				for _, c := range carts {
					if !c.collide && x == c.x && y == c.y {
						return &c
					}
				}
				return nil
			}()
			if cart != nil {
				fmt.Print(string(cart.c))
			} else if track == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(string(track))
			}
		}
		fmt.Println()
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
