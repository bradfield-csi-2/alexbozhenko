package main

import "fmt"

const OUTPUT_BYTE = 0x00
const INSTRUCTIONS_START_BYTE = 0x08

const MEMORY_SIZE = 256

const (
	LOAD  = 0x01
	STORE = 0x02
	ADD   = 0x03
	SUB   = 0x04
	ADDI  = 0x05
	SUBI  = 0x06
	JUMP  = 0x07
	BEQZ  = 0x08
	HALT  = 0xff
)

const R1 = 0x01
const R2 = 0x02
const PC = 0x00 // two general purpose registers and a Program Counter

func compute(memory []byte) {

	regs := map[byte]byte{R1: 0x00, R2: 0x00, PC: INSTRUCTIONS_START_BYTE}

	for {
		if regs[PC] > byte(cap(memory)-3) {
			return
		}
		operation := memory[regs[PC]]
		switch operation {
		case HALT:
			return
		case ADD:
			first_reg_arg := memory[regs[PC]+1]
			second_reg_arg := memory[regs[PC]+2]
			regs[first_reg_arg] += regs[second_reg_arg]
		case ADDI:
			reg_arg := memory[regs[PC]+1]
			value_to_add := memory[regs[PC]+2]
			regs[reg_arg] += value_to_add
		case SUB:
			first_reg_arg := memory[regs[PC]+1]
			second_reg_arg := memory[regs[PC]+2]
			regs[first_reg_arg] -= regs[second_reg_arg]
		case SUBI:
			reg_arg := memory[regs[PC]+1]
			value_to_substract := memory[regs[PC]+2]
			regs[reg_arg] -= value_to_substract
		case LOAD:
			reg_arg := memory[regs[PC]+1]
			address_to_load := memory[regs[PC]+2]
			regs[reg_arg] = memory[address_to_load]
		case STORE:
			reg_arg := memory[regs[PC]+1]
			where_to_store := memory[regs[PC]+2]
			if where_to_store >= INSTRUCTIONS_START_BYTE {
				// invalid instruction
				return
			}
			memory[where_to_store] = regs[reg_arg]
		case JUMP:
			jump_to_absolute_val := memory[regs[PC]+1]
			regs[PC] = jump_to_absolute_val
			continue
		case BEQZ:
			reg_arg := memory[regs[PC]+1]
			relative_offset := memory[regs[PC]+2]
			if regs[reg_arg] == 0 {
				regs[PC] = regs[PC] + relative_offset
			}
		default:
			return
		}
		regs[PC] += 3
	}

}

func main() {
	var memory []byte = []byte{0x00,
		0x03, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x01, 0x01,
		0x01, 0x02, 0x02,
		0x03, 0x01, 0x02,
		0x02, 0x01, 0x00,
		0xff}

	compute(memory)
	fmt.Println(memory)

}
