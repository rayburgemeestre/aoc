package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func runeToLowercase(char rune) rune {
	return char | ('a' - 'A')
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	polymer := ""
	for scanner.Scan() {
		polymer += scanner.Text()
	}

	runes := []rune(polymer)
	n := len(runes)
	minLength := len(polymer)
	minLengthChar := 'x'

	// Create map of all unique chars
	chars := map[rune]bool{}
	for index := 1; index < n; index++ {
		a := runeToLowercase(runes[index])
		chars[a] = true
	}

	// For each char to omit
	for char, _ := range chars {
		runes := []rune(polymer)
		n := len(runes)
		// Keep reducing until done
		for {
			changed := false
			current := 0
			index := 0
			previous := -1
			for index = 0; index < n; index++ {
				a := runeToLowercase(runes[index])
				if a != char {
					if previous == -1 {
						previous = index
					} else {
						same := runes[index] == runes[previous]
						b := runeToLowercase(runes[previous])
						if a == b && !same {
							changed = true
							index++
						} else {
							runes[current] = runes[previous]
							current++
						}
					}
					previous = index
				}
			}
			// Don't forget the last char
			if previous != -1 && previous < n && runes[previous] != char {
				runes[current] = runes[previous]
				current++
			}
			n = current
			if !changed {
				if current < minLength {
					minLength = current
					minLengthChar = char
				}
				break
			}
		}
	}

	fmt.Printf("For char %c the result is: %d\n", minLengthChar, minLength)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
