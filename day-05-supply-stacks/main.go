package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	diagramAnProcedure := bytes.Split(input, []byte("\n\n"))
	procedure := bytes.Split(diagramAnProcedure[1], []byte("\n"))

	// Parse diagram into []stack

	stacks := parseStacks(diagramAnProcedure[0])

	// Parse each move in producure and run it on the stacks

	for _, moveBytes := range procedure {
		move := parseMove(moveBytes)

		for count := 1; count <= move.count; count++ {
			crate, ok := stacks[move.from-1].pop()

			if !ok {
				fmt.Println("Tried to take from an empty stack")
				return
			}

			stacks[move.to-1].push(crate)
		}
	}

	// Get the top crate from each stack

	var topCrates string
	for _, stack := range stacks {
		crate, ok := stack.pop()
		if !ok {
			continue
		}

		topCrates += crate
	}

	// Output result

	fmt.Printf("Top crates: %s\n", topCrates)

	// Reset stacks

	stacks = parseStacks(diagramAnProcedure[0])

	for _, moveBytes := range procedure {
		move := parseMove(moveBytes)

		// Store crates temporarily then
		// insert in reverse order

		var crates []string

		for count := 1; count <= move.count; count++ {
			crate, ok := stacks[move.from-1].pop()

			if !ok {
				fmt.Println("Tried to take from an empty stack")
				return
			}

			crates = append(crates, crate)
		}

		for crateIndex := len(crates) - 1; crateIndex >= 0; crateIndex-- {
			stacks[move.to-1].push(crates[crateIndex])
		}
	}

	// Reset top crates

	topCrates = ""
	for _, stack := range stacks {
		crate, ok := stack.pop()
		if !ok {
			continue
		}

		topCrates += crate
	}

	// Output result

	fmt.Printf("Top crates using new strategy: %s\n", topCrates)
}

// Stack type for storing crates

type stack []string

func (s *stack) push(crate string) {
	*s = append(*s, crate)
}

func (s stack) isEmpty() bool {
	return len(s) == 0
}

func (s *stack) pop() (string, bool) {
	if s.isEmpty() {
		return "", false
	}

	index := len(*s) - 1
	crate := (*s)[index]
	*s = (*s)[:index]

	return crate, true
}

// Serialized move data

type move struct {
	count, from, to int
}

func parseStacks(diagram []byte) []stack {
	lines := bytes.Split(diagram, []byte("\n"))

	// Ignore x axis labels

	rows := lines[:len(lines)-1]

	var stacks []stack

	// Loop backwards over rows to ensure stacks are
	// created correctly

	for rowIndex := len(rows) - 1; rowIndex >= 0; rowIndex-- {
		row := string(rows[rowIndex])

		for charIndex, char := range row {
			if !unicode.IsLetter(char) {
				continue
			}

			// Normalize char index to column index

			columnIndex := (charIndex - 1) / 4

			// Ensure []stack is large enough

			if len(stacks) <= columnIndex {
				stacks = append(stacks, stack{})
			}

			stacks[columnIndex].push(string(char))
		}
	}

	return stacks
}

func parseMove(moveBytes []byte) move {
	words := bytes.Split(moveBytes, []byte(" "))

	// Assume position of operations

	count, _ := strconv.Atoi(string(words[1]))
	from, _ := strconv.Atoi(string(words[3]))
	to, _ := strconv.Atoi(string(words[5]))

	return move{
		count: count,
		from:  from,
		to:    to,
	}
}
