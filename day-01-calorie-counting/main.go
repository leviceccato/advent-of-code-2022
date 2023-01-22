package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	// Separate into groups of counts

	groups := bytes.Split(input, []byte("\n\n"))

	// Sum counts for each group

	var totals []int
	for _, group := range groups {

		counts := bytes.Split(group, []byte("\n"))
		var total int

		for _, count := range counts {

			countString := string(count)
			countInt, err := strconv.Atoi(countString)
			if err != nil {
				fmt.Printf("Unable to convert string '%s' to int", countString)
				return
			}

			total += countInt
		}

		totals = append(totals, total)
	}

	// Sort from lowers to highest

	sort.Ints(totals)

	// Output result

	highestTotal := totals[len(totals)-1]

	fmt.Printf("Highest total: %d\n", highestTotal)

	// Get top 3 highest totals

	highest3Totals := totals[len(totals)-3:]

	var highest3Sum int
	for _, highest3Total := range highest3Totals {
		highest3Sum += highest3Total
	}

	// Output result

	fmt.Printf("Highest 3 total: %d\n", highest3Sum)
}
