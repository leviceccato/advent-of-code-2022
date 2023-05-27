package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/exp/constraints"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	tunnels := parseTunnels(input)

	// Row for analysis
	y := 2_000_000

	impossiblePositions := Set[Position]{}

	for x := tunnels.minX; x <= tunnels.maxX; x++ {
		for _, d := range tunnels.detections {

			// On a beacon, occupied
			if x == d.beaconX && y == d.beaconY {
				continue
			}

			// Outside sensor distance, a beacon can be here
			if d.distance < manhattanDistance(x, y, d.sensorX, d.sensorY) {
				continue
			}

			impossiblePositions.add(Position{x, y})
		}
	}

	fmt.Printf("Positions that cannot contain a beacon: %d\n", len(impossiblePositions))
}

type Position [2]int

type Tunnels struct {
	detections []Detection
	maxX, minX int
}

type Detection struct {
	distance, sensorX, sensorY, beaconX, beaconY int
}

type Set[T comparable] map[T]struct{}

func (s *Set[T]) add(element T) {
	(*s)[element] = struct{}{}
}

func manhattanDistance[T Number](x1, y1, x2, y2 T) T {
	distanceX := abs(x1 - x2)
	distanceY := abs(y1 - y2)

	return distanceX + distanceY
}

func parseTunnels(input []byte) Tunnels {
	// Match all numbers, including negative ones with minus signs.
	numberRegex := regexp.MustCompile(`-?\d+`)

	detectionsBytes := bytes.Split(input, []byte("\n"))

	var detections []Detection
	var maxX, minX int

	for _, detectionBytes := range detectionsBytes {

		numbers := numberRegex.FindAll(detectionBytes, -1)

		// Assume we'll find 4 numbers
		sensorX, _ := strconv.Atoi(string(numbers[0]))
		sensorY, _ := strconv.Atoi(string(numbers[1]))
		beaconX, _ := strconv.Atoi(string(numbers[2]))
		beaconY, _ := strconv.Atoi(string(numbers[3]))

		distance := manhattanDistance(sensorX, sensorY, beaconX, beaconY)

		// Set max and min x values so we know how many positions
		// to check
		maxX = max(maxX, sensorX+distance)
		minX = min(minX, sensorX-distance)
		maxX = max(maxX, beaconX)
		minX = min(minX, beaconX)

		detections = append(detections, Detection{
			distance: distance,
			sensorX:  sensorX,
			sensorY:  sensorY,
			beaconX:  beaconX,
			beaconY:  beaconY,
		})
	}

	return Tunnels{
		detections: detections,
		maxX:       maxX,
		minX:       minX,
	}
}

// Simplify calculations

type Number interface {
	constraints.Integer | constraints.Float
}

func abs[T Number](a T) T {
	if a < T(0) {
		return -a
	}
	return a
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
