package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	// Define mapping of letters to their actual values

	letterMap := map[string]string{
		"A": "ROCK",
		"B": "PAPER",
		"C": "SCISSORS",
		"X": "ROCK",
		"Y": "PAPER",
		"Z": "SCISSORS",
	}

	// Define what each value scores

	scoreMap := map[string]int{
		"ROCK":     1,
		"PAPER":    2,
		"SCISSORS": 3,
	}

	// Define what value beats what

	defeatMap := map[string]string{
		"ROCK":     "SCISSORS",
		"PAPER":    "ROCK",
		"SCISSORS": "PAPER",
	}

	games := bytes.Split(input, []byte("\n"))

	var score int

	for _, game := range games {

		turns := bytes.Split(game, []byte(" "))

		// Map []byte("A") to "ROCK" etc

		turn := letterMap[string(turns[1])]
		opponentTurn := letterMap[string(turns[0])]

		// Add initial score from chosen value

		score = score + scoreMap[turn]

		// Calculate scores from winning & draw

		switch opponentTurn {
		case defeatMap[turn]:
			score += 6
		case turn:
			score += 3
		}
	}

	// Output result

	fmt.Printf("Total score: %d\n", score)

	// Calculate score again with new method
	// Redefine letterMap to accomodate new rules

	letterMap = map[string]string{
		"A": "ROCK",
		"B": "PAPER",
		"C": "SCISSORS",
	}

	resultScoreMap := map[string]int{
		"X": 0,
		"Y": 3,
		"Z": 6,
	}

	lossMap := map[string]string{
		"ROCK":     "PAPER",
		"PAPER":    "SCISSORS",
		"SCISSORS": "ROCK",
	}

	// Reset score

	score = 0

	for _, game := range games {

		turnAndResult := bytes.Split(game, []byte(" "))

		// Map []byte("A") to "ROCK" etc

		opponentTurn := letterMap[string(turnAndResult[0])]
		result := turnAndResult[1]

		// Add initial score from chosen value

		resultScore := resultScoreMap[string(result)]

		score += resultScore

		// Calculate scores from winning & draw

		switch resultScore {
		case 6:
			score += scoreMap[lossMap[opponentTurn]]
		case 3:
			score += scoreMap[opponentTurn]
		case 0:
			score += scoreMap[defeatMap[opponentTurn]]
		}
	}

	// Output result

	fmt.Printf("Total score using proper method: %d\n", score)
}
