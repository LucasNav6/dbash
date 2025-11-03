package instructions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func Args(args []string, logger *log.Logger) {

	if len(args) != 3 || args[1] != "=" {
		logger.Error("ARGS syntax error. Use: ARGS $nombre_variable = valor_variable")
		return
	}

	name := strings.TrimPrefix(args[0], "$")
	value := args[2]

	line := fmt.Sprintf("%s=%s\n", name, value)

	f, err := os.OpenFile(filepath.Join(".", "temp_vars.txt"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("Failed to open temp_vars.txt", "error", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		logger.Error("Failed to write variable to temp_vars.txt", "error", err)
		return
	}

	logger.Debug("Variable saved", "file", "temp_vars.txt", "variable", name, "value", value)
}
