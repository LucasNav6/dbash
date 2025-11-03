package utils

import (
	"os"
	"strings"

	"github.com/LucasNav6/dbash/models"
	instructions "github.com/LucasNav6/dbash/scripts"
	"github.com/charmbracelet/log"
)

func ReadInstruction(instruction *models.Instruction, logger *log.Logger) {
	switch instruction.Command {
	case "REQUIRE_SUDO":
		instructions.RequireSudo(logger)

	case "REQUIRE":
		instructions.Require(*instruction, logger)

	case "ARGS":
		instructions.Args(instruction.Args, logger)

	case "SET":
		// Unimos todos los argumentos menos el último con un espacio
		prompt := strings.Join(instruction.Args[:len(instruction.Args)-1], " ")
		// El nombre de la variable está en el último argumento
		varName := instruction.Args[len(instruction.Args)-1]
		instructions.Set(prompt, varName, logger)

	default:
		logger.Error("Invalid command. Please check the syntax.", "command not found", instruction.Command)
		os.Exit(1)
	}
}
