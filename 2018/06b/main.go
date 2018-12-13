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
type positionVal struct {
	id   int
	dist int
}

var world map[position]int

func main() {
	world = map[position]int{}

	threshold := 10000 // use 32 for input_test

	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	positions := map[position]positionVal{}
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
		positions[pos] = positionVal{id, 0}
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

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			index := position{x, y}
			for pos, _ := range positions {
				distance := int(math.Abs(float64(pos.x-x))) + int(math.Abs(float64(pos.y-y)))
				world[index] += distance
			}
			printTile(world[index], threshold)
		}
		fmt.Println("")
	}

	fmt.Println(getNumTilesBelowThreshold(threshold))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getNumTilesBelowThreshold(threshold int) int {
	total := 0
	for _, val := range world {
		if val < threshold {
			total++
		}
	}
	return total
}

func printTile(value int, threshold int) {
	if value < threshold {
		fmt.Printf("X ")
	} else {
		fmt.Printf(". ")
	}
}
