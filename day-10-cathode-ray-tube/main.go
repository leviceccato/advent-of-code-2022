package main

import (
	"bytes"
	"fmt"
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
			instruction.duration = 1
		case "addx":
			instruction.duration = 2

			value, _ := strconv.Atoi(string(instructionRaw[1]))
			instruction.value = value
		}

		instructions = append(instructions, instruction)
	}

	// Setup state & data

	lastReadCycle := 220
	readCycles := newSet(20, 60, 100, 140, 180, lastReadCycle)

	cycle := 1
	x := 1
	instructionDuration := instructions[0].duration
	var instructionIndex int
	var signalStrengths []int

	for cycle <= lastReadCycle && instructionIndex < len(instructions)-1 {
		cycle++

		if instructionDuration < 1 {
			instructionIndex++
			instructionDuration = instructions[instructionIndex].duration
		}

		if instructionDuration < 2 {
			x += instructions[instructionIndex].value
		}

		instructionDuration--

		if readCycles.has(cycle) {
			fmt.Println("cycle", cycle, "x", x, "strength", cycle*x)
			signalStrength := cycle * x
			signalStrengths = append(signalStrengths, signalStrength)
		}
	}

	// Sum signal strengths

	var signalStrengthSum int
	for _, signalStrength := range signalStrengths {
		signalStrengthSum += signalStrength
	}

	// Output result

	fmt.Printf("Sum of signal strengths: %d\n", signalStrengthSum)
}

type cpuInstruction struct {
	duration, value int
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
