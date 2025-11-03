package utils

import (
	"bufio"
	"os"
	"strings"

	"github.com/LucasNav6/dbash/models"
	"github.com/charmbracelet/log"
)

func ReadHushFile(localPath *string, logger *log.Logger) []models.Instruction {
	logger.Debug("Reading the provided .dbash file.", "file located", &localPath)

	instructions, err := ParseDBashFile(*localPath)

	if err != nil {
		logger.Error("An error ocurred reading the provided .dbash file")
		os.Exit(1)
	}

	logger.Debug("Read the provided .dbash file success", "total lines", len(instructions))
	return instructions
}

func ParseDBashFile(path string) ([]models.Instruction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var instructions []models.Instruction
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

		instructions = append(instructions, models.Instruction{
			Command: strings.ToUpper(parts[0]),
			Args:    parts[1:],
			Raw:     line,
		})
	}

	return instructions, scanner.Err()
}
