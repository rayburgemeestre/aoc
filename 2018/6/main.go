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

var world map[position]map[int]int

func main() {
	world = map[position]map[int]int{}

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
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	// just for fun, wrap around the coordinates at least by one..
	minX--
	minY--
	maxX++
	maxY++

	infinite := map[int]bool{}
	area := map[int]int{}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			index := position{x, y}

			if _, exists := world[index]; !exists {
				world[index] = map[int]int{}
			}

			for pos, val := range positions {
				distance := int(math.Abs(float64(pos.x-x))) + int(math.Abs(float64(pos.y-y)))
				world[index][val.id] = distance
			}

			minDist, minID := getIdOfClosest(index)

			minDist = getMinimalDistance(index, minDist)

			printTile(minDist, minID)

			if minDist != -1 {
				if x == minX || x == maxX || y == minY || y == maxY {
					infinite[minID] = true
				}
				_, exists := area[minID]
				if !exists {
					area[minID] = 0
				}
				area[minID]++
			}
		}
		fmt.Println("")
	}

	maxArea := 0
	for _, val := range positions {
		_, isInfinite := infinite[val.id]
		if isInfinite {
			continue
		}
		if area[val.id] > maxArea {
			maxArea = area[val.id]
		}
	}
	fmt.Println(maxArea)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func printTile(minDist int, minID int) {
	if minDist == -1 {
		fmt.Printf(".. ")
	} else {
		fmt.Printf("%d%d ", minID, minDist)
	}
}

func getIdOfClosest(index position) (int, int) {
	minDist := math.MaxInt32
	minID := -1
	for id, dist := range world[index] {
		if dist < minDist {
			minDist = dist
			minID = id
		}
	}
	return minDist, minID
}

func getMinimalDistance(index position, minDist int) int {
	seen := 0
	for _, dist := range world[index] {
		if dist == minDist {
			seen++
			if seen > 1 {
				minDist = -1
				break
			}
		}
	}
	return minDist
}
