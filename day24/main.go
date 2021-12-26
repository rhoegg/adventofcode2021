package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data1, err := os.ReadFile("input1.txt")
	if err != nil {
		panic(err)
	}
	data2, err := os.ReadFile("input2.txt")
	if err != nil {
		panic(err)
	}
	data3, err := os.ReadFile("input3.txt")
	if err != nil {
		panic(err)
	}
	data4, err := os.ReadFile("input4.txt")
	if err != nil {
		panic(err)
	}

	alu := NewALU()
	//alu.debug = true
	chunk1ALU := alu.Clone()
	chunk2ALU := chunk1ALU.Clone()
	chunk3ALU := chunk2ALU.Clone()
	//chunk2ALU.debug = true
	alu.Load(strings.Split(string(data1), "\n"))
	chunk1ALU.Load(strings.Split(string(data2), "\n"))
	chunk2ALU.Load(strings.Split(string(data3), "\n"))
	chunk3ALU.Load(strings.Split(string(data4), "\n"))
	for chunk1 := int64(17151); chunk1 < 29999; chunk1++ {
		// memoize forward at inp! save copies of Registers map
		input1 := strconv.FormatInt(chunk1, 10)
		if ! strings.ContainsRune(input1, '0') {
			alu.ResetRegisters()
			alu.Evaluate(input1)
			fmt.Println()
			for chunk2 := int64(11); chunk2 < 99; chunk2++ {
				input2 := strconv.FormatInt(chunk2, 10)
				if ! strings.ContainsRune(input2, '0') {
					chunk1ALU.Registers = alu.CopyRegisters()
					chunk1ALU.Evaluate(input2)
					fmt.Printf("chunk2 %d\n", chunk2)
					for chunk3 := int64(111); chunk3 < 999; chunk3++ {
						input3 := strconv.FormatInt(chunk3, 10)
						if ! strings.ContainsRune(input3, '0') {
							chunk2ALU.Registers = chunk1ALU.CopyRegisters()
							chunk2ALU.Evaluate(input3)
							if chunk2ALU.cacheHit {
								fmt.Print("c")
							} else {
								fmt.Print(".")
							}
							for lastChunk := int64(1111); lastChunk < 9999; lastChunk++ {
								lastInput := strconv.FormatInt(lastChunk, 10)
								if ! strings.ContainsRune(lastInput, '0') && lastInput[3] == '8' {
									chunk3ALU.Registers = chunk2ALU.CopyRegisters()
									chunk3ALU.Evaluate(lastInput)
									if chunk3ALU.Registers['z'] < 430 {
										fmt.Printf("%d %d %d %d %v\n", chunk1, chunk2, chunk3, lastChunk, chunk3ALU)
									}
									if chunk3ALU.Registers['z'] == 0 {
										fmt.Printf("Found a good one: %d %d %d %d \n", chunk1, chunk2, chunk3, lastChunk)
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Printf("ALU: %v\n", alu)
}

type Instruction func(*ALU)

type ALU struct {
	Registers map[rune]int64
	Program []Instruction
	inputBuffer string
	debug bool
	cacheHit bool
	cache map[string]map[rune]int64
}
func NewALU() *ALU {
	return &ALU{
		Registers: make(map[rune]int64),
		cache: make(map[string]map[rune]int64),
	}
}

func (alu ALU) String() string {
	return fmt.Sprintf("(%d, %d, %d, %d)", alu.Registers['w'], alu.Registers['x'], alu.Registers['y'], alu.Registers['z'])
}

func (alu ALU) Clone() *ALU {
	newALU := &alu
	newALU.cache = make(map[string]map[rune]int64)
	return newALU
}

func (alu *ALU) ResetRegisters() {
	alu.Registers = make(map[rune]int64)
}

func (alu *ALU) CopyRegisters() map[rune] int64{
	newRegisters := make(map[rune]int64)
	s := ""
	for r, v := range alu.Registers {
		newRegisters[r] = v
		s += string(r)
	}
	return newRegisters
}

func (alu *ALU) ResetProgram() {
	alu.Program = nil
}

func (alu *ALU) Load(instructions []string) {
	optimized := alu.optimize(instructions)
	if alu.debug {
		fmt.Printf("program(%d) optimized to %d lines\n", len(instructions), len(optimized))
		fmt.Println(strings.Join(optimized, "\n"))
	}
	for _, instructionText := range optimized {
		tokens := strings.Split(instructionText, " ")
		var instruction Instruction
		switch tokens[0] {
		case "inp":
			instruction = alu.ParseInput(tokens[1:])
		case "mul":
			instruction = alu.ParseMultiply(tokens[1:])
		case "eql":
			instruction = alu.ParseEquals(tokens[1:])
		case "add":
			instruction = alu.ParseAdd(tokens[1:])
		case "mod":
			instruction = alu.ParseModulo(tokens[1:])
		case "div":
			instruction = alu.ParseDivide(tokens[1:])
		case "mov":
			instruction = alu.ParseMove(tokens[1:])
		case "neq":
			instruction = alu.ParseNotEquals(tokens[1:])
		case "adm":
			instruction = alu.ParseAddMultiply(tokens[1:])
		case "mad":
			instruction = alu.ParseMultiplyAdd(tokens[1:])
		case "mda":
			instruction = alu.ParseModuloThenAdd(tokens[1:])
		case "mdaneq":
			instruction = alu.ParseModuloThenAddIsNotEqual(tokens[1:])
		case "st1":
			instruction = alu.ParseStrangeAlgo1(tokens[1:])
		case "st2":
			instruction = alu.ParseStrangeAlgo2(tokens[1:])
		default:
			panic("unknown instruction command " + tokens[0])
		}
		alu.Program = append(alu.Program, instruction)
	}
}

func (alu *ALU) optimize(instructions []string) []string {
	instructions = alu.optimizeDivideByOne(instructions)
	instructions = alu.optimizeAssignments(instructions)
	instructions = alu.optimizeNotEquals(instructions)
	instructions = alu.optimizeAddThenMultiply(instructions)
	instructions = alu.optimizeMultiplyThenAdd(instructions)
	instructions = alu.optimizeModuloThenAdd(instructions)
	instructions = alu.optimizeModuloThenAddIsNotEqual(instructions)
	instructions = alu.optimizeStrangeAlgorithm1(instructions)
	instructions = alu.optimizeStrangeAlgorithm2(instructions)
	return instructions
}

func (alu *ALU) optimizeDivideByOne(instructions []string) []string {
	var optimized []string
	for _, instruction := range instructions {
		tokens := strings.Split(instruction, " ")
		if tokens[0] == "div" && tokens[2] == "1" {
			continue
		}
		optimized = append(optimized, instruction)
	}
	return optimized
}

func (alu *ALU) optimizeAssignments(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 1{
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		if 	tokens[0] == "mul" && tokens[2] == "0" && line2Tokens[0] == "add" && tokens[1] == line2Tokens[1] {
			optimized = append(optimized, fmt.Sprintf("mov %s %s", tokens[1], line2Tokens[2]))
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeNotEquals(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 1 {
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		if tokens[0] == "eql" && line2Tokens[0] == "eql" && tokens[1] == line2Tokens[1] && line2Tokens[2] == "0" {
			optimized = append(optimized, fmt.Sprintf("neq %s %s", tokens[1], tokens[2]))
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeAddThenMultiply(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 1 {
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		if tokens[0] == "add" && line2Tokens[0] == "mul" &&
			tokens[1] == line2Tokens[1] {
			optimized = append(optimized, fmt.Sprintf("adm %s %s %s", tokens[1], tokens[2], line2Tokens[2]))
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeMultiplyThenAdd(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 1 {
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		if tokens[0] == "mul" && line2Tokens[0] == "add" &&
			tokens[1] == line2Tokens[1] {
			optimized = append(optimized, fmt.Sprintf("mad %s %s %s", tokens[1], tokens[2], line2Tokens[2]))
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeModuloThenAdd(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 1 {
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		if tokens[0] == "mod" && line2Tokens[0] == "add" &&
			tokens[1] == line2Tokens[1] {
			optimized = append(optimized, fmt.Sprintf("mda %s %s %s", tokens[1], tokens[2], line2Tokens[2]))
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeModuloThenAddIsNotEqual(instructions []string) []string {
	var optimized []string
	i := 0
	for i < len(instructions) - 2 {
		tokens := strings.Split(instructions[i], " ")
		line2Tokens := strings.Split(instructions[i+1], " ")
		line3Tokens := strings.Split(instructions[i+2], " ")
		if tokens[0] == "mov" && line2Tokens[0] == "mda" && line3Tokens[0] == "neq" &&
			tokens[1] == line2Tokens[1] && tokens[1] == line3Tokens[1] {
			optimized = append(optimized, fmt.Sprintf("mdaneq %s %s %s %s %s", tokens[1], tokens[2], line2Tokens[2], line2Tokens[3], line3Tokens[2]))
			i++
			i++
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeStrangeAlgorithm1(instructions []string) []string {
	detected := func(instructions []string) bool {
		if len(instructions) < 6 { return false	}
		if instructions[0] != "mov y 25" { return false }
		if instructions[1] != "mad y x 1" { return false }
		if instructions[2] != "mul z y" { return false }
		if instructions[3] != "mov y w" { return false } // seems useless
		if instructions[5] != "add z y" { return false }
		tokens := strings.Split(instructions[4], " ")
		if ! (tokens[0] == "adm" && tokens[1] == "y" && tokens[3] == "x") { return false }
		return true
	}
	var optimized []string
	i := 0
	for i < len(instructions) - 5 {
		if detected(instructions[i:]) {
			optimized = append(optimized, fmt.Sprintf("st1 %s", strings.Split(instructions[i+4], " ")[2]))
			i = i + 5
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized
}

func (alu *ALU) optimizeStrangeAlgorithm2(instructions []string) []string {
	detected := func(instructions []string) bool {
		if len(instructions) < 5 { return false	}
		if instructions[0] != "mov x z" { return false }
		if instructions[1] != "mod x 26" { return false }
		if instructions[2] != "div z 26" { return false }
		if instructions[4] != "neq x w" { return false }
		tokens := strings.Split(instructions[3], " ")
		if ! (tokens[0] == "add" && tokens[1] == "x") { return false }
		return true
	}
	var optimized []string
	i := 0
	for i < len(instructions) - 4 {
		if detected(instructions[i:]) {
			optimized = append(optimized, fmt.Sprintf("st2 %s", strings.Split(instructions[i+3], " ")[2]))
			i = i + 4
		} else {
			optimized = append(optimized, instructions[i])
		}
		i++
	}
	for i < len(instructions) {
		optimized = append(optimized, instructions[i])
		i++
	}
	return optimized

}

func (alu *ALU) Evaluate(input string) {
	alu.inputBuffer = input
	if alu.debug {
		fmt.Printf("Evaluating %s\n", input)
	}
	keyBuf := new(bytes.Buffer)
	_, err := keyBuf.WriteString(input)
	if err != nil {
		panic(err)
	}
	err = binary.Write(keyBuf, binary.LittleEndian, alu.Registers['z'])
	if err != nil {
		panic(err)
	}
	key := keyBuf.String()
	if cached, ok := alu.cache[key]; ok {
		alu.cacheHit = true
		if alu.debug {
			fmt.Printf("Cache hit %s\n", key)
		}
		alu.Registers = cached
		return
	}
	alu.cacheHit = false
	for _, instruction := range alu.Program {
		instruction(alu)
	}
	alu.cache[key] = alu.CopyRegisters()
}

func (alu *ALU) ParseInput(tokens []string) Instruction {
	if len(tokens) != 1 {
		panic(fmt.Sprintf("inp expets 1 parameter but found %d", len(tokens)))
	}
	var register rune
	switch tokens[0] {
	case "w":
		register = 'w'
	case "x":
		register = 'x'
	case "y":
		register = 'y'
	case "z":
		register = 'z'
	default:
		panic(fmt.Sprintf("bad register found for inp: %s", tokens[0]))
	}
	return func(alu *ALU) {
		if len(alu.inputBuffer) == 0 {
			fmt.Printf("inp %v failed, no more input\n", tokens)
		}
		v, _ := strconv.ParseInt(string(alu.inputBuffer[0]), 10, 64)
		if alu.debug {
			//fmt.Printf(" - z=%d, %s=%d\n", alu.Registers['z'], string(register), v)
		}
		alu.Registers[register] = v
		alu.inputBuffer = alu.inputBuffer[1:]
	}
}

func (alu *ALU) ParseAdd(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("add expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("add requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		//fmt.Printf("add %s=%d %s=%d\n", tokens[0], alu.read(tokens[0]), tokens[1], alu.read(tokens[1]))
		sum := alu.read(tokens[0]) + alu.read(tokens[1])
		alu.assign(tokens[0], sum)
	}
}

func (alu *ALU) ParseMultiply(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("mul expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mul requires a register in param 1 but found %s", tokens[0]))
	}
	//if tokens[0] == "x" {
		fmt.Printf("Found a mul %s\n", tokens[0])
	//}
	return func(alu *ALU) {
		//fmt.Printf("mul %s=%d %s=%d\n", tokens[0], alu.read(tokens[0]), tokens[1], alu.read(tokens[1]))
		product := alu.read(tokens[0]) * alu.read(tokens[1])
		alu.assign(tokens[0], product)
	}
}

func (alu *ALU) ParseDivide(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("div expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("div requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		//fmt.Printf("di %s=%d %s=%d\n", tokens[0], alu.read(tokens[0]), tokens[1], alu.read(tokens[1]))
		quotient := alu.read(tokens[0]) / alu.read(tokens[1])
		alu.assign(tokens[0], quotient)
	}
}

func (alu *ALU) ParseModulo(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("mod expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mod requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		//fmt.Printf("mod %s=%d %s=%d\n", tokens[0], alu.read(tokens[0]), tokens[1], alu.read(tokens[1]))
		modulus := alu.read(tokens[0]) % alu.read(tokens[1])
		alu.assign(tokens[0], modulus)
	}
}

func (alu *ALU) ParseEquals(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("eql expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("eql requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		equality := 0
		if alu.read(tokens[0]) == alu.read(tokens[1]) {
			equality = 1
		}
		alu.assign(tokens[0], int64(equality))
	}
}

func (alu *ALU) ParseNotEquals(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("neq expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("neq requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		inequality := 1
		if alu.read(tokens[0]) == alu.read(tokens[1]) {
			inequality = 0
		}
		alu.assign(tokens[0], int64(inequality))
	}
}

func (alu *ALU) ParseMove(tokens []string) Instruction {
	if len(tokens) != 2 {
		panic(fmt.Sprintf("mov expects 2 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mov requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		alu.assign(tokens[0], alu.read(tokens[1]))
	}
}

func (alu *ALU) ParseAddMultiply(tokens []string) Instruction {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("adm expects 3 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("adm requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		alu.assign(tokens[0], (alu.read(tokens[0]) + alu.read(tokens[1])) * alu.read(tokens[2]))
	}
}

func (alu *ALU) ParseMultiplyAdd(tokens []string) Instruction {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("mad expects 3 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mad requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		alu.assign(tokens[0], alu.read(tokens[0]) * alu.read(tokens[1]) + alu.read(tokens[2]))
	}
}

func (alu *ALU) ParseModuloThenAdd(tokens []string) Instruction {
	if len(tokens) != 3 {
		panic(fmt.Sprintf("mda expects 3 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mda requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		alu.assign(tokens[0], (alu.read(tokens[0]) % alu.read(tokens[1])) + alu.read(tokens[2]))
	}
}


func (alu *ALU) ParseModuloThenAddIsNotEqual(tokens []string) Instruction {
	if len(tokens) != 5 {
		panic(fmt.Sprintf("mdaneq expects 5 parameter but found %d", len(tokens)))
	}
	if ! alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("mdaneq requires a register in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		moduloThenAdd :=  (alu.read(tokens[1]) % alu.read(tokens[2])) + alu.read(tokens[3])
		inequality := 1
		if moduloThenAdd == alu.read(tokens[4]) {
			inequality = 0
		}
		//fmt.Printf("- mdaneq %d\n", inequality)
		alu.assign(tokens[0], int64(inequality))
	}
}

func (alu *ALU) ParseStrangeAlgo1(tokens []string) Instruction {
	if len(tokens) != 1 {
		panic(fmt.Sprintf("st1 expects 1 parameter but found %d", len(tokens)))
	}
	if alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("st1 requires a number in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		w, x, z := alu.Registers['w'], alu.Registers['x'], alu.Registers['z']
		y := 25 * x + 1 // not really storing in y register, cross your fingers
		if (alu.debug) {
			//fmt.Printf("  - st1 w=%d x=%d z=%d result=%d\n", w, x, z, z * y + x * (w + alu.read(tokens[0])))
		}
		alu.Registers['z'] = z * y + x * (w + alu.read(tokens[0]))
	}
}

func (alu *ALU) ParseStrangeAlgo2(tokens []string) Instruction {
	if len(tokens) != 1 {
		panic(fmt.Sprintf("st2 expects 1 parameter but found %d", len(tokens)))
	}
	if alu.isRegister(tokens[0]) {
		panic(fmt.Sprintf("st2 requires a number in param 1 but found %s", tokens[0]))
	}
	return func(alu *ALU) {
		w, z := alu.Registers['w'], alu.Registers['z']
		alu.Registers['z'] = z / 26
		x := z % 26 + alu.read(tokens[0])
		inequality := 1
		if (w == x) {
			inequality = 0
		}
		alu.Registers['x'] = int64(inequality)
	}
}

func (alu ALU) read(value string) int64 {
	switch value {
	case "w":
		return alu.Registers['w']
	case "x":
		return alu.Registers['x']
	case "y":
		return alu.Registers['y']
	case "z":
		return alu.Registers['z']
	default:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("failed to read value %s: %v", value, err))
		}
		return i
	}
}

func (alu *ALU) assign(target string, value int64) {
	if target == "y" {
		fmt.Printf(" - assign %s=%d\n", target, value)
	}

	var register rune
	switch target {
	case "w":
		register = 'w'
	case "x":
		register = 'x'
	case "y":
		register = 'y'
	case "z":
		register = 'z'
	default:
		panic(fmt.Sprintf("bad register found for inp: %s", target))
	}
	alu.Registers[register] = value
}

func (alu ALU) isRegister(token string) bool {
	switch token {
	case "w":
		return true
	case "x":
		return true
	case "y":
		return true
	case "z":
		return true
	default:
		return false
	}
}