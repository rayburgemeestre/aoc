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
		if false {
			fmt.Println(word)
		}
		re := regexp.MustCompile("#(?P<id>[0-9]+) @ (?P<x>[0-9]+),(?P<y>[0-9]+): (?P<width>[0-9]+)x(?P<height>[0-9]+)")
		result := make(map[string]string)
		match := re.FindStringSubmatch(word)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		width, _ := strconv.Atoi(result["width"])
		height, _ := strconv.Atoi(result["height"])
		x, _ := strconv.Atoi(result["x"])
		y, _ := strconv.Atoi(result["y"])

		for total_h := 0; total_h <= (y + height + y); total_h++ {
			for total_w := 0; total_w < (x + width + x); total_w++ {
				if total_w >= x && total_h >= y && total_w < (x+width) && total_h < (y+height) {
					t := tile{total_w, total_h}
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
	overlapping := 0
	for _, v := range tiles {
		if v > 0 {
			overlapping++
		}
	}
	fmt.Printf("Overlapping claims: %d\n", overlapping)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
