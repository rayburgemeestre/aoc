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

type logKey struct {
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

	dataSleepMinutes := make(map[logKey]int)

	var currentGuard int
	var sleepsAt int
	var sleepsUntil int
	re := regexp.MustCompile("#(?P<id>[0-9]+)")

	for _, logLine := range logline {
		minute, _ := strconv.Atoi(logLine[15:17])
		if strings.Contains(logLine, "Guard") {
			match := re.FindStringSubmatch(logLine)
			id, _ := strconv.Atoi(match[1])
			currentGuard = id
		}
		if strings.Contains(logLine, "falls asleep") {
			sleepsAt = minute
			continue
		}
		if strings.Contains(logLine, "wakes up") {
			sleepsUntil = minute
			for min := sleepsAt; min < sleepsUntil; min++ {
				dataSleepMinutes[logKey{currentGuard, min}]++
			}
		}
	}

	var maxMinute int
	var maxKey logKey
	for k, v := range dataSleepMinutes {
		if v > maxMinute {
			maxMinute = v
			maxKey = k
		}
	}
	fmt.Printf("Guard %d sleeps at minute %d for %d times\n", maxKey.guard, maxKey.minute, maxMinute)
	fmt.Printf("Thus the answer: %d\n", maxKey.guard*maxKey.minute)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
