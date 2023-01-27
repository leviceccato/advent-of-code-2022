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

	path := grid.aStar(startNode, endNode)

	// Calculate number of steps taken

	stepCount := len(path) - 1

	// Output result

	fmt.Printf("Least number of steps: %d\n", stepCount)

	// Get all lowest elevation nodes

	nodes := grid.values()
	lowestNodes := []*Node{}
	for _, node := range nodes {
		if node.elevation == 0 {
			lowestNodes = append(lowestNodes, node)
		}
	}

	// Run aStar on all lowest nodes and get the shortest path
	// Warning: this is very expensive

	var shortestStepCount int

	var paths [][]*Node
	for _, node := range lowestNodes {
		path := grid.aStar(node, endNode)
		if len(path) > 0 {
			paths = append(paths, path)
		}
	}
	sort.Slice(paths, func(a, b int) bool {
		return len(paths[a]) < len(paths[b])
	})

	shortestStepCount = len(paths[0]) - 1

	// Output result

	fmt.Printf("Least number of steps for all 0 elevations: %d\n", shortestStepCount)
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

func (g Grid) values() []*Node {
	values := []*Node{}

	for _, row := range g {
		for _, node := range row {
			values = append(values, node)
		}
	}

	return values
}

func (g Grid) getPath(parentMap map[*Node]*Node, node *Node) []*Node {
	path := []*Node{node}
	_, ok := parentMap[node]

	for ok {
		node, ok = parentMap[node]
		if ok {
			path = append([]*Node{node}, path...)
		}
	}

	return path
}

func (g Grid) aStar(startNode, endNode *Node) []*Node {
	parentMap := map[*Node]*Node{}
	scoreMap := map[*Node]float64{}
	startDistanceMap := map[*Node]float64{}
	openNodes := Set[*Node]{}

	getMapValue := func(valueMap map[*Node]float64, node *Node) float64 {
		value, ok := valueMap[node]
		if !ok {
			return math.Inf(1)
		}

		return value
	}

	heuristic := func(node *Node) float64 {
		startDistance := getMapValue(startDistanceMap, node)
		endDistance := float64(distance(node.x, node.y, endNode.x, endNode.y))

		return startDistance + endDistance
	}

	openNodes.add(startNode)
	startDistanceMap[startNode] = 0
	scoreMap[startNode] = heuristic(startNode)

	for len(openNodes) > 0 {
		// Get lowest scoring open node

		openNodesSlice := openNodes.values()
		sort.Slice(openNodesSlice, func(a, b int) bool {
			return getMapValue(scoreMap, openNodesSlice[a]) < getMapValue(scoreMap, openNodesSlice[b])
		})
		selectedNode := openNodesSlice[0]

		// If node is end return path from start

		if selectedNode == endNode {
			return g.getPath(parentMap, selectedNode)
		}

		openNodes.remove(selectedNode)

		neighbours := selectedNode.neighbours(g)
		for _, neighbour := range neighbours {
			// Ensure neighbour exists (for nodes on the edges)

			if neighbour == nil {
				continue
			}

			// Ensure we can traverse to this neighbour

			if !isTraversable(selectedNode.elevation, neighbour.elevation) {
				continue
			}

			tentativeStartDistance := getMapValue(startDistanceMap, selectedNode) + 1

			if tentativeStartDistance < getMapValue(startDistanceMap, neighbour) {
				parentMap[neighbour] = selectedNode
				startDistanceMap[neighbour] = tentativeStartDistance
				scoreMap[neighbour] = tentativeStartDistance + heuristic(neighbour)
				openNodes.add(neighbour)
			}
		}
	}

	// Failed to find path, return empty slice

	return []*Node{}
}

type Node struct {
	x, y      int
	elevation int
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
