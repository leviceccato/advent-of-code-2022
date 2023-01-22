package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	// Group by rucksacks

	rucksacks := bytes.Split(input, []byte("\n"))

	var prioritySum int

	for _, rucksack := range rucksacks {
		// Split rucksack in half and check for common characters

		halfIndex := len(rucksack) / 2
		firstCompartment := string(rucksack[:halfIndex])
		secondCompartment := string(rucksack[halfIndex:])

		var commonChar rune

		for _, char := range firstCompartment {
			if strings.Contains(secondCompartment, string(char)) {
				commonChar = char
			}
		}

		prioritySum += getCharPriority(commonChar)
	}

	// Output result

	fmt.Printf("Priority sum: %d\n", prioritySum)

	// Group rucksacks by elf group

	var group [][]byte
	var groups [][][]byte

	for _, rucksack := range rucksacks {
		group = append(group, rucksack)

		// Create a group once it is 3 rucksacks large

		if len(group) > 2 {
			groups = append(groups, group)
			group = [][]byte{}
		}
	}

	// Reset priority sum

	prioritySum = 0

	for _, group := range groups {

		// Find the common character

		var commonChar rune

		for _, char := range string(group[0]) {
			stringChar := string(char)

			if strings.Contains(string(group[1]), stringChar) &&
				strings.Contains(string(group[2]), stringChar) {
				commonChar = char
			}
		}

		prioritySum += getCharPriority(commonChar)
	}

	// Output result

	fmt.Printf("Priority sum of group item types: %d\n", prioritySum)
}

// Map alphabet letters to numbers, e.g:
// 'a' -> 1, 'b' -> 2, 'z' -> 26
// 'A" -> 27, 'B' -> 28, 'Z" -> 52

func getCharPriority(char rune) int {
	// If not a letter return 0

	if !unicode.IsLetter(char) {
		return 0
	}

	if unicode.IsLower(char) {
		return int(char) - int('a') + 1
	}

	return int(char) - int('A') + 27
}
