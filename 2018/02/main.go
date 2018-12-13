package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strings"
import "sort"

func parse_line(sorted_s string) (int, int) {
	previous_char := '-'
	times_seen := 1
	doubles := 0
	triples := 0
	for _, char := range sorted_s {
		if previous_char == char {
			times_seen++
		} else {
			times_seen = 1
		}
		if times_seen == 2 {
			doubles++
		} else if times_seen == 3 {
			doubles--
			triples++
		} else if times_seen == 4 {
			triples--
		}
		previous_char = char
	}
	return doubles, triples
}

func sort_word_alphabetically(s string) string {
	temp := strings.Split(s, "")
	sort.Strings(temp)
	return strings.Join(temp, "")
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	doubles_total := 0
	triples_total := 0
	for scanner.Scan() {
		sorted_s := sort_word_alphabetically(scanner.Text())
		doubles, triples := parse_line(sorted_s)
		if doubles > 0 {
			doubles_total++
		}
		if triples > 0 {
			triples_total++
		}
	}
	fmt.Println(doubles_total, triples_total)
	fmt.Println(doubles_total * triples_total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
