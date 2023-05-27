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
			if x == d.beacon[0] && y == d.beacon[1] {
				continue
			}

			position := Position{x, y}

			// Outside sensor distance, a beacon can be here
			if d.distance < d.sensor.manhattanDistance(position) {
				continue
			}

			impossiblePositions.add(position)
		}
	}

	fmt.Printf("Positions that cannot contain a beacon: %d\n", len(impossiblePositions))

	possiblePositions := Set[Position]{}
	maxXandY := 4_000_000
	minXandY := 0

	for _, d := range tunnels.detections {
		perimeter := d.sensor.manhattanPerimeter(d.distance + 1)

	perimeterLoop:
		for position := range perimeter {

			// Outside of bounds
			if (position[0] > maxXandY || position[1] > maxXandY) ||
				(position[0] < minXandY || position[1] < minXandY) {
				continue
			}

			for _, comparedDetections := range tunnels.detections {
				// Don't compare to current detection
				if d.sensor == comparedDetections.sensor {
					continue
				}

				// Position is within another sensor's area, ignore
				distance := position.manhattanDistance(comparedDetections.sensor)
				if distance < (comparedDetections.distance + 1) {
					continue perimeterLoop
				}
			}

			possiblePositions.add(position)
		}
	}

	var tuningFrequency int
	for position := range possiblePositions {
		tuningFrequency = position[0]*4_000_000 + position[1]
		// Immediately break because we assume there is only 1
		// possible position found
		break
	}

	fmt.Printf("Tuning frequency: %d\n", tuningFrequency)
}

type Tunnels struct {
	detections []Detection
	maxX, minX int
}

type Detection struct {
	distance       int
	sensor, beacon Position
}

type Set[T comparable] map[T]struct{}

func (s *Set[T]) add(element T) {
	(*s)[element] = struct{}{}
}

type Position [2]int

func (p Position) manhattanDistance(p2 Position) int {
	distanceX := abs(p[0] - p2[0])
	distanceY := abs(p[1] - p2[1])

	return distanceX + distanceY
}

func (p Position) manhattanPerimeter(distance int) Set[Position] {
	perimeter := Set[Position]{}

	// Add top and bottom points for every column
	for x := p[0] - distance; x <= p[0]+distance; x++ {
		yOffset := abs(p[0]-x) - distance

		perimeter.add(Position{x, p[1] + yOffset})
		perimeter.add(Position{x, p[1] - yOffset})
	}

	return perimeter
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

		sensor := Position{sensorX, sensorY}
		beacon := Position{beaconX, beaconY}

		distance := sensor.manhattanDistance(beacon)

		// Set max and min x values so we know
		// how many positions to check
		maxX = max(maxX, sensorX+distance)
		minX = min(minX, sensorX-distance)
		maxX = max(maxX, beaconX)
		minX = min(minX, beaconX)

		detections = append(detections, Detection{
			distance: distance,
			sensor:   sensor,
			beacon:   beacon,
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
