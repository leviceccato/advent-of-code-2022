package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	// Parse input into matrix

	var treeHeights matrix[int]

	rows := bytes.Split(input, []byte("\n"))

	for _, rowBytes := range rows {
		row := string(rowBytes)

		var treeHeightRow []int
		for _, char := range row {
			height, _ := strconv.Atoi(string(char))
			treeHeightRow = append(treeHeightRow, height)
		}

		treeHeights = append(treeHeights, treeHeightRow)
	}

	// Create mirrored matrix to indicate which trees are visible
	// from the sides

	treeVisibilities := newMatrix[bool](treeHeights.getDimensions())

	// Calculate visible trees for each side

	for rowIndex, row := range treeHeights {
		if rowIndex == 0 {
			for columnIndex := range row {
				startPoint := point{x: columnIndex, y: 0}
				endPoint := point{x: columnIndex, y: len(treeHeights) - 1}

				column := treeHeights.getLine(startPoint, endPoint)
				columnReversed := treeHeights.getLine(endPoint, startPoint)

				largestHeight := -1

				for heightIndex, height := range column {
					if height > largestHeight {
						largestHeight = height
						treeVisibilities[heightIndex][columnIndex] = true
					}
				}

				largestHeight = -1

				for reversedHeightIndex, height := range columnReversed {
					if height > largestHeight {
						largestHeight = height
						heightIndex := len(column) - 1 - reversedHeightIndex
						treeVisibilities[heightIndex][columnIndex] = true
					}
				}
			}
		}

		largestHeight := -1

		startPoint := point{x: 0, y: rowIndex}
		endPoint := point{x: len(row) - 1, y: rowIndex}

		rowReversed := treeHeights.getLine(endPoint, startPoint)

		for heightIndex, height := range row {
			if height > largestHeight {
				largestHeight = height
				treeVisibilities[rowIndex][heightIndex] = true
			}
		}

		largestHeight = -1

		for reversedHeightIndex, height := range rowReversed {
			if height > largestHeight {
				largestHeight = height
				heightIndex := len(row) - 1 - reversedHeightIndex
				treeVisibilities[rowIndex][heightIndex] = true
			}
		}
	}

	// Sum visible trees

	var visibleTreeCount int

	for _, row := range treeVisibilities {
		for _, isVisible := range row {
			if isVisible {
				visibleTreeCount++
			}
		}
	}

	// Output result

	fmt.Printf("Visible tree count: %d\n", visibleTreeCount)

	// Calculate highest scenic score for each tree

	var highestPossibleScenicScore int

	for rowIndex, row := range treeHeights {
		for heightIndex, height := range row {
			var viewDistanceTop int
			var viewDistanceBottom int
			var viewDistanceLeft int
			var viewDistanceRight int

			topLine := treeHeights.getLine(
				point{x: heightIndex, y: rowIndex - 1},
				point{x: heightIndex, y: 0},
			)
			bottomLine := treeHeights.getLine(
				point{x: heightIndex, y: rowIndex + 1},
				point{x: heightIndex, y: len(treeHeights) - 1},
			)
			leftLine := treeHeights.getLine(
				point{x: heightIndex - 1, y: rowIndex},
				point{x: 0, y: rowIndex},
			)
			rightLine := treeHeights.getLine(
				point{x: heightIndex + 1, y: rowIndex},
				point{x: len(row) - 1, y: rowIndex},
			)

			for _, lineHeight := range topLine {
				viewDistanceTop++
				if lineHeight >= height {
					break
				}
			}
			for _, lineHeight := range bottomLine {
				viewDistanceBottom++
				if lineHeight >= height {
					break
				}
			}
			for _, lineHeight := range leftLine {
				viewDistanceLeft++
				if lineHeight >= height {
					break
				}
			}
			for _, lineHeight := range rightLine {
				viewDistanceRight++
				if lineHeight >= height {
					break
				}
			}

			scenicScore := viewDistanceTop * viewDistanceBottom * viewDistanceLeft * viewDistanceRight
			highestPossibleScenicScore = int(math.Max(float64(highestPossibleScenicScore), float64(scenicScore)))
		}
	}

	// Output result

	fmt.Printf("Highest possible scenic score: %d\n", highestPossibleScenicScore)
}

// Define point for matrix selection

type point struct {
	x, y int
}

// Define matrix to store trees

type matrix[T comparable] [][]T

// Allow for easy creation of matrix

func newMatrix[T comparable](width, height int) matrix[T] {
	m := make(matrix[T], height)

	for index := 0; index < height; index++ {
		m[index] = make([]T, width)
	}

	return m
}

func (m matrix[T]) getDimensions() (int, int) {
	return len(m), len(m[0])
}

func (m matrix[T]) isPointInside(p point) bool {
	return p.x < len(m[0]) && p.x >= 0 && p.y < len(m) && p.y >= 0
}

// Get a line from a matrix in any order

func (m matrix[T]) getLine(p1, p2 point) []T {
	var line []T

	// Ensure lines are straight

	if p1.x != p2.x && p1.y != p2.y {
		return line
	}

	// Ensure line isn't outside matrix

	if !m.isPointInside(p1) || !m.isPointInside(p2) {
		return line
	}

	if p1.x == p2.x {
		direction := 1

		if math.Signbit(float64(p2.y - p1.y)) {
			direction = -1
		}

		for index := p1.y; index != (p2.y + direction); index += direction {
			line = append(line, m[index][p1.x])
		}
	}

	if p1.y == p2.y {
		direction := 1

		if math.Signbit(float64(p2.x - p1.x)) {
			direction = -1
		}

		for index := p1.x; index != (p2.x + direction); index += direction {
			line = append(line, m[p1.y][index])
		}
	}

	return line
}
