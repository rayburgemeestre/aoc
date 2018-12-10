package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	solution("input_test", 2, 0)
	solution("input", 5, 60)
}

func solution(filename string, numWorkers int, minTaskDuration int) {
	puzzle := puzzle{map[edge]bool{}, map[string]bool{}, []string{}, map[string]bool{}, map[string]bool{}}
	workers := workers{make([]int, numWorkers), make([]string, numWorkers)}

	forEachLineInFile(filename, puzzle.ingestLineOfInput)
	sort.Strings(puzzle.lettersUsed) // turns out this is simply A-Z

	for time := 0; ; time++ {
		printProgress(time, &puzzle)
		workers.workForOneSecond(&puzzle)
		lettersReady := getLettersReadyForProcessing(&puzzle)

		// Provide letters to workers that are open for processing
		for workerIndex, secondsTodo := range workers.workers {
			workerIsFree := secondsTodo == 0
			workAvailable := len(lettersReady) > 0

			if workerIsFree && workAvailable {
				var letter string
				letter = pop(&lettersReady)
				duration := minTaskDuration + int(letter[0]-'A') + 1

				workers.assignWorkload(workerIndex, letter, duration)

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

type puzzle struct {
	graph           map[edge]bool
	graphLetterDone map[string]bool

	lettersUsed       []string
	lettersInProgress map[string]bool
	lettersDone       map[string]bool
}

type edge struct {
	from string
	to   string
}

func (puzzle *puzzle) ingestLineOfInput(line string) {
	words := strings.Fields(line)
	edge := edge{words[1], words[7]}
	puzzle.graph[edge] = true
	add := func(num string) {
		if _, alreadyAdded := puzzle.graphLetterDone[num]; !alreadyAdded {
			puzzle.lettersUsed = append(puzzle.lettersUsed, num)
		}
		puzzle.graphLetterDone[num] = true
	}
	add(edge.to)
	add(edge.from)
}

type workers struct {
	workers      []int
	workerLetter []string
}

func (w *workers) workForOneSecond(puzzle *puzzle) {
	for index, _ := range w.workers {
		if w.workers[index] > 0 {
			w.workers[index]--
			// When done, mark the letter as done, this frees up the worker again
			if w.workers[index] == 0 {
				puzzle.lettersDone[w.workerLetter[index]] = true
			}
		}
	}
}

func (w *workers) assignWorkload(index int, letter string, duration int) {
	w.workers[index] = duration
	w.workerLetter[index] = letter
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
	for edge := range puzzle.graph {
		if edge.to == letter {
			if _, letterRequirementIsDone := puzzle.lettersDone[edge.from]; !letterRequirementIsDone {
				return false
			}
		}
	}
	return true
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

func pop(lettersReady *[]string) (letter string) {
	letter = (*lettersReady)[0]
	*lettersReady = (*lettersReady)[1:]
	return
}
