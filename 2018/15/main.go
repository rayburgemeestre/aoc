package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type Unit struct {
	c  rune
	HP int
	x  int
	y  int
}

type Tile struct {
	c    rune
	unit *Unit
	x    int
	y    int
	cost int

	px int
	py int
}

func main() {
	world := [][]Tile{}
	units := []*Unit{}

	y := 0
//	forEachLineInFile("input_test", func(s string) {
		forEachLineInFile("input_test", func(s string) {
		world = append(world, []Tile{})
		x := 0
		for _, c := range s {
			//fmt.Println(string(c))
			switch c {
			case '#':
				world[y] = append(world[y], Tile{c, nil, x, y, 0, 0, 0})
			default:
				var unit *Unit
				if c == 'E' || c == 'G' {
					units = append(units, &Unit{c, 200, x, y})
					unit = units[len(units)-1]
					unit.x, unit.y = x, y
					unit.y = y
				}
				world[y] = append(world[y], Tile{'.', unit, x, y, 0, 0, 0})

			}
			x++
		}
		y++
	})

	fmt.Println("Initially:")
	visualize(&world)
	for round := 0; ; round++ {

		// sort units
		sort.Slice(units, func(i, j int) bool {
			if units[i].y == units[j].y {
				return units[i].x < units[j].x
			}
			return units[i].y < units[j].y
		})

		for i, attacker := range units {
			fmt.Println("IN: attacker x,y",attacker.x,attacker.y, " ",string(attacker.c))
			if attacker.HP <= 0 {
				fmt.Println("DEad")
				continue
			}
			//if attacker.x == 7 && attacker.y == 3 {
			//
			//} else {
			//	continue
			//}
			// TODO: we can skip units we already compared, convert into manual loop I suppose
			// for now let's keep it like this.

			min := math.MaxInt32
			//bestVictim := -1
			var nextMove CoordCost
			hasAdjecentEnemy := false
			for _, victim := range units {
				if victim.HP <= 0 {
					continue
				}
				if victim.c == attacker.c {
					continue // same team
				}
				if victim.x == attacker.x && victim.y == attacker.y {
					continue // cannot attack yourself
				}
				// TODO: check if this can be simplified!
				//fmt.Printf("Checking distance between (%d,%d) - (%d,%d)!\n", attacker.x, attacker.y, victim.x, victim.y)

				checkAttackVector := func(x int, y int, h int) {
					if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
						if world[y][x].c != '.' {
							return
						}
						if world[y][x].unit != nil {
							return
						}
						w, next := distance_to(world[attacker.y][attacker.x], world[y][x], &world)
//						fmt.Println(w, next)
						if w > 0 && w < min {
							fmt.Println("attacker",string(attacker.c), attacker,"found victim side:",x,y,"health",h, "at dist:",w, hasAdjecentEnemy, "next:",next)
							min = w
							nextMove = CoordCost{next.x, next.y, next.cost}
						}
					}
				}
				//fmt.Println("Check above")
				checkAttackVector(victim.x, victim.y - 1, victim.HP)
				//fmt.Println("Check left")
				checkAttackVector(victim.x - 1, victim.y, victim.HP)
				//fmt.Println("Check right of")
				checkAttackVector(victim.x + 1, victim.y, victim.HP)
				//fmt.Println("Check underneath")
				checkAttackVector(victim.x, victim.y + 1, victim.HP)
			}

			checkAdjecentEnemy := func(x int, y int, attacker rune) bool {
				if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
					if world[y][x].unit != nil && world[y][x].unit.c != attacker && world[y][x].unit.HP > 0 {
						return true
					}
				}
				return false
			}
			if checkAdjecentEnemy(attacker.x, attacker.y - 1, attacker.c) || checkAdjecentEnemy(attacker.x - 1, attacker.y, attacker.c) || checkAdjecentEnemy(attacker.x + 1, attacker.y, attacker.c) || checkAdjecentEnemy(attacker.x, attacker.y + 1, attacker.c) {
				hasAdjecentEnemy = true
			}
			//hasAdjecentEnemy = false

			if min != math.MaxInt32 && !hasAdjecentEnemy {
				// TODO: This is not correct, gather all moves, an d pick, but this is to test:
				if world[nextMove.y][nextMove.x].unit != nil {
					panic("Already occupied")
				} else {
					fmt.Println("Checking:",string(attacker.c), attacker.x,attacker.y, "moving to",nextMove.x,nextMove.y)
					if world[attacker.y][attacker.x].unit == nil {
						fmt.Println(attacker)
						fmt.Println(world[attacker.y][attacker.x])
						var x rune
						x = 46
						fmt.Println(string(x))
						visualize(&world)
						panic("WHY")
					} else {
						//units[i].x = nextMove.x
						//units[i].y = nextMove.y
						world[nextMove.y][nextMove.x].unit = world[attacker.y][attacker.x].unit
						if world[nextMove.y][nextMove.x].unit == nil {
							panic("TWI")
						}
						////world[nextMove.y][nextMove.x].unit.x = nextMove.x
						////world[nextMove.y][nextMove.x].unit.y = nextMove.y
						fmt.Println("Unsetting:", world[attacker.y][attacker.x].x, world[attacker.y][attacker.x].y)
						world[attacker.y][attacker.x].unit = nil
						fmt.Println("Unsetting2:", world[nextMove.y][nextMove.x].x, world[nextMove.y][nextMove.x].y)
						if world[nextMove.y][nextMove.x].unit == nil {
							panic("TWI2")
						}

						attacker.x = nextMove.x
						attacker.y = nextMove.y

						//world[attacker.y][attacker.x].c = '.'
						if i == 1 {

						}

						// reset the world
						for y := 0; y < len(world); y++ {
							for x := 0; x < len(world[y]); x++ {
								world[y][x].cost = 0
							}
						}
						//fmt.Println(nextMove)
						//visualize(&world)
					}
				}
			}
		}
		// sort units
		sort.Slice(units, func(i, j int) bool {
			if units[i].y == units[j].y {
				return units[i].x < units[j].x
			}
			return units[i].y < units[j].y
		})
		for _, attacker := range units {
			fmt.Println("IN2: attacker x,y", attacker.x, attacker.y, string(attacker.c))
		}


		// sort units
		sort.Slice(units, func(i, j int) bool {
			if units[i].y == units[j].y {
				return units[i].x < units[j].x
			}
			return units[i].y < units[j].y
		})

		fmt.Println(units)



		for i, attacker := range units {
			if attacker.HP <= 0 {
				continue
			}

			minFound := math.MaxInt32
			var unit * Unit
			findWithLowestHealth := func(x int, y int) {
				if !(x >= 0 && y >= 0 && y < len(world) && x < len((world)[y])) {
					return
				}
				if world[y][x].unit != nil {
					if world[y][x].unit.HP <= 0 {
						return
					}
					if world[y][x].unit.c == attacker.c {
						return
					}

					if world[y][x].unit.HP < minFound {
						minFound = world[y][x].unit.HP
						unit = world[y][x].unit
						fmt.Println("Found new target:",string(attacker.c),string(unit.c),"which is at:",unit.x,unit.y)
					} else if unit != nil {
						fmt.Println("Got old? target:",string(attacker.c),string(unit.c),"which is at:",unit.x,unit.y)

					}
				}
			}
			if i == 1 {

			}
			findWithLowestHealth(attacker.x, attacker.y - 1)
			findWithLowestHealth(attacker.x - 1, attacker.y)
			findWithLowestHealth(attacker.x + 1, attacker.y)
			findWithLowestHealth(attacker.x, attacker.y + 1)
			if unit != nil {
				//fmt.Println("looking up",unit.x,unit.y)
				//for j, victim := range units {
				//	if victim.HP <= 0 {
				//		fmt.Println("Skip 1")
				//		continue
				//	}
				//	fmt.Println("compare:",victim.x,victim.y,unit.x,unit.y)
				//	if victim.x == unit.x && victim.y == unit.y {
				//		fmt.Println("attacked!")
						unit.HP -= 3
						if unit.HP < 0 {
							fmt.Println("IT DIEDED!!")
						}
				//		world[unit.y][unit.x].unit = nil
				//		panic("DEath")
				//		break
				//	} else {
				//		fmt.Println("Skip 2")
				//	}
				//}
			}

			fmt.Println(minFound, unit)
		}
		fmt.Println("Res:",units)
		for _, x := range units {
			fmt.Println("UNIT:",x.c,x.HP)
		}
		fmt.Printf("After %d round(s)\n", round + 1)
		visualize(&world)

	}
}

// An TileHeap is a min-heap of ints.
type TileHeap []Tile

func (h TileHeap) Len() int { return len(h) }
func (h TileHeap) Less(i, j int) bool {
	if h[i].cost == h[j].cost {
		if h[i].y == h[j].y {
			return h[i].x < h[j].x
		}
		return h[i].y < h[j].y
	}
	return h[i].cost < h[j].cost
}
func (h TileHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *TileHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Tile))
}

func (h *TileHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}


type CoordCost struct {
	x int
	y int
	cost int
}

func distance_to(from Tile, to Tile, world *[][]Tile) (distance int, next CoordCost) {
	exploredNodes := map[CoordCost]bool{}
	frontier := &TileHeap{}
	heap.Init(frontier)
	from.cost = 0
	from.px = -1
	from.py = -1
	heap.Push(frontier, from)
	debug := false
	if from.x == 7 && from.y == 3 {
		//debug = true
	}
	pred := map[CoordCost]CoordCost{}
	for ; frontier.Len() != 0;  {
		currentp := &(*frontier)[0]
		// COPY LIKE THIS:
		current := Tile{currentp.c, currentp.unit, currentp.x, currentp.y, currentp.cost, currentp.px, currentp.y}
		heap.Pop(frontier)

		//// Solution found
		if current.x == to.x && current.y == to.y {
			//fmt.Println("Solution found:", current.cost)

			lookup := CoordCost{to.x, to.y, current.cost}
			if debug {
				fmt.Println("BEGIN LOOKING UP:", lookup)
			}
			for {
				if _, exists := pred[lookup]; !exists {
					break
				}
				nextlookup := pred[lookup]
				if nextlookup.x == from.x && nextlookup.y == from.y {
					if debug {
						fmt.Println("NEXT:",nextlookup)
					}
					break
				}
				lookup = nextlookup
				if debug {
					fmt.Println("OK:", lookup)
				}
			}
			//fmt.Println()
			if debug {
				fmt.Println("BOO")
				for k, v := range pred {
					fmt.Println("PRED:", k, v)
				}
			}
			return current.cost, lookup
		}

		exploredNodes[CoordCost{current.x, current.y, -1}] = true

		exploreNeighbour := func(x int, y int, cost int) bool {

			// out side of map
			if !(x >= 0 && y >= 0 && y < len(*world) && x < len((*world)[y])) {
				return false
			}

			// neighbour is wall or other unit
			neighbour := (*world)[y][x]
			if neighbour.c == '#' || (neighbour.unit != nil && neighbour.unit != to.unit) {
				//fmt.Printf("OR %c neighbour.unit = %p\n", neighbour.c, neighbour.unit)
				return false
			}

			// already explored
			if _, alreadyExplored := exploredNodes[CoordCost{neighbour.x, neighbour.y, -1}]; alreadyExplored {
				return false
			}


			// find in frontier
			found := false
			var existingNeighbour *Tile
			for _, v := range *frontier {
				if v.x == neighbour.x && v.y == neighbour.y {
					found = true
					existingNeighbour = &v
					break
				}
			}
			if !found {
				neighbour.cost = current.cost + cost
				heap.Push(frontier, Tile{neighbour.c, neighbour.unit, neighbour.x, neighbour.y, neighbour.cost, neighbour.x, neighbour.y})
				if debug {
					fmt.Println("added1:", CoordCost{neighbour.x, neighbour.y, neighbour.cost}, CoordCost{current.x, current.y, current.cost})
				}
				pred[CoordCost{neighbour.x, neighbour.y, neighbour.cost}]=CoordCost{current.x, current.y, current.cost}
			} else if existingNeighbour.cost > current.cost+cost {
				existingNeighbour.cost = current.cost + cost // this doesn't change the right thing
				for k, v := range *frontier { // this does :-]
					if v.x == neighbour.x && v.y == neighbour.y {
						(*frontier)[k].cost = current.cost + cost
						(*frontier)[k].px = neighbour.x
						(*frontier)[k].py = neighbour.y
						pred[CoordCost{v.x, v.y, (*frontier)[k].cost}]=CoordCost{current.x, current.y, current.cost}
						if debug {
							fmt.Println("added2:", CoordCost{v.x, v.y, (*frontier)[k].cost}, CoordCost{current.x, current.y, current.cost})
						}
						break
					}
				}
			}
			return true
		}

		//fmt.Println("Current:",current)
		exploreNeighbour(current.x, current.y-1, 1) // UP
		exploreNeighbour(current.x-1, current.y, 1) // LEFT
		exploreNeighbour(current.x+1, current.y, 1) // RIGHT
		exploreNeighbour(current.x, current.y+1, 1) // DOWN
		//fmt.Println("UP",exploreNeighbour(current.x, current.y-1, 1)) // UP
		//fmt.Println("LEFT",exploreNeighbour(current.x-1, current.y, 1)) // LEFT
		//fmt.Println("RIGHT",exploreNeighbour(current.x+1, current.y, 1)) // RIGHT
		//fmt.Println("DOWN",exploreNeighbour(current.x, current.y+1, 1)) // DOWN

		//exploreNeighbour(current.x-1, current.y-1, 2) // TOPLEFT
		//exploreNeighbour(current.x-1, current.y+1, 2) // TOPRIGHT
		//exploreNeighbour(current.x+1, current.y-1, 2) // BOTTOMLEFT
		//exploreNeighbour(current.x+1, current.y+1, 2) // BOTTOMRIGHT
	}
	return
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

func visualize(world *[][]Tile) {
	for y := 0; y < len(*world); y++ {
		for x := 0; x < len((*world)[y]); x++ {
			if tile := (*world)[y][x]; tile.unit != nil && tile.unit.HP > 0 {
				fmt.Print(string(tile.unit.c))
			} else {
				fmt.Print(string(tile.c))
			}
		}
		fmt.Println()
	}
}

//				// check up
//				x, y := victim.x, victim.y - 1
//				//inRange := [4]int{-1, -1, -1, -1}
//				if world[y][x].c == '.' && world[y][x].unit == nil {
////					(1, 1), (3, 3)
//					var dist *int
//					calculate := func (offsetX int, offsetY int) (*int) {
//						diffX, diffY := victim.x - attacker.x + offsetX, victim.y - attacker.y + offsetY
//						absX, absY := diffX, diffY
//						if diffX < 0 {
//							diffX *= -1
//						}
//						if diffY < 0 {
//							diffY *= -1
//						}
//						if t := world[attacker.x + diffX][attacker.y + diffY]; t.c == '#' || t.unit != nil {
//							return nil // occupied
//						}
//						dist := absX + absY - 1 // start counting from zero
//						fmt.Printf("attacker at %d,%d , victim at %d,%d, dist = %d\n", attacker.x, attacker.y, victim.x, victim.y, dist)
//						return &dist
//					}
//					dist = calculate(0, -1) // up
//					if dist != nil {
//						fmt.Println(*dist)
//					} else {
//						fmt.Println("NULL")
//					}
//					dist = calculate(-1, 0) // left
//					if dist != nil {
//						fmt.Println(*dist)
//					} else {
//						fmt.Println("NULL")
//					}
//					dist = calculate(1, 0) // right
//					if dist != nil {
//						fmt.Println(*dist)
//					} else {
//						fmt.Println("NULL")
//					}
//					dist = calculate(0, 1) // bottom
//					if dist != nil {
//						fmt.Println(*dist)
//					} else {
//						fmt.Println("NULL")
//					}
//					os.Exit(0)
//				}
