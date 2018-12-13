package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type Translation struct {
	leftleft   bool
	left       bool
	plantstate bool
	right      bool
	rightright bool
	nextplant  bool
}

func main() {
	plants := map[int]bool{}
	nextplants := map[int]bool{}
	var translations []Translation
	forEachLineInFile("input", func(s string) {
		if strings.Contains(s, "initial") {
			for pos, c := range s[len("initial state: "):len(s)] {
				plants[pos] = c == '#'
			}
		} else if strings.Contains(s, "=>") {
			fields := strings.Fields(s)
			translations = append(translations, Translation{
				fields[0][0] == '#',
				fields[0][1] == '#',
				fields[0][2] == '#',
				fields[0][3] == '#',
				fields[0][4] == '#',
				fields[2] == "#",
			})
		}
	})
	fmt.Println(plants)
	fmt.Println(translations)

	get_stats := func(plants *map[int]bool) intStats {
		stats := newIntStats()
		for index, plant := range *plants {
			if plant {
				stats.add(index)
			}
		}
		return *stats
	}

	visualize := func(plants *map[int]bool, stats intStats) {
		for i := stats.min; i <= stats.max; i++ {
			if (*plants)[i] == true {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
	}

	// Step zero
	stats := get_stats(&plants)
	visualize(&plants, stats)
	fmt.Println()

	previousNumPlants := 0
	numPlantsAddedSincePreviousGeneration := math.MinInt32
	stableForGenerations := 0
	step := 1
	for step = 1; step <= 100; step++ {
		minPos := stats
		stats.reset()
		for i := minPos.min - 2; i <= minPos.max+2; i++ {
			for _, translation := range translations {
				if patternMatches(&translation, &plants, i) {
					nextplants[i] = translation.nextplant
					if translation.nextplant {
						stats.add(i)
					}
				}
			}
		}
		plants = nextplants
		nextplants = map[int]bool{}
		fmt.Printf("%d - ", step)
		visualize(&plants, stats)
		numPlants := getPlantsValue(&plants)
		fmt.Printf("%d - [%d - %d] - %d ans: %d, added: %d\n", stats.max-stats.min, stats.min, stats.max, step, numPlants, numPlants-previousNumPlants)
		newvalue := numPlants - previousNumPlants
		if newvalue == numPlantsAddedSincePreviousGeneration {
			stableForGenerations++
			fmt.Println(stableForGenerations)
		} else {
			fmt.Println(stableForGenerations)
		}
		numPlantsAddedSincePreviousGeneration = newvalue
		previousNumPlants = numPlants
		if stableForGenerations == 100 {
			break
		}
	}

	numPlants := getPlantsValue(&plants)
	fmt.Printf("Part one answer: %d\n", numPlants)

	remain := 50000000000 - step
	remain *= numPlantsAddedSincePreviousGeneration
	numPlants += remain
	fmt.Printf("Part two answer: %d\n", numPlants)
}

func patternMatches(translation *Translation, plants *map[int]bool, i int) bool {
	flag := translation.plantstate == (*plants)[i] &&
		translation.left == (*plants)[i-1] &&
		translation.leftleft == (*plants)[i-2] &&
		translation.right == (*plants)[i+1] &&
		translation.rightright == (*plants)[i+2]
	return flag
}

func getPlantsValue(plants *map[int]bool) (numPlants int) {
	for k, v := range *plants {
		if v {
			numPlants += k
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

type intStats struct {
	min int
	max int
}

func newIntStats() *intStats {
	return &intStats{min: math.MaxInt32, max: math.MinInt32}
}

func (i *intStats) reset() {
	i.min = math.MaxInt32
	i.max = math.MinInt32
}
func (i *intStats) add(num int) {
	if num < i.min {
		i.min = num
	}
	if num > i.max {
		i.max = num
	}
}
