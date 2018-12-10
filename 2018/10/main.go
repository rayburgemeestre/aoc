package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	solution("input_test")
	solution("input")
}

type position struct {
	x int
	y int
}

func solution(filename string) {

	var positions []position
	var velocities []position

	forEachLineInFile(filename, func(s string) {
		posX, posY, velX, velY := parseLine(s)
		positions = append(positions, position{posX, posY})
		velocities = append(velocities, position{velX, velY})
	})

	lastSecond := math.MaxInt32
	minWidthFound := math.MaxInt32
	previousUniverse := map[position]bool{}
	for second := 0; ; second++ {
		universe := map[position]bool{}
		// apply velocities
		X := newIntStats()
		Y := newIntStats()
		for i, pos := range positions {
			velocity := velocities[i]
			pos.x += velocity.x
			pos.y += velocity.y
			X.add(pos.x)
			Y.add(pos.y)
			universe[pos] = true
			positions[i] = pos
		}
		if X.max-X.min < minWidthFound {
			minWidthFound = X.max - X.min
			lastSecond = second
		} else if X.max-X.min > minWidthFound {
			minWidthFound = X.max - X.min
			lastSecond = second
			fmt.Println("Part one answer:")
			for y := Y.min; y <= Y.max; y++ {
				for x := X.min; x <= X.max; x++ {
					_, exists := previousUniverse[position{x, y}]
					if exists {
						fmt.Printf("#")
					} else {
						fmt.Printf(".")
					}
				}
				fmt.Println()
			}
			fmt.Println("Part two answer:", lastSecond)
			break
		}
		previousUniverse = universe
	}
}

func parseLine(word string) (int, int, int, int) {
	re := regexp.MustCompile("position=<(?P<posX>[^,]+), (?P<posY>[^>]+)> velocity=<(?P<velX>[^,]+), (?P<velY>[^>]+)>")
	result := make(map[string]string)
	match := re.FindStringSubmatch(word)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	posX, _ := strconv.Atoi(strings.TrimLeft(result["posX"], " "))
	posY, _ := strconv.Atoi(strings.TrimLeft(result["posY"], " "))
	velX, _ := strconv.Atoi(strings.TrimLeft(result["velX"], " "))
	velY, _ := strconv.Atoi(strings.TrimLeft(result["velY"], " "))
	return posX, posY, velX, velY
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

type intStats struct {
	min int
	max int
}

func newIntStats() *intStats {
	return &intStats{min: math.MaxInt32, max: math.MinInt32}
}

func (i *intStats) add(num int) {
	if num < i.min {
		i.min = num
	}
	if num > i.max {
		i.max = num
	}
}
