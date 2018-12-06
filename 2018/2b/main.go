package main

import "bufio"
import "fmt"
import "log"
import "os"

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	for _, word := range words {
		for _, word2 := range words {
			diff := 0
			result := ""
			for pos, _ := range word {
				char1 := word[pos]
				char2 := word2[pos]
				if char1 != char2 {
					diff++
				} else {
					result = result + string(char1)
				}
			}
			if diff == 1 {
				fmt.Println(result)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
