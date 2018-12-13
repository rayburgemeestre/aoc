package main

import (
	"fmt"
	"math"
)

type Coord struct {
	x int
	y int
}

func main() {
	gridSize := 300
	{
		maxSum, maxCoord := solution1(gridSize)
		fmt.Println("Total power:", maxSum)
		fmt.Println("X, Y coordinate:", maxCoord)
		fmt.Printf("Answer, XY: %d,%d\n", maxCoord.x, maxCoord.y)
	}

	{
		maxSum, maxCoord, maxGridSize := solution2(gridSize)
		fmt.Println("Total power:", maxSum)
		fmt.Println("Grid size:", maxGridSize)
		fmt.Println("X, Y coordinate:", maxCoord)
		fmt.Printf("Answer part 2, XY: %d,%d,%d\n", maxCoord.x, maxCoord.y, maxGridSize)
	}
}

func createMap(gridSize int, input int) (world map[Coord]int) {
	world = map[Coord]int{}
	for y := 1; y <= gridSize; y++ {
		for x := 1; x <= gridSize; x++ {
			value := getFuelValue(x, y, input)
			world[Coord{x, y}] = value
		}
	}
	return
}

// For part 2 I moved to use regular arrays in an attempt to speed things up
// With the hash map it was unbearably slow :-)
func createMap2(gridSize int, input int) (world [301][301]int) {
	world = [301][301]int{}
	for y := 1; y <= gridSize; y++ {
		for x := 1; x <= gridSize; x++ {
			value := getFuelValue(x, y, input)
			world[y][x] = value
		}
	}
	return
}

func solution1(gridSize int) (int, Coord) {
	world := createMap(gridSize, 7400)
	scan := []Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{0, 2},
		{1, 2},
		{2, 2},
	}
	maxSum := math.MinInt32
	maxCoord := Coord{0, 0}
	for y := 1; y <= gridSize-2; y++ {
		for x := 1; x <= gridSize-2; x++ {
			sum := 0
			for _, coord := range scan {
				sum += world[Coord{x + coord.x, y + coord.y}]
			}
			if sum > maxSum {
				maxSum = sum
				maxCoord = Coord{x, y}
			}
		}
	}
	return maxSum, maxCoord
}

func solution2(gridSize int) (int, Coord, int) {
	world := createMap2(gridSize, 7400)
	maxSum := math.MinInt32
	maxCoord := Coord{0, 0}
	maxGridSize := 0
	fmt.Println(gridSize)
	for y := 1; y <= gridSize; y++ {
		for x := 1; x <= gridSize; x++ {
			maxSubGridSize := getMaxGridSize(gridSize, y, x)
			for subGrid := 1; subGrid <= maxSubGridSize; subGrid++ {
				sum := 0
				for vy := 0; vy < subGrid; vy++ {
					for vx := 0; vx < subGrid; vx++ {
						sum += world[y+vy][x+vx]
					}
				}
				if sum > maxSum {
					fmt.Printf("New max %d at %d,%d,%d\n", sum, x, y, subGrid)
					maxSum = sum
					maxCoord.x = x
					maxCoord.y = y
					maxGridSize = subGrid
				}
			}
		}
		fmt.Printf("Row %d of %d\n", y, 300)
	}
	fmt.Println()

	return maxSum, maxCoord, maxGridSize
}

func getMaxGridSize(gridSize int, y int, x int) int {
	maxSubGridSize := gridSize - y + 1
	if gridSize-x < maxSubGridSize {
		maxSubGridSize = gridSize - x + 1
	}
	return maxSubGridSize
}

func getFuelValue(x int, y int, input int) (value int) {
	value = x + 10
	value *= y
	value += input // serial
	value *= x + 10
	// So this is what makes me spent a few extra hours, First I wrote:
	//   for ; value > 1000; value -= 1000 {}
	// But I *should* have typed:
	for ; value >= 1000; value -= 1000 {
	}
	value = value / 100
	// Because in this case 1000 would still be 1000, and 1000 / 100 = 10, and not 0 obviously..
	// So every now and then when there is a value of exactly 1000 in the 300x300 grid, I would
	// get a different value, and this screwed up my end result.
	// A nicer approach is:
	//   value %= 1000
	//   value = value / 100
	// Or as a colleague pointed out:
	//   value /= 100
	//   value %= 10
	value -= 5
	return
}
