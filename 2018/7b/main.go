package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// For input use:
const numWorkers = 5
const minTaskDuration = 60

// For input_test use:
//const numWorkers = 2
//const minTaskDuration = 0

type pair struct {
	from string
	to   string
}

func main() {
	tree := map[pair]bool{}
	treeLetterAdded := map[string]bool{}

	workers := [numWorkers]int{}
	workerLetter := [numWorkers]string{}

	lettersUsed := []string{}
	lettersInProgress := map[string]bool{}
	lettersDone := map[string]bool{}

	forEachLineInFile("input", func(line string) {
		words := strings.Fields(line)
		p := pair{words[1], words[7]}
		tree[p] = true
		add := func(num string) {
			if _, alreadyAdded := treeLetterAdded[num]; !alreadyAdded {
				lettersUsed = append(lettersUsed, num)
			}
			treeLetterAdded[num] = true
		}
		add(p.to)
		add(p.from)
	})
	sort.Strings(lettersUsed) // turns out this is simply A-Z

	time := 0
	for {
		// Verbose output
		fmt.Printf("Second: %d, Done: ", time)
		for letter, _ := range lettersDone {
			fmt.Printf("%s", letter)
		}
		fmt.Printf(", Items todo: %d\n", len(lettersUsed)-len(lettersDone))

		// Process the workers
		for index, _ := range workers {
			if workers[index] > 0 {
				workers[index]--
				// When done, mark the letter as done, this frees up the worker again
				if workers[index] == 0 {
					lettersDone[workerLetter[index]] = true
				}
			}
		}

		// Gather letters that are ready to be processed (all their inputs are done)
		lettersReady := []string{}
		for _, letter := range lettersUsed {
			if _, isWorkedOn := lettersInProgress[letter]; isWorkedOn {
				continue
			}
			if isLetterReady(tree, letter, lettersDone) {
				lettersReady = append(lettersReady, letter)
			}
		}

		// Provide letters to workers that are open for processing
		for index, _ := range workers {
			if workers[index] == 0 && len(lettersReady) > 0 {
				// Chop off one letter from the ready list
				letter := lettersReady[0]
				lettersReady = lettersReady[1:]

				// Assign this workload to the worker
				duration := minTaskDuration + int(letter[0]-'A') + 1
				workers[index] = duration
				workerLetter[index] = letter

				// Make sure we won't assign it again to another worker
				lettersInProgress[letter] = true
			}
		}

		if len(lettersUsed)-len(lettersDone) == 0 {
			fmt.Println("All done! Answer =", time)
			break
		}

		time++
	}

}

// TODO: Further refactoring (all these parameters are ugly)
func isLetterReady(tree map[pair]bool, letter string, lettersDone map[string]bool) bool {
	letterRequirementsAreOk := true
	for p, _ := range tree {
		if p.to == letter {
			if _, letterRequirementIsDone := lettersDone[p.from]; !letterRequirementIsDone {
				letterRequirementsAreOk = false
			}
		}
	}
	return letterRequirementsAreOk
}

func forEachLineInFile(filename string, callback func(string) ) {
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
