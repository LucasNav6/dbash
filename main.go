package main

import (
	"flag"
	"os"

	"github.com/LucasNav6/dbash/internal/loggers"
	"github.com/LucasNav6/dbash/internal/utils"
)

func main() {
	localPath := flag.String("local", "", "Path to .dbash file")
	isDebugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	logger := loggers.Init(*isDebugFlag)

	logger.Debug("Debug mode is ON. Remove '--debug' to turn off.")

	if *localPath == "" {
		logger.Error("No .dbash file specified. Use --local=<path> to provide one.")
		os.Exit(1)
	}

	logger.Info("Running the provided .dbash file.")

	instructions := utils.ReadHushFile(localPath, logger)

	for i, instruction := range instructions {
		logger.Debug("New instruction readed", "line", i, "command", instruction.Command, "args", instruction.Args)
		utils.ReadInstruction(&instruction, logger)
	}

	logger.Info("Execution .dbash file finished successfully")
}
