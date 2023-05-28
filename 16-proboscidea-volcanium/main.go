package main

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	network := parseValves(input)
}

type Valves map[string]Valve

type Valve struct {
	connections []string
	flowRate    int
}

func parseValves(input []byte) Valves {
	numberRegex := regexp.MustCompile(`\d+`)
	valveNameRegex := regexp.MustCompile(`[A-Z]{2,}`)

	valvesBytes := bytes.Split(input, []byte("\n"))

	valves := map[string]Valve{}

	for _, valveBytes := range valvesBytes {
		flowRateBytes := numberRegex.Find(valveBytes)
		valveNamesBytes := valveNameRegex.FindAll(valveBytes, -1)
		flowRate, _ := strconv.Atoi(string(flowRateBytes))

		var currentValveName string
		valve := Valve{
			flowRate: flowRate,
		}

		for i, valveNameBytes := range valveNamesBytes {
			valveName := string(valveNameBytes)

			if i == 0 {
				currentValveName = valveName
				continue
			}

			valve.connections = append(valve.connections, valveName)
		}

		valves[currentValveName] = valve
	}

	return valves
}
