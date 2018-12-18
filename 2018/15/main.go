package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"sync"
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
}

var (
	wg    sync.WaitGroup
	sides []CoordCost
)

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	//flag.Parse()
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//}

	// Part one
	partOneAnswer, _ := solve("input", 3)
	fmt.Println(partOneAnswer)
	//pprof.StopCPUProfile()

	// Part two
	const maxConcurrentRoutines = 8
	wg = sync.WaitGroup{}
	type Result struct {
		attack int
		answer int
		killed int
	}
	ret := make(chan Result)

	var answer *Result
	minFound := math.MaxInt32
	for check := 4; ; {
		// check batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < maxConcurrentRoutines; i++ {
				res := <-ret
				if res.killed == 0 {
					fmt.Println("All elves survived with attack power:", res.attack)
				}
				if res.killed == 0 && res.attack < minFound {
					minFound = res.attack
					answer = &res
				}
			}
		}()

		// launch batch
		for i := 0; i < maxConcurrentRoutines; i++ {
			wg.Add(1)
			go func(check int) {
				defer wg.Done()
				fmt.Println("Executing game with elves attack power of", check)
				answer, elvesKilled := solve("input", check)
				ret <- Result{check, answer, elvesKilled}
			}(check)
			check++
		}
		wg.Wait()
		if minFound != math.MaxInt32 {
			fmt.Println("Answer of part two:")
			fmt.Println(answer.answer)
			break
		}
	}
}

func solve(filename string, elves_attack_power int) (answer int, elvesKilled int) {

	var (
		world [][]Tile
		units []*Unit
	)

	debug := false
	units = []*Unit{}

	world = [][]Tile{}
	readMap(filename, &world, &units)

	sides = []CoordCost{
		{0, -1, 0}, // up
		{-1, 0, 0}, // left
		{+1, 0, 0}, // right
		{0, +1, 0}, // down
	}

	//	debug = true
	if debug {
		fmt.Println("Initial round:")
		visualize(&world, &units)
	}

	for round := 0; ; round++ {
		sortUnitsByCoord(&units)

		// Move units
		for _, attacker := range units {
			if attacker.HP <= 0 {
				continue
			}

			// Find shortest path to a vulnerable side of a victim
			min := math.MaxInt32
			var nextMove CoordCost
			for _, victim := range units {
				if victim.HP <= 0 || victim.c == attacker.c || victim.x == attacker.x && victim.y == attacker.y {
					continue // cannot target dead unit, same team or self
				}
				for _, side := range sides {
					func(x int, y int, h int) {
						if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
							if world[y][x].c != '.' || world[y][x].unit != nil {
								return
							}
							w, next := calculateDistance(world[attacker.y][attacker.x], world[y][x], &world)
							if w > 0 && w < min {
								min, nextMove = w, CoordCost{next.x, next.y, next.cost}
							}
						}
					}(victim.x+side.x, victim.y+side.y, victim.HP)
				}
			}

			// Determine if we have an enemy adjecent to us
			hasAdjecentEnemy := false
			for _, side := range sides {
				if func(x int, y int, attacker rune) bool {
					if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
						if world[y][x].unit != nil && world[y][x].unit.c != attacker && world[y][x].unit.HP > 0 {
							return true
						}
					}
					return false
				}(attacker.x+side.x, attacker.y+side.y, attacker.c) {
					hasAdjecentEnemy = true
				}
			}

			// If we found a shortest path and are not already engaged in battle..
			if min != math.MaxInt32 && !hasAdjecentEnemy {
				// Sanity check..
				if world[nextMove.y][nextMove.x].unit != nil {
					panic("Already occupied")
				}

				// Move the unit on the map
				world[nextMove.y][nextMove.x].unit = world[attacker.y][attacker.x].unit
				world[attacker.y][attacker.x].unit = nil

				// Update the x,y in the unit as well
				attacker.x = nextMove.x
				attacker.y = nextMove.y
			}

			// For each possible attacker find the first victim with the lowest health
			minFound := math.MaxInt32
			var unit *Unit
			for _, side := range sides {
				func(x int, y int) {
					if !(x >= 0 && y >= 0 && y < len(world) && x < len((world)[y])) {
						return // outside of the map
					}
					if world[y][x].unit != nil {
						if world[y][x].unit.HP <= 0 || world[y][x].unit.c == attacker.c {
							return // unit dead or same team
						}
						if world[y][x].unit.HP < minFound {
							minFound = world[y][x].unit.HP
							unit = world[y][x].unit
						}
					}
				}(attacker.x+side.x, attacker.y+side.y)
			}
			// We found a victim, slay it!!
			if unit != nil {
				if attacker.c == 'E' {
					unit.HP -= elves_attack_power
				} else {
					unit.HP -= 3
				}
				if unit.HP <= 0 {
					if unit.c == 'E' {
						elvesKilled++
						// Small optimization
						if elves_attack_power != 3 {
							return // early exit!
						}
					}
					world[unit.y][unit.x].unit = nil
				}
			}
		}
		if debug {
			fmt.Printf("After %d round(s)\n", round+1)
			visualize(&world, &units)
		}
		goblins, elves := 0, 0
		for _, u := range units {
			if u.HP > 0 && u.c == 'G' {
				goblins++
			} else if u.HP > 0 && u.c == 'E' {
				elves++
			}
		}
		if elves == 0 || goblins == 0 {
			sum := 0
			for _, s := range units {
				if s.HP > 0 {
					sum += s.HP
				}
			}
			if debug {
				fmt.Printf("Inside round %d - cannot continue\n", round)
				visualize(&world, &units)
				fmt.Println(sum)
				fmt.Print(sum * round)
			}
			answer = sum * round
			return
		}
	}
	return
}

func readMap(filename string, world *[][]Tile, units *[]*Unit) {
	y := 0
	forEachLineInFile(filename, func(s string) {
		*world = append(*world, []Tile{})
		x := 0
		for _, c := range s {
			switch c {
			case '#':
				(*world)[y] = append((*world)[y], Tile{c, nil, x, y, 0})
			default:
				var unit *Unit
				if c == 'E' || c == 'G' {
					*units = append(*units, &Unit{c, 200, x, y})
					unit = (*units)[len(*units)-1]
					unit.x, unit.y = x, y
					unit.y = y
				}
				(*world)[y] = append((*world)[y], Tile{'.', unit, x, y, 0})
			}
			x++
		}
		y++
	})
}

func sortUnitsByCoord(units *[]*Unit) {
	// sort units
	sort.Slice(*units, func(i, j int) bool {
		if (*units)[i].y == (*units)[j].y {
			return (*units)[i].x < (*units)[j].x
		}
		return (*units)[i].y < (*units)[j].y
	})
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
	x    int
	y    int
	cost int
}

func calculateDistance(from Tile, to Tile, world *[][]Tile) (distance int, next CoordCost) {
	explored := map[CoordCost]bool{}
	frontier := &TileHeap{}
	heap.Init(frontier)
	heap.Push(frontier, from)
	pred := map[CoordCost]CoordCost{}
	for frontier.Len() != 0 {
		p := &(*frontier)[0]
		current := Tile{p.c, p.unit, p.x, p.y, p.cost}
		heap.Pop(frontier)

		if current.x == to.x && current.y == to.y {
			lookup := CoordCost{to.x, to.y, current.cost}
			for {
				if _, exists := pred[lookup]; !exists {
					break
				}
				nextlookup := pred[lookup]
				if nextlookup.x == from.x && nextlookup.y == from.y {
					break
				}
				lookup = nextlookup
			}
			return current.cost, lookup
		}

		explored[CoordCost{current.x, current.y, -1}] = true

		for _, side := range sides {
			func(x int, y int, cost int) {
				// outside of map
				if !(x >= 0 && y >= 0 && y < len(*world) && x < len((*world)[y])) {
					return
				}
				// neighbour is wall or other unit
				neighbour := (*world)[y][x]
				if neighbour.c == '#' || (neighbour.unit != nil && neighbour.unit != to.unit) {
					return
				}
				// already explored
				if _, explored := explored[CoordCost{neighbour.x, neighbour.y, -1}]; explored {
					return
				}

				// find in frontier
				foundInFrontier := false
				for _, v := range *frontier {
					if v.x == neighbour.x && v.y == neighbour.y {
						foundInFrontier = true
						break
					}
				}

				if !foundInFrontier {
					neighbour.cost = current.cost + cost
					heap.Push(frontier, neighbour)
					neighbourCoordCost := CoordCost{neighbour.x, neighbour.y, neighbour.cost}
					currentCoordCost := CoordCost{current.x, current.y, current.cost}
					pred[neighbourCoordCost] = currentCoordCost
				}
				return
			}(current.x+side.x, current.y+side.y, 1)
		}
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

func visualize(world *[][]Tile, units *[]*Unit) {
	for y := 0; y < len(*world); y++ {
		for x := 0; x < len((*world)[y]); x++ {
			if tile := (*world)[y][x]; tile.unit != nil {
				fmt.Print(string(tile.unit.c))
			} else {
				fmt.Print(string(tile.c))
			}
		}
		fmt.Println()
	}
	for _, u := range *units {
		fmt.Println(u.x, u.y, string(u.c), u.HP)
	}
}
