package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	// Split into pairs

	pairs := bytes.Split(input, []byte("\n"))

	var containingAssignmentCount int

	for _, pair := range pairs {

		sectionRanges := bytes.Split(pair, []byte(","))

		firstRawRange := sectionRanges[0]
		secondRawRange := sectionRanges[1]

		firstRange := parseRange((string(firstRawRange)))
		secondRange := parseRange((string(secondRawRange)))

		if rangeContains(firstRange, secondRange) {
			containingAssignmentCount++
			continue
		}

		if rangeContains(secondRange, firstRange) {
			containingAssignmentCount++
		}
	}

	// Output result

	fmt.Printf("Containing assignments: %d\n", containingAssignmentCount)

	// Calculate overlapping assignments

	var overlappingAssignmentCount int

	for _, pair := range pairs {
		sectionRanges := bytes.Split(pair, []byte(","))

		firstRawRange := sectionRanges[0]
		secondRawRange := sectionRanges[1]

		firstRange := parseRange((string(firstRawRange)))
		secondRange := parseRange((string(secondRawRange)))

		if rangeOverlaps(firstRange, secondRange) {
			overlappingAssignmentCount++
		}
	}

	// Output result

	fmt.Printf("Overlapping assignments: %d\n", overlappingAssignmentCount)
}

// Convert sectionRange in form of "9-12" to
// to a []int{} in the form []int{9, 12}

func parseRange(sectionRange string) []int {
	startAndEnd := strings.Split(sectionRange, "-")

	start := startAndEnd[0]
	end := startAndEnd[1]

	startInt, _ := strconv.Atoi(string(start))
	endInt, _ := strconv.Atoi(string(end))

	return []int{startInt, endInt}
}

// Calculate if ranges in the form []int{9, 12} contain each other

func rangeContains(sectionRange1, sectionRange2 []int) bool {
	return sectionRange2[0] >= sectionRange1[0] && sectionRange2[1] <= sectionRange1[1]
}

// Calculate if ranges overlap anywhere

func rangeOverlaps(sectionRange1, sectionRange2 []int) bool {
	return sectionRange1[1] >= sectionRange2[0] && sectionRange1[0] <= sectionRange2[1]
}
