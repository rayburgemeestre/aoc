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
}

var world [][]Tile
var units []*Unit
var sides []CoordCost

func main() {
	units = []*Unit{}

	readMap()

	sides = []CoordCost{
		{0, -1, 0}, // up
		{-1, 0, 0}, // left
		{+1, 0, 0}, // right
		{0, +1, 0}, // down
	}

	fmt.Println("Round Initially:")
	visualize(&world)

	for round := 0; ; round++ {
		sortUnitsByCoord()

		// Move units
		for _, attacker := range units {
			fmt.Println("processing unit:", attacker.x, attacker.y)
			if attacker.HP <= 0 {
				continue
			}

			// Find shortest path to a vulnerable side of a victim
			min := math.MaxInt32
			var nextMove CoordCost
			hasAdjecentEnemy := false
			var foundVictim *Unit
			for _, victim := range units {
				if victim.HP <= 0 || victim.c == attacker.c || victim.x == attacker.x && victim.y == attacker.y {
					continue // cannot target dead unit, same team or self
				}
				checkAttackVector := func(x int, y int, h int) {
					if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
						if world[y][x].c != '.' || world[y][x].unit != nil {
							return
						}
						w, next := calculateDistance(world[attacker.y][attacker.x], world[y][x], &world)
						if w > 0 && w < min {
							min, nextMove = w, CoordCost{next.x, next.y, next.cost}
							foundVictim = victim
						}
					}
				}
				for _, side := range sides {
					checkAttackVector(victim.x+side.x, victim.y+side.y, victim.HP)
				}
			}

			// Determine if we have an enemy adjecent to us
			checkAdjecentEnemy := func(x int, y int, attacker rune) bool {
				if x >= 0 && y >= 0 && y < len(world) && x < len((world)[y]) {
					if world[y][x].unit != nil && world[y][x].unit.c != attacker && world[y][x].unit.HP > 0 {
						return true
					}
				}
				return false
			}
			for _, side := range sides {
				if checkAdjecentEnemy(attacker.x+side.x, attacker.y+side.y, attacker.c) {
					hasAdjecentEnemy = true
				}
			}

			// If we found a shortest path and are not already engaged in battle..
			if min != math.MaxInt32 && !hasAdjecentEnemy {
				// Sanity check..
				if world[nextMove.y][nextMove.x].unit != nil {
					panic("Already occupied")
				}

				fmt.Println(attacker.x, attacker.y, string(attacker.c), "move towards victim", foundVictim.x, foundVictim.y, "next step", nextMove.x, nextMove.y)

				// Move the unit on the map
				world[nextMove.y][nextMove.x].unit = world[attacker.y][attacker.x].unit
				world[attacker.y][attacker.x].unit = nil

				// Update the x,y in the unit as well
				attacker.x = nextMove.x
				attacker.y = nextMove.y
			}

			// sortUnitsByCoord()

			// Calculate damages
			//for _, attacker := range units {
			//    if attacker.HP <= 0 {
			//        continue
			//    }

			// For each possible attacker find the first victim with the lowest health
			minFound := math.MaxInt32
			var unit *Unit
			findWithLowestHealth := func(x int, y int) {
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
			}
			for _, side := range sides {
				findWithLowestHealth(attacker.x+side.x, attacker.y+side.y)
			}
			// We found a victim, slay it!!
			if unit != nil {
				unit.HP -= 3
				if unit.HP <= 0 {
					world[unit.y][unit.x].unit = nil
				}
			}
		}

		fmt.Printf("After %d round(s)\n", round+1)
		visualize(&world)

		numGoblins := 0
		numElves := 0
		for _, u := range units {
			if u.HP > 0 && u.c == 'G' {
				numGoblins++
			} else if u.HP > 0 && u.c == 'E' {
				numElves++
			}
		}
		if numElves == 0 || numGoblins == 0 {
			fmt.Println(round)
			sum := 0
			for _, s := range units {
				if s.HP > 0 {
					sum += s.HP
				}
			}
			fmt.Printf("Inside %d round(s)\n", round)
			visualize(&world)
			fmt.Println(sum)
			fmt.Print(sum * round)
			os.Exit(0)
		}
	}
}

func readMap() {
	world = [][]Tile{}
	y := 0
	forEachLineInFile("input", func(s string) {
		world = append(world, []Tile{})
		x := 0
		for _, c := range s {
			switch c {
			case '#':
				world[y] = append(world[y], Tile{c, nil, x, y, 0})
			default:
				var unit *Unit
				if c == 'E' || c == 'G' {
					units = append(units, &Unit{c, 200, x, y})
					unit = units[len(units)-1]
					unit.x, unit.y = x, y
					unit.y = y
				}
				world[y] = append(world[y], Tile{'.', unit, x, y, 0})
			}
			x++
		}
		y++
	})
}

func sortUnitsByCoord() {
	// sort units
	sort.Slice(units, func(i, j int) bool {
		if units[i].y == units[j].y {
			return units[i].x < units[j].x
		}
		return units[i].y < units[j].y
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
	exploredNodes := map[CoordCost]bool{}
	frontier := &TileHeap{}
	heap.Init(frontier)
	from.cost = 0
	heap.Push(frontier, from)
	pred := map[CoordCost]CoordCost{}
	for frontier.Len() != 0 {
		currentp := &(*frontier)[0]
		current := Tile{currentp.c, currentp.unit, currentp.x, currentp.y, currentp.cost}
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

		exploredNodes[CoordCost{current.x, current.y, -1}] = true

		exploreNeighbour := func(x int, y int, cost int) bool {
			// outside of map
			if !(x >= 0 && y >= 0 && y < len(*world) && x < len((*world)[y])) {
				return false
			}
			// neighbour is wall or other unit
			neighbour := (*world)[y][x]
			if neighbour.c == '#' || (neighbour.unit != nil && neighbour.unit != to.unit) {
				return false
			}
			// already explored
			if _, alreadyExplored := exploredNodes[CoordCost{neighbour.x, neighbour.y, -1}]; alreadyExplored {
				return false
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
				heap.Push(frontier, Tile{neighbour.c, neighbour.unit, neighbour.x, neighbour.y, neighbour.cost})
				pred[CoordCost{neighbour.x, neighbour.y, neighbour.cost}] = CoordCost{current.x, current.y, current.cost}
			}
			return true
		}

		for _, side := range sides {
			exploreNeighbour(current.x+side.x, current.y+side.y, 1)
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

func visualize(world *[][]Tile) {
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
	for _, u := range units {
		fmt.Println(u.x, u.y, string(u.c), u.HP)
	}
}
