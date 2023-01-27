package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	// Parse nodes

	lines := bytes.Split(input, []byte("\n"))

	var grid Grid
	var startNode *Node
	var endNode *Node

	for lineIndex, line := range lines {
		grid = append(grid, []*Node{})

		for charIndex, char := range string(line) {
			node := Node{
				x:         charIndex,
				y:         lineIndex,
				elevation: charToElevation(char),
			}

			switch char {
			case 'S':
				startNode = &node
			case 'E':
				endNode = &node
			}

			grid[lineIndex] = append(grid[lineIndex], &node)
		}
	}

	openNodes := newSet(startNode)
	closedNodes := newSet[*Node]()

	for {
		openNodesSorted := openNodes.values()
		sort.Slice(openNodesSorted, func(a, b int) bool {
			return openNodesSorted[a].score(*endNode) < openNodesSorted[b].score(*endNode)
		})

		selectedNode := openNodesSorted[0]
		openNodes.remove(selectedNode)
		closedNodes.add(selectedNode)

		if selectedNode == endNode {
			break
		}

		neighbours := selectedNode.neighbours(grid)
		for _, neighbour := range neighbours {
			// Ensure neighbour exists (for nodes on the edges)

			if neighbour == nil {
				continue
			}

			if closedNodes.has(neighbour) || !isTraversable(selectedNode.elevation, neighbour.elevation) {
				continue
			}

			isNeighbourOpen := openNodes.has(neighbour)

			if !isNeighbourOpen || neighbour.pathLength() < selectedNode.pathLength()+1 {
				neighbour.distanceFromStart = neighbour.distanceTo(*startNode)
				neighbour.parent = selectedNode
				if !isNeighbourOpen {
					openNodes.add(neighbour)
				}
			}
		}
	}

	// Calculate number of steps taken

	var test [][]int
	for i := 1; i <= len(grid); i++ {
		var test2 []int
		for j := 1; j <= len(grid[0]); j++ {
			test2 = append(test2, 0)
		}
		test = append(test, test2)
	}

	var stepCount int
	selectedNode := endNode
	for selectedNode.parent != nil {
		test[selectedNode.y][selectedNode.x] = stepCount + 1
		stepCount++
		selectedNode = selectedNode.parent
	}

	var output []byte
	for _, a := range test {
		for _, b := range a {
			output = append(output, []byte(fmt.Sprintf("%*d ", 3, b))...)
		}
		output = append(output, []byte("\n")...)
	}

	// Output result

	fmt.Printf("Least number of steps: %d\n", stepCount)
	os.WriteFile("output.txt", output, 0644)
}

type Grid [][]*Node

func (g Grid) getNode(x, y int) *Node {
	if y > len(g)-1 || math.Signbit(float64(x)) || math.Signbit(float64(y)) {
		return nil
	}

	row := g[y]
	if x > len(row)-1 {
		return nil
	}

	return row[x]
}

type Node struct {
	parent    *Node
	elevation int
	x, y      int

	distanceFromStart int
}

func (n Node) distanceTo(n2 Node) int {
	return distance(n.x, n.y, n2.x, n2.y)
}

func (n *Node) score(n2 Node) int {
	return int(n.distanceFromStart) + int(n.distanceTo(n2))
}

func (n Node) pathLength() int {
	var length int
	current := &n

	for current.parent != nil {
		length++
		current = current.parent
	}

	return length
}

func (n *Node) neighbours(grid Grid) []*Node {
	return []*Node{
		grid.getNode(n.x-1, n.y),
		grid.getNode(n.x, n.y-1),
		grid.getNode(n.x+1, n.y),
		grid.getNode(n.x, n.y+1),
	}
}

type Set[T comparable] map[T]struct{}

func newSet[T comparable](values ...T) Set[T] {
	s := Set[T]{}

	for _, v := range values {
		s[v] = struct{}{}
	}

	return s
}

func (s Set[T]) values() []T {
	var elements []T

	for key := range s {
		elements = append(elements, key)
	}

	return elements
}

func (s *Set[T]) add(element T) {
	(*s)[element] = struct{}{}
}

func (s *Set[T]) remove(element T) {
	delete(*s, element)
}

func (s Set[T]) has(element T) bool {
	_, ok := s[element]
	return ok
}

func charToElevation(char rune) int {
	switch char {
	case 'S':
		char = 'a'
	case 'E':
		char = 'z'
	}

	return int(char) - int('a')
}

func isTraversable(elevation1, elevation2 int) bool {
	return elevation2 <= elevation1+1
}

func distance(x1, y1, x2, y2 int) int {
	return int(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}
