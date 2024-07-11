package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	// Build symbol table
	symbol_table := map[string]string{
		"KBD":    "110000000000000",
		"SCREEN": "100000000000000",
		"R0":     "000000000000000",
		"R1":     "000000000000001",
		"R2":     "000000000000010",
		"R3":     "000000000000011",
		"R4":     "000000000000100",
		"R5":     "000000000000101",
		"R6":     "000000000000110",
		"R7":     "000000000000111",
		"R8":     "000000000001000",
		"R9":     "000000000001001",
		"R10":    "000000000001010",
		"R11":    "000000000001011",
		"R12":    "000000000001100",
		"R13":    "000000000001101",
		"R14":    "000000000001110",
		"R15":    "000000000001111",
		"SP":     "000000000000000",
		"LCL":    "000000000000001",
		"ARG":    "000000000000010",
		"THIS":   "000000000000011",
		"THAT":   "000000000000100",
	}

	dest_table := map[string]string{
		"":    "000",
		"M":   "001",
		"D":   "010",
		"MD":  "011",
		"A":   "100",
		"AM":  "101",
		"AD":  "110",
		"AMD": "111",
	}
	jmp_table := map[string]string{
		"":    "000",
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}
	comp_table := map[string]string{
		"0":   "0101010",
		"1":   "0111111",
		"-1":  "0111010",
		"D":   "0001100",
		"A":   "0110000",
		"!D":  "0001101",
		"!A":  "0110001",
		"-D":  "0001111",
		"-A":  "0110011",
		"D+1": "0011111",
		"A+1": "0110111",
		"D-1": "0001110",
		"A-1": "0110010",
		"D+A": "0000010",
		"D-A": "0010011",
		"A-D": "0000111",
		"D&A": "0000000",
		"D|A": "0010101",
		"M":   "1110000",
		"!M":  "1110001",
		"-M":  "1110011",
		"M+1": "1110111",
		"M-1": "1110010",
		"D+M": "1000010",
		"D-M": "1010011",
		"M-D": "1000111",
		"D&M": "1000000",
		"D|M": "1010101",
	}

	// Open assembly file
	file := open_file("rect/Rect.asm")
	defer file.Close()

	first_pass(file, &symbol_table)

	file.Seek(0, 0)

	output, err := os.Create("Rect.hack")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	second_pass(file, &symbol_table, &comp_table, &dest_table, &jmp_table, output)
	defer output.Close()
}

func open_file(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	return file
}

func first_pass(file *os.File, symbol_table *map[string]string) {
	r := bufio.NewReader(file)
	instruction := 0
	for {
		l, _, err := r.ReadLine()
		line := strings.TrimSpace(string(l))
		if err != nil {
			break
		}
		if len(line) > 0 {
			if line[0] == '@' || line[0] == '0' || line[0] == 'D' || line[0] == 'M' || line[0] == 'A' || line[0] == '-' || line[0] == '1' || line[0] == '!' {
				instruction++
			}
			if line[0] == '(' {
				var t string
				i := 1
				for line[i] != ')' {
					t += string(line[i])
					i++
				}
				if _, ok := (*symbol_table)[t]; !ok {
					(*symbol_table)[t] = fmt.Sprintf("%015b", instruction)
				}
			}
		}
	}
}

func second_pass(file *os.File, symbol_table *map[string]string, comp_table *map[string]string, dest_table *map[string]string, jmp_table *map[string]string, output_file *os.File) {
	r := bufio.NewReader(file)
	var n uint16 = 16
	for {
		l, _, err := r.ReadLine()
		line := strings.TrimSpace(string(l))
		if err != nil {
			break
		}
		if len(line) > 0 {

			if line[0] == '@' {
				var t string
				t += string('0')
				b := line[1:]
				if val, ok := (*symbol_table)[b]; ok {
					t += val
				} else {
					if containsNumber(b) {
						numeric, err := strconv.Atoi(b)
						if err != nil {
							log.Fatalf("Error converting number: %v", err)
						}
						t += fmt.Sprintf("%015b", numeric)
					} else {
						(*symbol_table)[b] = fmt.Sprintf("%015b", n)
						t += fmt.Sprintf("%015b", n)
						n++
					}
				}
				t += string('\n')
				output_file.WriteString(t)
			} else if line[0] == '0' || line[0] == 'D' || line[0] == 'M' || line[0] == 'A' || line[0] == '-' || line[0] == '1' || line[0] == '!' {
				// C-Instruction
				var t string
				t += "111"

				var dest, comp, jmp string

				var tmp string
				for i, char := range line {
					if char == '=' {
						dest = tmp
						tmp = ""
					} else if char == ';' {
						jmp = line[i+1:]
						break
					} else if char == ' ' {
						continue
					} else {
						tmp += string(char)
					}
				}
				comp = tmp

				t += (*comp_table)[comp]
				t += (*dest_table)[dest]
				t += (*jmp_table)[jmp]
				t += string('\n')
				output_file.WriteString(t)
			}
		}
	}

}

func containsNumber(s string) bool {
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
