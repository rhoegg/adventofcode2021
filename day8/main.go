package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Display struct {
	wires WireMapping
	patterns [10]string
	outputValue [4]string
}

type WireMapping struct {
	A, B, C, D, E, F, G string
}

func (d *Display) countEasyOutputDigits() int {
	count := 0
	for _, v := range d.outputValue {
		switch len(v) {
		case 2: // 1
			count++
		case 3: // 7
			count++
		case 4: // 4
			count++
		case 7: // 8
			count++
		}
	}
	return count
}


func (d *Display) identifySignalSegments() {
	// 9 doesn't have e
	var oneSignal, sevenSignal, fourSignal, eightSignal string // unique lengths
//	var twoSignal, fiveSignal, threeSignal string // length 5
//	var zeroSignal, sixSignal, nineSignal string // length 6
	var fiveWires, sixWires []string
	for _, p := range d.patterns {
		switch len(p) {
		case 2:
			oneSignal = p
		case 3:
			sevenSignal = p
		case 4:
			fourSignal = p
		case 5:
			fiveWires = append(fiveWires, p)
		case 6:
			sixWires = append(sixWires, p)
		case 7:
			eightSignal = p
		}
	}

	bWire, eWire, fWire := d.identifyBEFWires()
	aWire, cWire := d.identifyACWires(oneSignal, sevenSignal, fWire)
	dWire := d.identifyDWire(fourSignal, bWire, cWire, fWire)
	gWire := d.identifyGWire(eightSignal, aWire, bWire, cWire, dWire, eWire, fWire)
	d.wires = WireMapping{aWire, bWire, cWire, dWire, eWire, fWire, gWire}

}

func (d *Display) identifyBEFWires() (string, string, string) {
	// b is 6
	// e is 4
	// f is 9
	var b, e, f string
	counts := make(map[rune]int)
	for _, p := range d.patterns {
		for _, r := range p {
			counts[r] += 1
		}
	}
	for r, count := range counts {
		switch count {
		case 4:
			e = string(r)
		case 6:
			b = string(r)
		case 9:
			f = string(r)
		}
	}
	return b, e, f
}

func (d *Display) identifyACWires(oneSignal, sevenSignal, fWire string) (string, string) {
	var aWire, cWire string
	for _, wire := range sevenSignal {
		if strings.ContainsRune(oneSignal, wire) {
			if string(wire) != fWire {
				cWire = string(wire)
			}
		} else {
			aWire = string(wire)
		}
	}
	return aWire, cWire
}

func (d *Display) identifyDWire(fourSignal, bWire, cWire, fWire string) (string) {
	var dWire string
	for _, wire := range fourSignal {
		switch string(wire) {
		case bWire:
		case cWire:
		case fWire:
		default:
			dWire = string(wire)
		}
	}
	return dWire
}

func (d *Display) identifyGWire(eightSignal, aWire, bWire, cWire, dWire, eWire, fWire string) string {
	for _, wire := range eightSignal {
		switch string(wire) {
		case aWire:
		case bWire:
		case cWire:
		case dWire:
		case eWire:
		case fWire:
		default:
			return string(wire)
		}
	}
	return ""
}
func (d *Display) decodePattern(pattern string) int {
	digits := make(map[string]int)
	digits[SortString(d.wires.A + d.wires.B + d.wires.C + d.wires.E + d.wires.F + d.wires.G)] = 0
	digits[SortString(d.wires.C + d.wires.F)] = 1
	digits[SortString(d.wires.A + d.wires.C + d.wires.D + d.wires.E + d.wires.G)] = 2
	digits[SortString(d.wires.A + d.wires.C + d.wires.D + d.wires.F + d.wires.G)] = 3
	digits[SortString(d.wires.B + d.wires.C + d.wires.D + d.wires.F)] = 4
	digits[SortString(d.wires.A + d.wires.B + d.wires.D + d.wires.F + d.wires.G)] = 5
	digits[SortString(d.wires.A + d.wires.B + d.wires.D + d.wires.E + d.wires.F + d.wires.G)] = 6
	digits[SortString(d.wires.A + d.wires.C + d.wires.F)] = 7
	digits[SortString(d.wires.A + d.wires.B + d.wires.C + d.wires.D + d.wires.E + d.wires.F + d.wires.G)] = 8
	digits[SortString(d.wires.A + d.wires.B + d.wires.C + d.wires.D + d.wires.F + d.wires.G)] = 9

	fmt.Println("map: ", digits)
	fmt.Printf("Decoded %s: %s -- %v\n", pattern, SortString(pattern), digits)
	return digits[SortString(pattern)]
}

func (d *Display) decodedOutput() int {
	result := 0
	for i, v := range d.outputValue {
		fmt.Printf("Decoded pattern %s: %d\n", v, d.decodePattern(v))
		result += d.decodePattern(v) * int(math.Pow10(3 - i))
	}
	return result
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	inputLines := strings.Split(string(data), "\n")

	var displays []Display

	for _, inputLine := range inputLines {
		parts := strings.Split(inputLine, " | ")
		patterns := strings.Split(parts[0], " ")
		outputValue := strings.Split(parts[1], " ")
		display := Display{}
		for i := range display.patterns {
			display.patterns[i] = patterns[i]
		}
		for i := range display.outputValue {
			display.outputValue[i] = outputValue[i]
		}
		displays = append(displays, display)
	}

	totalEasyDigits := 0
	totalOutputValues := 0
	for _, d := range displays {
		totalEasyDigits += d.countEasyOutputDigits()
		d.identifySignalSegments()
		totalOutputValues += d.decodedOutput()
	}
	fmt.Printf("Total easy digits is %d\n", totalEasyDigits)
	fmt.Printf("Total output values is %d\n", totalOutputValues)

}