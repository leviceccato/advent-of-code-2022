package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	// Parse input

	monkeys1 := parseMonkeys(input)
	monkeys2 := parseMonkeys(input)

	// Run rounds

	roundCount := 20
	business := startMonkeyBusiness(3, roundCount, monkeys1, 0)

	// Output result

	fmt.Printf("Monkey business level for %d rounds: %d\n", roundCount, business)

	// Rerun rounds

	divisibleBysCommonMultiple := 1
	for _, monkey := range monkeys2 {
		divisibleBysCommonMultiple *= monkey.test.divisibleBy
	}

	roundCount = 10_000
	business = startMonkeyBusiness(0, roundCount, monkeys2, divisibleBysCommonMultiple)

	// Output result
	fmt.Printf("Monkey business level for %d rounds: %d\n", roundCount, business)
}

type monkey struct {
	inspectionCount int
	itemWorryLevels queue[int]
	operation       monkeyOperation
	test            monkeyTest
}

type monkeyOperation struct {
	operand, operator string
}

type monkeyTest struct {
	divisibleBy, trueMonkey, falseMonkey int
}

func parseMonkeys(input []byte) []*monkey {
	var monkeys []*monkey

	for _, paragraph := range bytes.Split(input, []byte("\n\n")) {
		lines := bytes.Split(paragraph, []byte("\n"))

		// Parse starting item worry levels

		levelsBytes := bytes.Split(lines[1], []byte(": "))
		levelsBytesSlice := bytes.Split(levelsBytes[1], []byte(", "))

		var itemWorryLevels queue[int]

		for _, levelsBytes := range levelsBytesSlice {
			level, _ := strconv.Atoi(string(levelsBytes))
			itemWorryLevels.enqueue(level)
		}

		// Parse operation

		operationBytes := bytes.Split(lines[2], []byte(" "))
		operator := operationBytes[len(operationBytes)-2]
		operand := operationBytes[len(operationBytes)-1]

		operation := monkeyOperation{
			operator: string(operator),
			operand:  string(operand),
		}

		// Parse test

		divisibleByBytes := bytes.Split(lines[3], []byte(" "))
		divisibleBy, _ := strconv.Atoi(string(divisibleByBytes[len(divisibleByBytes)-1]))

		trueMonkeyBytes := bytes.Split(lines[4], []byte(" "))
		trueMonkey, _ := strconv.Atoi(string(trueMonkeyBytes[len(trueMonkeyBytes)-1]))

		falseMonkeyBytes := bytes.Split(lines[5], []byte(" "))
		falseMonkey, _ := strconv.Atoi(string(falseMonkeyBytes[len(falseMonkeyBytes)-1]))

		test := monkeyTest{
			divisibleBy: divisibleBy,
			trueMonkey:  trueMonkey,
			falseMonkey: falseMonkey,
		}

		// Add new monkey

		newMonkey := monkey{
			itemWorryLevels: itemWorryLevels,
			operation:       operation,
			test:            test,
		}

		monkeys = append(monkeys, &newMonkey)
	}

	return monkeys
}

func startMonkeyBusiness(relief, roundCount int, monkeys []*monkey, commonMultiple int) int {
	for round := 1; round <= roundCount; round++ {
		for _, monkey := range monkeys {
			worryLevels := monkey.itemWorryLevels.values()
			for _, worryLevel := range worryLevels {
				monkey.itemWorryLevels.dequeue()
				monkey.inspectionCount++

				// Set operand

				var operand int
				switch monkey.operation.operand {
				case "old":
					operand = worryLevel
				default:
					parsedOperand, _ := strconv.Atoi(monkey.operation.operand)
					operand = parsedOperand
				}

				// Calculate new worry level

				var newLevel int
				switch monkey.operation.operator {
				case "+":
					newLevel = worryLevel + operand
				case "*":
					newLevel = worryLevel * operand
				}

				if commonMultiple == 0 {
					newLevel /= relief
				} else {
					newLevel %= commonMultiple
				}

				// Calculate target monkey

				targetMonkeyIndex := monkey.test.falseMonkey
				if newLevel%monkey.test.divisibleBy == 0 {
					targetMonkeyIndex = monkey.test.trueMonkey
				}

				// Pass to target monkey

				targetMonkey := monkeys[targetMonkeyIndex]
				targetMonkey.itemWorryLevels.enqueue(newLevel)
			}
		}
	}

	// Calculate monkey business

	sort.Slice(monkeys, func(a, b int) bool {
		return monkeys[a].inspectionCount > monkeys[b].inspectionCount
	})
	monkeyBusiness := 1
	for _, monkey := range monkeys[:2] {
		monkeyBusiness *= monkey.inspectionCount
	}

	return monkeyBusiness
}

type queue[T comparable] struct {
	items []T
}

func (q queue[T]) values() []T {
	return q.items
}

func (q *queue[T]) enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *queue[T]) dequeue() (T, bool) {
	if len(q.items) == 0 {
		return *new(T), false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}
