package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type tile struct {
	x int
	y int
}

func parseWord(word string) (int, int, int, int, string) {
	re := regexp.MustCompile("#(?P<id>[0-9]+) @ (?P<x>[0-9]+),(?P<y>[0-9]+): (?P<width>[0-9]+)x(?P<height>[0-9]+)")
	result := make(map[string]string)
	match := re.FindStringSubmatch(word)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	w, _ := strconv.Atoi(result["width"])
	h, _ := strconv.Atoi(result["height"])
	x, _ := strconv.Atoi(result["x"])
	y, _ := strconv.Atoi(result["y"])
	return w, h, x, y, result["id"]
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	tiles := make(map[tile]int)

	for _, word := range words {
		width, height, x, y, _ := parseWord(word)
		for totalH := 0; totalH <= (y + height + y); totalH++ {
			for totalW := 0; totalW < (x + width + x); totalW++ {
				if totalW >= x && totalH >= y && totalW < (x+width) && totalH < (y+height) {
					t := tile{totalW, totalH}
					_, exists := tiles[t]
					if !exists {
						tiles[t] = 0
					} else {
						tiles[t]++
					}
				}
			}
		}
	}

	// loop again and find the one that only has tiles with zero overlap
	for _, word := range words {
		width, height, x, y, id := parseWord(word)

		allBlocksFreeStanding := true
		for total_h := 0; total_h <= (y + height + y); total_h++ {
			for total_w := 0; total_w < (x + width + x); total_w++ {
				if total_w >= x && total_h >= y && total_w < (x+width) && total_h < (y+height) {
					t := tile{total_w, total_h}
					if tiles[t] != 0 {
						allBlocksFreeStanding = false
					}
				}
			}
		}
		if allBlocksFreeStanding {
			fmt.Printf("Found ID: %s!\n", id)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
