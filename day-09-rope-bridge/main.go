package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	// Parse input into commands

	var commands []string

	lines := bytes.Split(input, []byte("\n"))

	for _, line := range lines {
		commandRaw := bytes.Split(line, []byte(" "))
		steps, _ := strconv.Atoi(string(commandRaw[1]))
		direction := string(commandRaw[0])

		for step := 0; step < steps; step++ {
			commands = append(commands, direction)
		}
	}

	// Setup head and tail states

	headPoint := point{}
	tailPoint := point{}

	tailHistory := set[point]{}
	tailHistory.add(tailPoint)

	// Move points based on commands

	for _, command := range commands {
		switch command {
		case "U":
			headPoint.y++
		case "D":
			headPoint.y--
		case "L":
			headPoint.x--
		case "R":
			headPoint.x++
		}

		tailPoint.moveToward(headPoint)
		tailHistory.add(tailPoint)
	}

	// Output result

	fmt.Printf("Unique tail positions 1: %d\n", len(tailHistory))

	// Reset state

	var points []*point
	for len(points) < 10 {
		points = append(points, &point{})
	}

	tailHistory.clear()

	// Move points based on commands with more tail sections

	for _, command := range commands {
		for pointIndex, point := range points {
			// Move head point per commands

			if pointIndex == 0 {
				switch command {
				case "U":
					point.y++
				case "D":
					point.y--
				case "L":
					point.x--
				case "R":
					point.x++
				}

				continue
			}

			// Move to next closest point

			point.moveToward(*points[pointIndex-1])

			// Track last point positions

			if pointIndex == len(points)-1 {
				tailHistory.add(*point)
			}
		}
	}

	// Output result
	fmt.Printf("Unique tail positions 2: %d\n", len(tailHistory))
}

type point struct {
	x, y int
}

func (p *point) moveToward(p2 point) {
	offsetX := p2.x - p.x
	offsetY := p2.y - p.y

	distanceX := int(math.Abs(float64(offsetX)))
	distanceY := int(math.Abs(float64(offsetY)))

	if distanceX > 1 || (distanceX > 0 && distanceY > 1) {
		p.x += 1 * (offsetX / distanceX)
	}

	if distanceY > 1 || (distanceY > 0 && distanceX > 1) {
		p.y += 1 * (offsetY / distanceY)
	}
}

type set[T comparable] map[T]struct{}

func (s *set[T]) clear() {
	for element := range *s {
		delete(*s, element)
	}
}

func (s *set[T]) add(element T) {
	(*s)[element] = struct{}{}
}
