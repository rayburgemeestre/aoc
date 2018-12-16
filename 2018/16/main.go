package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Instruction struct {
	registers   []int
	instruction []int
	after       []int
}

const (
	addr = iota // 0
	addi = iota // 1
	mulr = iota // 2
	muli = iota // 3
	banr = iota // 4
	bani = iota // 5
	borr = iota // 6
	bori = iota // 7
	setr = iota // 8
	seti = iota // 9
	gtir = iota // 10
	gtri = iota // 11
	gtrr = iota // 12
	eqir = iota // 13
	eqri = iota // 14
	eqrr = iota // 15
)
const (
	A = 0
	B = 1
	C = 2
)

func next(instruction int, input []int, register []int) (out []int) {
	input = input[1:]

	for _, reg := range register {
		out = append(out, reg)
	}

	switch instruction {
	case addr:
		out[input[C]] = register[input[A]] + register[input[B]]
	case addi:
		out[input[C]] = register[input[A]] + input[B]
	case mulr:
		out[input[C]] = register[input[A]] * register[input[B]]
	case muli:
		out[input[C]] = register[input[A]] * input[B]
	case banr:
		out[input[C]] = register[input[A]] & register[input[B]]
	case bani:
		out[input[C]] = register[input[A]] & input[B]
	case borr:
		out[input[C]] = register[input[A]] | register[input[B]]
	case bori:
		out[input[C]] = register[input[A]] | input[B]
	case setr:
		out[input[C]] = register[input[A]]
	case seti:
		out[input[C]] = input[A]
	case gtir:
		if input[A] > register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case gtri:
		if register[input[A]] > input[B] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case gtrr:
		if register[input[A]] > register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case eqir:
		if input[A] == register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case eqri:
		if register[input[A]] == input[B] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case eqrr:
		if register[input[A]] == register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	}
	return
}

func main() {
	//for instruction := addr; instruction <= eqrr; instruction++ {
	//	register := []int{3, 2, 1, 1}
	//	input := []int{11, 2, 3, 3}
	//	expected := []int{3, 2, 1, 0}
	//	out := next(instruction, input, register)
	//	if reflect.DeepEqual(out, expected) {
	//		fmt.Println("Instruction", instruction, "behaves as expected")
	//	}
	//	//fmt.Println("1:", expected)
	//	//fmt.Println("2:", out)
	//}

	foundThreeOrMore, _ := solve(nil)
	fmt.Println("Part one answer:", foundThreeOrMore)

	deduced := map[int]int{}
	for {
		before := len(deduced)
		_, partTwoAnswer := solve(&deduced)
		if len(deduced) <= before {
			fmt.Println("Part two answer:", partTwoAnswer)
			break
		}
	}
}

func solve(deduced *map[int]int) (foundThreeOrMore int, partTwoAnswer int) {
	state := 0
	instr := Instruction{}
	forEachLineInFile("input", func(s string) bool {
		switch state {
		case 0: // Before
			if s == "" {
				// Second blank line detected!
				state = 4 //
				return true
			} else {
				registers := extractNumbers(s)
				instr.registers = registers
			}
			state++
		case 1: // Instruction
			var instruction []int
			for _, f := range strings.Fields(s) {
				i, _ := strconv.Atoi(f)
				instruction = append(instruction, i)
			}
			instr.instruction = instruction
			state++
		case 2: // After
			registers := extractNumbers(s)
			instr.after = registers
			matches, lastInstr, _ := checkInstruction(instr.registers, instr.instruction, instr.after, deduced)
			if matches == 1 && deduced != nil {
				if _, exists := (*deduced)[lastInstr]; !exists {
					(*deduced)[lastInstr] = instr.instruction[0]
					return false // we have to start over again with the knowledge we have know
				}
			} else if matches >= 3 {
				foundThreeOrMore++
			}
			state++
		case 3: // Blank line, repeat.
			state = 0
		case 4: // Another blank line before we start processing the actual instructions
			// Initialize empty registers for processing
			instr.registers = []int{0, 0, 0, 0}
			instr.after = []int{0, 0, 0, 0}
			state++
		case 5: // Process instruction
			var instruction []int
			for _, f := range strings.Fields(s) {
				i, _ := strconv.Atoi(f)
				instruction = append(instruction, i)
			}
			instr.instruction = instruction
			_, _, lastOut := checkInstruction(instr.registers, instr.instruction, instr.after, deduced)
			instr.registers = lastOut
			partTwoAnswer = lastOut[0]
		}
		return true
	})
	return
}

func checkInstruction(register []int, input []int, expected []int, deduced *map[int]int) (matches int, lastInstr int, lastOut []int) {
	for instruction := addr; instruction <= eqrr; instruction++ {
		if deduced != nil {
			if _, exists := (*deduced)[instruction]; exists {
				if (*deduced)[instruction] != input[0] {
					continue
				}

			}
		}
		out := next(instruction, input, register)
		if reflect.DeepEqual(out, expected) {
			lastInstr = instruction
			matches++

		}
		lastOut = out
	}
	return
}

func extractNumbers(s string) []int {
	strNumbers := strings.Split(s, ": ")[1]
	// The next line was missing, and I wasn't checking strconv.Atoi's err below
	// so I always got zero for the first number!
	strNumbers = strings.TrimPrefix(strNumbers, " ")
	strNumbers = strings.Replace(strNumbers, "[", "", 1)
	strNumbers = strings.Replace(strNumbers, "]", "", 1)
	fields := strings.Split(strNumbers, ", ")
	var numbers []int
	for _, num := range fields {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, n)
	}
	return numbers
}

func forEachLineInFile(filename string, callback func(string) bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !callback(scanner.Text()) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
