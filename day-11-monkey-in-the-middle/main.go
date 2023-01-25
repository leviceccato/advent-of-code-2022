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

	// Set data

	roundCount := 20
	relief := 3

	// Run rounds

	for round := 1; round <= roundCount; round++ {
		for _, monkey := range monkeys {
			for _, worryLevel := range monkey.itemWorryLevels.values() {
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

				newLevel = newLevel / relief

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

	// Output result

	fmt.Printf("Monkey business level: %d\n", monkeyBusiness)
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
