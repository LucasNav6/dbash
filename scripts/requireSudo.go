package instructions

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func RequireSudo(logger *log.Logger) {
	logger.Debug("Verifying if the user executed the command with 'sudo'...")

	var choice string

	prompt := huh.NewSelect[string]().
		Title("This command requires sudo privileges. Do you want to continue with sudo?").
		Options(
			huh.NewOption("[YES]", "yes"),
			huh.NewOption("[NO]", "no"),
		).
		Value(&choice)

	if err := prompt.Run(); err != nil {
		logger.Warn("Execution aborted: failed to confirm sudo privileges.")
		return
	}

	if choice != "yes" {
		logger.Warn("Execution aborted: sudo privileges not granted by the user.")
		return
	}

	logger.Info("Sudo confirmed by user. Proceeding with execution...")
}
