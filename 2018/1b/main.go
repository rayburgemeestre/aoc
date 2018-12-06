package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strconv"

func main() {
    result := 0
    seen := map[int]bool{}
    for {
        file, err := os.Open("input")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            s := scanner.Text()
            i, err := strconv.Atoi(s)
            if err != nil {
                fmt.Println("Could not convert", s, "to number")
                os.Exit(1)
            }
            _, present := seen[result]
            if present {
                fmt.Println("The result", result)
                os.Exit(0)
            }
            seen[result] = true
            result += i
        }
        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }
    }
}
