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

	// Parse input into instructions

	lines := bytes.Split(input, []byte("\n"))

	var instructions []cpuInstruction

	for _, line := range lines {
		instructionRaw := bytes.Split(line, []byte(" "))
		name := string(instructionRaw[0])

		var instruction cpuInstruction

		switch name {
		case "noop":
			instruction.cycles = 1
		case "addx":
			instruction.cycles = 2

			value, _ := strconv.Atoi(string(instructionRaw[1]))
			instruction.value = value
		}

		instructions = append(instructions, instruction)
	}

	// Setup state & data

	lastReadCycle := 220
	readCycles := newSet(20, 60, 100, 140, 180, lastReadCycle)

	x := 1
	var cycle int
	var signalStrengths []int

	for _, instruction := range instructions {
		for instructionCycle := 1; instructionCycle <= instruction.cycles; instructionCycle++ {
			cycle++

			if readCycles.has(cycle) {
				signalStrengths = append(signalStrengths, cycle*x)
			}
		}

		x += instruction.value
	}

	// Sum signal strengths

	var signalStrengthSum int
	for _, signalStrength := range signalStrengths {
		signalStrengthSum += signalStrength
	}

	// Output result

	fmt.Printf("Sum of signal strengths: %d\n", signalStrengthSum)

	// Reset state

	spritePadding := 1
	rowLength := 40
	litPixel := []byte("#")
	darkPixel := []byte(".")

	var pixels []byte

	cycle = 0
	x = 1

	for _, instruction := range instructions {
		for instructionCycle := 1; instructionCycle <= instruction.cycles; instructionCycle++ {
			cycle++

			pixelIndex := normaliseRange(cycle, 1, rowLength+1) - 1

			pixel := darkPixel
			difference := int(math.Abs(float64(x) - float64(pixelIndex)))
			if difference-spritePadding <= 0 {
				pixel = litPixel
			}

			pixels = append(pixels, pixel...)
		}

		x += instruction.value
	}

	// Create screen from pixels

	var rows [][]byte
	for index := 0; index < len(pixels); index += rowLength {
		rows = append(rows, pixels[index:index+rowLength])
	}

	// // Output result

	fmt.Printf("CRT Screen:\n")
	for _, row := range rows {
		fmt.Printf("%s\n", row)
	}
}

type cpuInstruction struct {
	cycles, value int
}

type set[T comparable] map[T]struct{}

func newSet[T comparable](values ...T) set[T] {
	s := set[T]{}
	for _, v := range values {
		s[v] = struct{}{}
	}
	return s
}

func (s set[T]) has(element T) bool {
	_, ok := s[element]
	return ok
}

func normaliseRange(value, start, end int) int {
	width := float64(end - start)
	offset := float64(value - start)

	return int(offset - (math.Floor(offset/width) * width) + float64(start))
}
