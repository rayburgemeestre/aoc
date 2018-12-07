package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type pair struct {
	from string
	to string
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tree := map[pair]bool{}
	seen := map[string]bool{}
	letters := make([]string, 0, 1000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmp := strings.Fields(scanner.Text())
		p := pair{tmp[1], tmp[7]}
		tree[p] = true
		letters = append(letters, tmp[1])
		letters = append(letters, tmp[7])
	}

	fmt.Println(tree)

	fmt.Println("Result:")
	for {
		nextLetters := make([]string, 0, 100)
		for i := 0; i < len(letters); i++ {
			letter := letters[i]
			if _, haveSeen := seen[letter]; haveSeen {
				continue
			}

			isReady := true
			for p, _ := range tree {
				if p.to == letter {
					if _, haveSeen := seen[p.from]; !haveSeen {
						isReady = false
					}
				}
			}
			if isReady {
				nextLetters = append(nextLetters, letter)
			}
		}
		sort.Strings(nextLetters)
		if len(nextLetters) > 0 {
			fmt.Printf("%s", nextLetters[0])
			seen[nextLetters[0]] = true
		} else {
			break
		}
	}
	fmt.Println("")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

