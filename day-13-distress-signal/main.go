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

	// Parse into packet pairs

	pairs := parsePairs(input)

	// Sum ordered pairs

	var orderedPairSum int
	for pairIndex, pair := range pairs {
		if packetsOrder(pair[0], pair[1]) == orderOrdered {
			orderedPairSum += pairIndex + 1
		}
	}

	// Output result

	fmt.Printf("Sum of correctly ordererd pair numbers: %d\n", orderedPairSum)

	// Create divider packets

	dividers := parsePairs([]byte("[[2]]\n[[6]]"))

	// Transform pairs & divider pairs to slice of packets

	var packets []*Packet
	for _, pair := range append(pairs, dividers...) {
		for _, packet := range pair {
			packets = append(packets, packet)
		}
	}

	// Sort packets

	sort.Slice(packets, func(a, b int) bool {
		return packetsOrder(packets[a], packets[b]) == orderOrdered
	})

	// Calculate decode key from indexes of divider packets

	decoderKey := 1
	for packetIndex, rootPacket := range packets {
		packetSlice := rootPacket.values()
		if len(packetSlice) != 1 {
			continue
		}

		packet, isPacket := packetSlice[0].(*Packet)
		if !isPacket {
			continue
		}

		values := packet.values()
		if len(values) != 1 {
			continue
		}

		value, isValue := values[0].(Value)
		if !isValue {
			continue
		}

		switch value {
		case 2, 6:
			decoderKey *= (packetIndex + 1)
		}
	}

	// Output result

	fmt.Printf("Decoder key: %d\n", decoderKey)
}

type Order int

const (
	orderNotOrdered Order = iota - 1
	orderSame
	orderOrdered
)

func packetsOrder(leftValues, rightValues Values) Order {
	leftSlice := leftValues.values()
	rightSlice := rightValues.values()

	for len(leftSlice) > 0 && len(rightSlice) > 0 {
		// Pop first values

		left := leftSlice[0]
		leftSlice = leftSlice[1:]
		right := rightSlice[0]
		rightSlice = rightSlice[1:]

		// Return if unequal values otherwise continue loop

		leftValue, isLeftValue := left.(Value)
		rightValue, isRightValue := right.(Value)

		if isLeftValue && isRightValue {
			if leftValue < rightValue {
				return orderOrdered
			}
			if leftValue > rightValue {
				return orderNotOrdered
			}
			continue
		}

		// If at least 1 list convert

		order := packetsOrder(left, right)
		if order != orderSame {
			return order
		}
	}

	if len(leftSlice) < len(rightSlice) {
		return orderOrdered
	}
	if len(leftSlice) > len(rightSlice) {
		return orderNotOrdered
	}
	return orderSame
}

type Pair [2]*Packet

type Values interface {
	values() []Values
}

type Packet struct {
	parent   *Packet
	elements []Values
}

func (p *Packet) add(value Values) {
	p.elements = append(p.elements, value)
}

func (p Packet) values() []Values {
	return p.elements
}

type Value int

func (v Value) values() []Values {
	return []Values{v}
}

func parsePairs(inputBytes []byte) []Pair {
	pairsBytes := bytes.Split(inputBytes, []byte("\n\n"))

	var pairs []Pair

	for _, pairBytes := range pairsBytes {
		packetsBytes := bytes.Split(pairBytes, []byte("\n"))

		var pair Pair

		for packetBytesIndex, packetBytes := range packetsBytes[:2] {
			var value []rune
			var packet *Packet

			// Add value as sum of runes and reset

			addValueToPacket := func() {
				if len(value) < 1 {
					return
				}
				number, _ := strconv.Atoi(string(value))
				packet.add(Value(number))
				value = []rune{}
			}

		packetLoop:
			for _, char := range string(packetBytes) {
				switch char {
				case '[':
					// Add nested packet or root packet if
					// none exists

					if packet == nil {
						packet = &Packet{}
						continue
					}
					newPacket := &Packet{
						parent: packet,
					}
					packet.add(newPacket)
					packet = newPacket
				case ']':
					// Finish packet and move to parent

					addValueToPacket()
					parent := packet.parent
					if parent == nil {
						break packetLoop
					}
					packet = parent
				case ',':
					addValueToPacket()
				default:
					// Assume number char

					value = append(value, char)
				}
			}

			pair[packetBytesIndex] = packet
		}

		pairs = append(pairs, pair)
	}

	return pairs
}
