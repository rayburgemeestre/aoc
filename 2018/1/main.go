package main

import "bufio"
import "fmt"
import "log"
import "os"
import "strconv"
//import "strings"

func main() {
    file, err := os.Open("input")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    result := 0
    for scanner.Scan() {
        //s := strings.TrimSpace(scanner.Text())
        s := scanner.Text()
        i, err := strconv.Atoi(s)
        if err != nil {
            fmt.Println("Could not convert", s, "to number")
            os.Exit(1)
        }
        result += i
    }
    fmt.Println(result)

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
