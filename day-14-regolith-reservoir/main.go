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

	caveMap1, sandSourcePoint1 := createCaveMap(input)

	// Simulate sand fall

	var restingSandCount int
	for {
		point := sandSourcePoint1
		var sandState SandState

		for sandState == sandStateFalling {
			point, sandState = caveMap1.moveSand(point)
		}

		if sandState == sandStateLost {
			break
		}

		restingSandCount++
	}

	// Output result

	fmt.Printf("Resting sand count: %d\n", restingSandCount)

	caveMap2, sandSourcePoint2 := createCaveMap(input)

	// Add space above floor

	caveMap2.addPoint(Point{
		x: int(caveMap2.minX),
		y: int(caveMap2.maxY) + 1,
	}, materialAir)

	// Simulate sand fall with floor

	caveMap2.hasFloor = true
	restingSandCount = 0
	for {
		point := sandSourcePoint2
		var sandState SandState

		for sandState != sandStateStopped {
			point, sandState = caveMap2.moveSand(point)
			// fmt.Println(point, sandState)
		}

		restingSandCount++

		if point == sandSourcePoint2 {
			break
		}
	}

	// Output result

	fmt.Printf("Resting sand count with floor: %d\n", restingSandCount)

	// Debug points to materials

	// var grid [][][]byte

	// for y := 0; y <= caveMap2.height(); y++ {
	// 	row := [][]byte{}

	// 	for x := 0; x <= caveMap2.width(); x++ {
	// 		row = append(row, []byte("."))
	// 	}

	// 	grid = append(grid, row)
	// }

	// for point, material := range caveMap2.elements {
	// 	x := point.x - int(caveMap2.minX)
	// 	y := point.y - int(caveMap2.minY)

	// 	switch material {
	// 	case materialSand:
	// 		grid[y][x] = []byte("o")
	// 	case materialRock:
	// 		grid[y][x] = []byte("#")
	// 	}
	// }

	// for _, row := range grid {
	// 	for _, materialBytes := range row {
	// 		fmt.Print(string(materialBytes))
	// 	}
	// 	fmt.Print("\n")
	// }
}

type SandState int

const (
	sandStateFalling SandState = iota
	sandStateStopped
	sandStateLost
)

type Material int

const (
	materialAir Material = iota
	materialRock
	materialSand
)

type Point struct {
	x, y int
}

type MaterialMap struct {
	maxX, maxY, minX, minY float64
	elements               map[Point]Material
	hasFloor               bool
}

func (m MaterialMap) width() int {
	return int(m.maxX - m.minX)
}

func (m MaterialMap) height() int {
	return int(m.maxY - m.minY)
}

// Helper function for adding points with
// setting min/max vars as side effect

func (m *MaterialMap) addPoint(point Point, material Material) {
	xFloat := float64(point.x)
	yFloat := float64(point.y)

	m.maxX = math.Max(float64(m.maxX), xFloat)
	m.maxY = math.Max(float64(m.maxY), yFloat)
	m.minX = math.Min(float64(m.minX), xFloat)
	m.minY = math.Min(float64(m.minY), yFloat)

	m.elements[point] = material
}

func (m *MaterialMap) moveSand(point Point) (Point, SandState) {
	yBelow := point.y + 1

	if float64(yBelow) > m.maxY {
		if m.hasFloor {
			return point, sandStateStopped
		} else {
			return point, sandStateLost
		}
	}

	pointsBelow := []Point{
		{x: point.x, y: yBelow},
		{x: point.x - 1, y: yBelow},
		{x: point.x + 1, y: yBelow},
	}

	for _, pointBelow := range pointsBelow {
		material, ok := m.elements[pointBelow]
		if material == materialAir || !ok {
			m.elements[point] = materialAir
			m.addPoint(pointBelow, materialSand)

			return pointBelow, sandStateFalling
		}
	}

	m.elements[point] = materialSand
	return point, sandStateStopped
}

// Convert two ints to range of values between them (inclusive)
// E.g. 5, -5 -> 5, 4, 3, 2, 1, 0, -1, -2, -3, -4, -5

func rangeToValues(start, end int) []int {
	offset := end - start
	distance := int(math.Abs(float64(offset)))

	var direction int
	if distance != 0 {
		direction = offset / distance
	}

	values := []int{start}
	for current := start + direction; current != end+direction; current += direction {
		values = append(values, current)
	}

	return values
}

func parsePointCoords(pointBytes []byte) (int, int) {
	xAndY := bytes.Split(pointBytes, []byte(","))
	x, _ := strconv.Atoi(string(xAndY[0]))
	y, _ := strconv.Atoi(string(xAndY[1]))

	return x, y
}

func createCaveMap(input []byte) (MaterialMap, Point) {
	// Define map

	materialMap := MaterialMap{
		maxX:     math.Inf(-1),
		maxY:     math.Inf(-1),
		minX:     math.Inf(1),
		minY:     math.Inf(1),
		elements: map[Point]Material{},
	}

	// Parse points

	lines := bytes.Split(input, []byte("\n"))

	for _, line := range lines {
		pointsBytes := bytes.Split(line, []byte(" -> "))

		for pointBytesIndex, pointBytes := range pointsBytes {
			currentX, currentY := parsePointCoords(pointBytes)

			// Get next point bytes if they exist

			var nextPointBytes []byte
			if pointBytesIndex < len(pointsBytes)-1 {
				nextPointBytes = pointsBytes[pointBytesIndex+1]
			}

			nextX := currentX
			nextY := currentY
			if len(nextPointBytes) > 0 {
				nextX, nextY = parsePointCoords(nextPointBytes)
			}

			// Create sequence of points to next if it exists

			xValues := rangeToValues(currentX, nextX)
			yValues := rangeToValues(currentY, nextY)

			// Rocks are never diagonal so we can use two loops here

			for _, x := range xValues {
				materialMap.addPoint(Point{x: x, y: currentY}, materialRock)
			}

			for _, y := range yValues {
				materialMap.addPoint(Point{x: currentX, y: y}, materialRock)
			}
		}
	}

	// Add sand source

	sandSourcePoint := Point{x: 500, y: 0}
	materialMap.addPoint(sandSourcePoint, materialAir)

	// Fill in map with air

	for x := 0; x < materialMap.width(); x++ {
		for y := 0; y < materialMap.height(); y++ {
			point := Point{
				x: int(materialMap.minX) + x,
				y: int(materialMap.minY) + y,
			}

			// There is material at this point already, continue

			_, ok := materialMap.elements[point]
			if ok {
				continue
			}

			materialMap.addPoint(point, materialAir)
		}
	}

	return materialMap, sandSourcePoint
}
