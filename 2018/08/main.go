package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//const filename = "input_test"
const filename = "input"

type tree struct {
	children []tree
	metaData []int
}

func ingest(index int, numbers []int, currentTree *tree) (int) {
	numChildren := numbers[index]
	index++
	numMetadata := numbers[index]
	index++
	for child:=0; child < numChildren; child++ {
		subtree := tree{[]tree{}, []int{}}
		newIndex := ingest(index, numbers, &subtree)
		index += newIndex - index
		currentTree.children = append(currentTree.children, subtree)
	}
	for meta:=0; meta<numMetadata; meta++ {
		currentTree.metaData = append(currentTree.metaData, numbers[index])
		index++
	}
	return index
}

func sumMetaDataRecursively(t *tree) (sum int) {
	for _, md := range t.metaData {
		sum += md
	}
	for _, child := range t.children {
		sum += sumMetaDataRecursively(&child)
	}
	return sum
}

func sumValuesOfNodesRecursively(t *tree) (sum int) {
	if len(t.children) == 0 {
		for _, md := range t.metaData {
			sum += md
		}
		return
	}

	for _, childIndex := range t.metaData {
		if childIndex > 0 && childIndex  <= len(t.children) {
			a := sumValuesOfNodesRecursively(&t.children[childIndex - 1])
			sum +=a
		}
	}
	return sum

}

func main() {
	myTree := tree{[]tree{}, []int{}}
	ingest(0, getFileNumbers(filename), &myTree)

	partOneAnswer := sumMetaDataRecursively(&myTree)
	partTwoAnswer := sumValuesOfNodesRecursively(&myTree)

	fmt.Println("Part one answer:", partOneAnswer)
	fmt.Println("Part two answer:", partTwoAnswer)
}


func getFileNumbers(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var numbers []int
	for scanner.Scan() {
		for _, number := range strings.Fields(scanner.Text()) {
			num, err := strconv.Atoi(number)
			if err != nil {
				panic(err.Error())
			}
			numbers = append(numbers, num)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return numbers
}
