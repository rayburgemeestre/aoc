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

type puzzle struct {
	tree            map[pair]bool
	treeLetterAdded map[string]bool

	lettersUsed       []string
	lettersInProgress map[string]bool
	lettersDone       map[string]bool
}

func (puzzle *puzzle) ingestLineOfInput(line string) {
	words := strings.Fields(line)
	p := pair{words[1], words[7]}
	puzzle.tree[p] = true
	add := func(num string) {
		if _, alreadyAdded := puzzle.treeLetterAdded[num]; !alreadyAdded {
			puzzle.lettersUsed = append(puzzle.lettersUsed, num)
		}
		puzzle.treeLetterAdded[num] = true
	}
	add(p.to)
	add(p.from)
}

type workers struct {
	workerSecondsOfWork [numWorkers]int
	workerLetter        [numWorkers]string
}

func (w *workers) workForOneSecond(puzzle *puzzle) {
	for index, _ := range w.workerSecondsOfWork {
		if w.workerSecondsOfWork[index] > 0 {
			w.workerSecondsOfWork[index]--
			// When done, mark the letter as done, this frees up the worker again
			if w.workerSecondsOfWork[index] == 0 {
				puzzle.lettersDone[w.workerLetter[index]] = true
			}
		}
	}
}

func (w *workers) assignWorkload(index int, letter string, duration int) {
	w.workerSecondsOfWork[index] = duration
	w.workerLetter[index] = letter
}

func main() {
	puzzle := puzzle{map[pair]bool{},
		map[string]bool{},
		[]string{}, map[string]bool{}, map[string]bool{}}
	workers := workers{[numWorkers]int{}, [numWorkers]string{}}

	forEachLineInFile("input", puzzle.ingestLineOfInput)

	sort.Strings(puzzle.lettersUsed) // turns out this is simply A-Z

	for time := 0; ; time++ {
		printProgress(time, &puzzle)

		workers.workForOneSecond(&puzzle)

		lettersReady := getLettersReadyForProcessing(&puzzle)

		// Provide letters to workerSecondsOfWork that are open for processing
		for index := range workers.workerSecondsOfWork {

			workerIsFree := workers.workerSecondsOfWork[index] == 0
			workAvailable := len(lettersReady) > 0

			if workerIsFree && workAvailable {
				var letter string
				letter, lettersReady = pop(lettersReady)

				duration := minTaskDuration + int(letter[0]-'A') + 1
				workers.assignWorkload(index, letter, duration)

				// Make sure we won't assign it again to another worker
				puzzle.lettersInProgress[letter] = true
			}
		}

		if len(puzzle.lettersUsed)-len(puzzle.lettersDone) == 0 {
			fmt.Println("All done! Answer =", time)
			break
		}
	}
}

func getLettersReadyForProcessing(puzzle *puzzle) (lettersReady []string) {
	for _, letter := range puzzle.lettersUsed {
		if _, isWorkedOn := puzzle.lettersInProgress[letter]; isWorkedOn {
			continue
		}
		if isLetterReady(letter, puzzle) {
			lettersReady = append(lettersReady, letter)
		}
	}
	return lettersReady
}

func isLetterReady(letter string, puzzle *puzzle) bool {
	letterRequirementsAreOk := true
	for p, _ := range puzzle.tree {
		if p.to == letter {
			if _, letterRequirementIsDone := puzzle.lettersDone[p.from]; !letterRequirementIsDone {
				letterRequirementsAreOk = false
			}
		}
	}
	return letterRequirementsAreOk
}

func printProgress(time int, puzzle *puzzle) {
	fmt.Printf("Second: %d, Done: ", time)
	for letter := range puzzle.lettersDone {
		fmt.Printf("%s", letter)
	}
	fmt.Printf(", Items todo: %d\n", len(puzzle.lettersUsed)-len(puzzle.lettersDone))
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

func pop(lettersReady []string) (string, []string) {
	return lettersReady[0], lettersReady[1:]
}
