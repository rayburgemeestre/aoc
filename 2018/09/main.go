package main

import (
	"container/list"
	"fmt"
)

// Test values
//const players = 9
//const lastMarblePartOne = 25
//const players = 10
//const lastMarblePartOne = 1618

// Puzzle values
// 476 players; last marble is worth 71657 points
const players = 476
const lastMarblePartOne = 71657
const lastMarblePartTwo = 71657 * 100

func solution(lastMarble int) {
	marble := 0
	player := 0
	game := list.New()
	current := game.PushBack(marble)
	scores := map[int]int{}

	for marble = 1; marble <= lastMarble; marble++ {
		isDivisibleBy23 := marble%23 == 0

		if isDivisibleBy23 {
			// take current marble as score
			scores[player] += marble

			// move current counter-clockwise seven times
			for i := 0; i < 7; i++ {
				current = current.Prev()
				if current == nil {
					current = game.Back()
				}
			}

			// remove marble and take it as score
			scores[player] += current.Value.(int)
			newCurrent := current.Next()
			game.Remove(current)
			current = newCurrent

		} else {
			// move current clockwise two times
			for i := 0; i < 2; i++ {
				current = current.Next()
				if current == nil {
					current = game.Front()
				}
			}

			// insert marble
			game.InsertBefore(marble, current)
			current = current.Prev()
			if current == nil {
				current = game.Front()
			}
		}

		//fmt.Printf("[%d] %d - ", player, marble)
		//for e := game.Front(); e != nil; e = e.Next() {
		//	// do something with e.Value
		//	fmt.Printf("%d! ", e.Value)
		//}
		//if isDivisibleBy23 {
		//	fmt.Printf("DIVISIBLE BY 23")
		//}
		//fmt.Printf("Current = %d", current.Value)
		//fmt.Println()

		player++
		player %= players
	}

	maxScore := 0
	bestPlayer := 0
	for player, score := range scores {
		if score >= maxScore {
			maxScore = score
			bestPlayer = player
		}
	}

	fmt.Printf("Winner is player %d with score %d!\n", bestPlayer, maxScore)
}

func main() {
	fmt.Println("Part one:")
	solution(lastMarblePartOne)
	fmt.Println("Part two:")
	solution(lastMarblePartTwo)
}
