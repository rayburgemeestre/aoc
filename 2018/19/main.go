package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	A = 0
	B = 1
	C = 2
)

type Instruction struct {
	instruction  string
	instructions []int
}

func main() {
	instrPtrPosition := 0

	var instructions []Instruction

	forEachLineInFile("input", func(s string) bool {
		if s[0] == '#' {
			ip, err := strconv.Atoi(s[len("#ip "):])
			instrPtrPosition = ip
			if err != nil {
				panic(err)
			}
			fmt.Println(ip)
			return true
		}
		f := strings.Fields(s)
		instr := f[0]
		a, err := strconv.Atoi(f[1])
		if err != nil {
			panic(err)
		}
		b, err := strconv.Atoi(f[2])
		if err != nil {
			panic(err)
		}
		c, err := strconv.Atoi(f[3])
		if err != nil {
			panic(err)
		}
		fmt.Println(instr, a, b, c)

		instruction := Instruction{
			instruction:  instr,
			instructions: []int{a, b, c},
		}
		instructions = append(instructions, instruction)
		return true
	})

	// part one:
	//register := []int{0, 0, 0, 0, 0, 0}
	// part two:
	register := []int{1, 0, 0, 0, 0, 0}

	prev := -1
	for IP := 0; ; {

		register[instrPtrPosition] = IP
		if IP < 0 || IP >= len(instructions) {
			fmt.Println("IP invalid: ", IP)
			fmt.Println(register)
			break
		}
		instr := instructions[IP]

		if prev != register[0] {
			fmt.Printf("%s[%02d]{%2d %2d %2d} reg: %10d %10d %10d %10d %10d %10d\n", instr.instruction, IP,
				instr.instructions[0],
				instr.instructions[1],
				instr.instructions[2],
				register[0], register[1], register[2], register[3], register[4], register[5])
		}
		prev = register[0]

		_, NewRegister := next(instr.instruction, instrPtrPosition, instr.instructions, register)
		register = NewRegister

		NewRegister[instrPtrPosition]++
		IP = NewRegister[instrPtrPosition]
	}

	// I calculated part 2 by hand, inspecting when register[0] changed in the part 1 example:

	//addi[00]{ 5 16  5} reg:          0          0          0          0          0          0
	//addi[08]{ 2  1  2} reg:          1          1       1025       1025          1          8
	//addi[08]{ 2  1  2} reg:          6          5        205       1025          1          8
	//addi[08]{ 2  1  2} reg:         31         25         41       1025          1          8
	//addi[08]{ 2  1  2} reg:         72         41         25       1025          1          8
	//addi[08]{ 2  1  2} reg:        277        205          5       1025          1          8
	//addi[08]{ 2  1  2} reg:       1302       1025          1       1025          1          8

	// Ran part two the same way, it got stuck after a few lines:

	//addi[00]{ 5 16  5} reg:          1          0          0          0          0          0
	//seti[35]{ 0  4  5} reg:          0          0          0   10551425   10550400         35
	//addi[08]{ 2  1  2} reg:          1          1   10551425   10551425          1          8
	//addi[08]{ 2  1  2} reg:          6          5    2110285   10551425          1          8
	//addi[08]{ 2  1  2} reg:         31         25     422057   10551425          1          8

	// From this I deduced after staring at the output for a looooong time, that there register[0]
	// was the cumulative sum of register[1] and since register[1] * register[2] was always register[3]
	// with no remainder, it turns out it was calculating all possible factor numbers for register[4]
	// So manually edited above table into this, and I had my answer:

	//addi[08]{ 2  1  2} reg:          1          1   10551425   10551425          1          8
	//addi[08]{ 2  1  2} reg:          6          5    2110285   10551425          1          8
	//addi[08]{ 2  1  2} reg:         31         25     422057   10551425          1          8
	//addi[08]{ 2  1  2} reg:     422088     422057         25   10551425          1          8
	//addi[08]{ 2  1  2} reg:    2532373    2110285          5   10551425          1          8
	//addi[08]{ 2  1  2} reg:   13083798   10551425          1   10551425          1          8

}

// Execute a given instruction against the current register state and produce the new registers state (out)
// Note that the input instruction contains the instruction number as well, but since we need to speculate which
// instruction number is which operation, we provide our own as a first parameter.
func next(instruction string, IP int, input []int, register []int) (NewIP int, out []int) {
	NewIP = IP

	// initialize the new state of the register with the current state
	for _, reg := range register {
		out = append(out, reg)
	}

	switch instruction {
	case "addr":
		out[input[C]] = register[input[A]] + register[input[B]]
	case "addi":
		out[input[C]] = register[input[A]] + input[B]
	case "mulr":
		out[input[C]] = register[input[A]] * register[input[B]]
	case "muli":
		out[input[C]] = register[input[A]] * input[B]
	case "banr":
		out[input[C]] = register[input[A]] & register[input[B]]
	case "bani":
		out[input[C]] = register[input[A]] & input[B]
	case "borr":
		out[input[C]] = register[input[A]] | register[input[B]]
	case "bori":
		out[input[C]] = register[input[A]] | input[B]
	case "setr":
		out[input[C]] = register[input[A]]
	case "seti":
		out[input[C]] = input[A]
	case "gtir":
		if input[A] > register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case "gtri":
		if register[input[A]] > input[B] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case "gtrr":
		if register[input[A]] > register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case "eqir":
		if input[A] == register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case "eqri":
		if register[input[A]] == input[B] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	case "eqrr":
		if register[input[A]] == register[input[B]] {
			out[input[C]] = 1
		} else {
			out[input[C]] = 0
		}
	}
	return
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
