package main

import (
	"bufio"
	"os"
	"strings"
)

type Instruction struct {
	Command string
	Args    []string
	Raw     string
}

func ParseDBashFile(path string) ([]Instruction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		instructions = append(instructions, Instruction{
			Command: strings.ToUpper(parts[0]),
			Args:    parts[1:],
			Raw:     line,
		})
	}

	return instructions, scanner.Err()
}
