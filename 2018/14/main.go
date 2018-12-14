package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	puzzle := []int{3, 7}
	elf1, elf2 := 0, 1
	puzzleInput := "084601"
	foundAnswer1, foundAnswer2 := false, false
	for {
		// Create new numbers
		sum := puzzle[elf1] + puzzle[elf2]
		added := 1
		if sum >= 10 {
			d, r := sum/10, sum%10
			puzzle = append(puzzle, d)
			puzzle = append(puzzle, r)
			added++
		} else {
			puzzle = append(puzzle, sum)
		}
		// Move elves
		elf1 = (elf1 + 1 + puzzle[elf1]) % len(puzzle)
		elf2 = (elf2 + 1 + puzzle[elf2]) % len(puzzle)
		//visualize(&puzzle, elf1, elf2)

		if !foundAnswer1 {
			if inputNum, _ := strconv.Atoi(puzzleInput); len(puzzle) == inputNum+10+added {
				fmt.Println("Answer part one:")
				ans := puzzle[inputNum : len(puzzle)-added]
				for _, num := range ans {
					fmt.Print(num)
				}
				fmt.Println()
				foundAnswer1 = true
			}
		}

		if !foundAnswer2 {
			sliceToFind := make([]int, 0, len(puzzleInput))
			for _, r := range puzzleInput {
				sliceToFind = append(sliceToFind, int(r)-'0')
			}
			for offset := 0; offset < added; offset++ {
				if len(puzzle)-offset < len(puzzleInput) {
					continue
				}
				i := len(puzzle) - len(puzzleInput) - offset
				j := len(puzzle) - offset
				if slice := puzzle[i:j]; reflect.DeepEqual(slice, sliceToFind) {
					fmt.Println("Answer part two:")
					fmt.Println(len(puzzle) - len(slice) - offset)
					fmt.Println()
					foundAnswer2 = true
				}
			}
		}
		if foundAnswer1 && foundAnswer2 {
			break
		}
	}
}

func visualize(puzzle *[]int, elf1 int, elf2 int) {
	for i, p := range *puzzle {
		if elf1 == i {
			fmt.Print("(")
		} else if elf2 == i {
			fmt.Print("[")
		} else {
			fmt.Print("")
		}
		fmt.Print(p)
		if elf1 == i {
			fmt.Print(") ")
		} else if elf2 == i {
			fmt.Print("] ")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
