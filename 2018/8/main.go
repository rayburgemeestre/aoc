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
	numChildren int
	numMetadata int
	children []tree
	metaData []int
}

func ingest(index int, numbers []int, currentTree *tree) (int) {
	currentTree.numChildren = numbers[index]
	index++
	currentTree.numMetadata = numbers[index]
	index++
	for child:=0; child < currentTree.numChildren; child++ {
		subtree := tree{0, 0, []tree{}, []int{}}
		newIndex := ingest(index, numbers, &subtree)
		index += newIndex - index
		currentTree.children = append(currentTree.children, subtree)
	}
	for meta:=0; meta<currentTree.numMetadata; meta++ {
		currentTree.metaData = append(currentTree.metaData, numbers[index])
		index++
	}
	//fmt.Println(currentTree)
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
	if t.numChildren == 0 {
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
	myTree := tree{0, 0, []tree{}, []int{}}
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
	numbers := []int{}
	for scanner.Scan() {
		for _, number := range strings.Fields(scanner.Text()) {
			num, _ := strconv.Atoi(number)
			numbers = append(numbers, num)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return numbers
}
