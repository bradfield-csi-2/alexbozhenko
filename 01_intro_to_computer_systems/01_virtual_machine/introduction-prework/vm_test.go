package vm

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

type vmCase struct{ x, y, out byte }
type vmTest struct {
	name  string
	asm   string
	cases []vmCase
}

var mainTests = []vmTest{
	// Do nothing, just halt
	{
		name: "Halt",
		asm: `
halt`,
		cases: []vmCase{{0, 0, 0}},
	},
	// Move a value from input to output
	{
		name: "LoadStore",
		asm: `
load r1 1
store r1 0
halt`,
		cases: []vmCase{
			{1, 0, 1},
			{255, 0, 255},
		},
	},
	// Add two unsigned integers together
	{
		name: "Add",
		asm: `
load r1 1
load r2 2
add r1 r2
store r1 0
halt`,
		cases: []vmCase{
			{1, 2, 3},     // 1 + 2 = 3
			{254, 1, 255}, // support max int
			{255, 1, 0},   // correctly overflow
		},
	},
	{
		name: "Subtract",
		asm: `
load r1 1
load r2 2
sub r1 r2
store r1 0
halt`,
		cases: []vmCase{
			{5, 3, 2},
			{0, 1, 255}, // correctly overflow backwards
		},
	},
}

var stretchGoalTests = []vmTest{
	// Support a basic jump, ie skipping ahead to a particular location
	{
		name: "Jump",
		asm: `
load r1 1
jump 16
store r1 0
halt`,
		cases: []vmCase{{42, 0, 0}},
	},
	// Support a "branch if equal to zero" with relative offsets
	{
		name: "Beqz",
		asm: `
load r1 1
load r2 2
beqz r2 3
store r1 0
halt`,
		cases: []vmCase{
			{42, 0, 0},  // r2 is zero, so should branch over the store
			{42, 1, 42}, // r2 is nonzero, so should store back 42
		},
	},
	// Support adding immediate values
	{
		name: "Addi",
		asm: `
load r1 1
addi r1 3
addi r1 5
store r1 0
halt`,
		cases: []vmCase{
			{0, 0, 8},   // 0 + 3 + 5 = 8
			{20, 0, 28}, // 20 + 3 + 5 = 8
		},
	},
	// Calculate the sum of first n numbers (using subi to decrement loop index)
	{
		name: "Sum to n",
		asm: `
load r1 1
beqz r1 8
add r2 r1
subi r1 1
jump 11
store r2 0
halt`,
		cases: []vmCase{
			{0, 0, 0},
			{1, 0, 1},
			{5, 0, 15},
			{10, 0, 55},
		},
	},
}

func TestCompute(t *testing.T) {
	for _, test := range mainTests {
		t.Run(test.name, func(t *testing.T) { testCompute(t, test) })
	}
	if os.Getenv("STRETCH") != "true" {
		println("Skipping stretch goal tests. Run `STRETCH=true go test` to include them.")
	} else {
		for _, test := range stretchGoalTests {
			t.Run(test.name, func(t *testing.T) { testCompute(t, test) })
		}
	}
}

// Given some assembly code and test cases, construct a program
// according to the required memory structure, and run in each
// case through the virtual machine
func testCompute(t *testing.T, test vmTest) {
	// assemble code and load into memory
	memory := make([]byte, 256)
	copy(memory[8:], assemble(test.asm))
	// for each case, set inputs and run vm
	for _, c := range test.cases {
		memory[1] = c.x
		memory[2] = c.y

		compute(memory)

		actual := memory[0]
		if actual != c.out {
			t.Fatalf("Expected f(%d, %d) to be %d, not %d", c.x, c.y, c.out, actual)
		}

		memory[1] = 0
		memory[2] = 0
	}
}

func reg(s string) (b byte) {
	return map[string]byte{
		"r1": 0x01,
		"r2": 0x02,
	}[s]
}

func mem(s string) (b byte) {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return byte(i)
}

func imm(s string) (b byte) {
	// for now, immediate values and memory addresses are both just ints
	return mem(s)
}

// Assemble the given assembly code to machine code
func assemble(asm string) []byte {
	mc := []byte{}
	asm = strings.TrimSpace(asm)
	for _, line := range strings.Split(asm, "\n") {
		parts := strings.Split(strings.TrimSpace(line), " ")
		switch parts[0] {
		case "load":
			mc = append(mc, []byte{0x01, reg(parts[1]), mem(parts[2])}...)
		case "store":
			mc = append(mc, []byte{0x02, reg(parts[1]), mem(parts[2])}...)
		case "add":
			mc = append(mc, []byte{0x03, reg(parts[1]), reg(parts[2])}...)
		case "sub":
			mc = append(mc, []byte{0x04, reg(parts[1]), reg(parts[2])}...)
		case "addi":
			mc = append(mc, []byte{0x05, reg(parts[1]), imm(parts[2])}...)
		case "subi":
			mc = append(mc, []byte{0x06, reg(parts[1]), imm(parts[2])}...)
		case "jump":
			mc = append(mc, []byte{0x07, imm(parts[1])}...)
		case "beqz":
			mc = append(mc, []byte{0x08, reg(parts[1]), imm(parts[2])}...)
		case "halt":
			mc = append(mc, 0xff)
		default:
			panic("Invalid operation: " + parts[0])
		}
	}
	return mc
}
