package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func runeToLowercase(char rune) (rune) {
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

	for {
		var previousChar rune
		var polymerResolved string
		reacted := false

		for _, char := range polymer {
			same := char == previousChar
			a := runeToLowercase(char)
			b := runeToLowercase(previousChar)
			if a == b && !same {
				reacted = true
				previousChar = 0
			} else {
				if previousChar != 0 {
					polymerResolved += string(previousChar)
				}
				previousChar = char
			}
		}
		if previousChar != 0 {
			polymerResolved += string(previousChar)
		}
		polymer, polymerResolved = polymerResolved, ""
		if !reacted {
			break
		}
	}

	fmt.Printf("RESULT: %d\n", len(polymer))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
