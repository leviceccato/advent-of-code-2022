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
