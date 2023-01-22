package main

import (
	"fmt"
	"os"
)

func main() {
	// Read input file

	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("input.txt not found")
		return
	}

	datastream := []rune(string(input))

	var markerPosition int

	marker := queue[rune]{
		capacity: 4,
	}

	for charIndex, char := range datastream {
		marker.enqueue(char)

		// Ignore check until marker is 4 characters

		if !marker.isFull() {
			continue
		}

		if marker.isUnique() {
			markerPosition = charIndex + 1
			break
		}
	}

	// Output result

	fmt.Printf("Start-of-packet marker position: %d\n", markerPosition)

	// Reset data

	markerPosition = 0

	marker = queue[rune]{
		capacity: 14,
	}

	for charIndex, char := range datastream {
		marker.enqueue(char)

		// Ignore check until marker is 4 characters

		if !marker.isFull() {
			continue
		}

		if marker.isUnique() {
			markerPosition = charIndex + 1
			break
		}
	}

	// Output result

	fmt.Printf("Start-of-message marker position: %d\n", markerPosition)
}

type set[T comparable] map[T]struct{}

type queue[T comparable] struct {
	capacity int
	items    []T
}

func (q queue[T]) isFull() bool {
	return len(q.items) == q.capacity
}

func (q queue[T]) isUnique() bool {
	itemSet := set[T]{}

	for _, item := range q.items {
		itemSet[item] = struct{}{}
	}

	return len(itemSet) == len(q.items)
}

func (q *queue[T]) enqueue(item T) {
	if q.isFull() {
		q.dequeue()
	}

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
