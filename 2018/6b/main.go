package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type position struct {
	x int
	y int
}
type position_value struct {
	id   int
	dist int
}

var world map[position]int

func main() {
	world = map[position]int{}

	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	positions := map[position]position_value{}
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	id := 0
	for scanner.Scan() {
		tmp := strings.Split(scanner.Text(), ", ")
		x, _ := strconv.Atoi(tmp[0])
		y, _ := strconv.Atoi(tmp[1])
		pos := position{x, y}
		positions[pos] = position_value{id, 0}
		id++
		if x < minX { minX = x }
		if x > maxX { maxX = x }
		if y < minY { minY = y }
		if y > maxY { maxY = y }
	}

	// just for fun, wrap around the coordinates at least by one..
	minX--
	minY--
	maxX++
	maxY++

	//threshold := 32
	threshold := 10000

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			index := position{x, y}

			for pos, _ := range positions {
				distance := int(math.Abs(float64(pos.x - x))) + int(math.Abs(float64(pos.y - y)))
				world[index] += distance
			}

			//minDist, minID := getIdOfClosest(index)

			//minDist = getMinimalDistance(index, minDist)

			printTile(world[index], threshold)


			//if minDist != -1 {
			//	if x == minX || x == maxX || y == minY || y == maxY {
			//		infinite[minID] = true
			//	}
			//	_, exists := area[minID]
			//	if !exists {
			//		area[minID] = 0
			//	}
			//	area[minID]++
			//}
		}
		fmt.Println("")
	}

	total := 0
	for _, val := range world {
		if val < threshold {
			total++
		}
	}
	fmt.Println(total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func printTile(value int, threshold int) {
	if value < threshold {
		fmt.Printf("X ")
	} else {
		fmt.Printf(". ")
	}
}

func getTotal(value int, threshold int) {
	if value < threshold {
		fmt.Printf("X ")
	} else {
		fmt.Printf(". ")
	}
}

