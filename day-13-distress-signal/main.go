package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")

	// Parse into packet pairs

	pairs := parsePairs(input)

	fmt.Println(pairs)
}

type Pair [2]Packet

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
					addValueToPacket()
					parent := packet.parent
					if parent == nil {
						break packetLoop
					}
					packet = parent
				case ',':
					addValueToPacket()
				default:
					value = append(value, char)
				}
			}

			pair[packetBytesIndex] = *packet
		}

		pairs = append(pairs, pair)
	}

	return pairs
}
