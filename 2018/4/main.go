package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type log_key struct {
	guard  int
	minute int
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var logline []string
	for scanner.Scan() {
		logline = append(logline, scanner.Text())
	}

	dataSleepMinutes := make(map[log_key]int)
	dataMostSleep := make(map[int]int)

	currentGuard := -1
	sleepsAt := -1
	sleepsUntil := -1
	re := regexp.MustCompile("#(?P<id>[0-9]+)")

	for _, logLine := range logline {
		minute, _ := strconv.Atoi(logLine[15:17])
		if strings.Contains(logLine, "Guard") {
			match := re.FindStringSubmatch(logLine)
			id, _ := strconv.Atoi(match[1])
			currentGuard = id
		}
		if currentGuard == -1 {
			continue
		}
		if strings.Contains(logLine, "falls asleep") {
			sleepsAt = minute
			continue
		}
		if strings.Contains(logLine, "wakes up") {
			sleepsUntil = minute
			dataMostSleep[currentGuard] += sleepsUntil - sleepsAt
			for min := sleepsAt; min < sleepsUntil; min++ {
				dataSleepMinutes[log_key{currentGuard, min}]++
			}
		}
	}

	max := -1
	theGuard := -1
	for k, v := range dataMostSleep {
		if v > max {
			max = v
			theGuard = k
		}
	}
	maxMinute := -1
	theMinute := -1
	for k, v := range dataSleepMinutes {
		if theGuard != k.guard {
			continue
		}
		if v > maxMinute {
			maxMinute = v
			theMinute = k.minute
		}
	}
	fmt.Printf("Guard %d sleeps at minute %d for %d times\n", theGuard, theMinute, maxMinute)
	fmt.Printf("Thus the answer: %d\n", theGuard*theMinute)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
