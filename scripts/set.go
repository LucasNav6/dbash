package instructions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func Set(promptTemplate string, varName string, logger *log.Logger) {

	varName = strings.TrimPrefix(varName, "$")
	if varName == "" {
		logger.Error("SET: empty variable name")
		os.Exit(1)
	}

	vars, err := loadVars()
	if err != nil {
		logger.Error("SET: failed reading temp_vars.txt", "error", err)
		os.Exit(1)
	}

	promptText, err := expandPlaceholders(promptTemplate, vars)
	if err != nil {
		logger.Error("SET: unknown variable in prompt", "error", err)
		os.Exit(1)
	}

	var newValue string
	input := huh.NewInput().
		Title(promptText).
		Value(&newValue)

	if err := input.Run(); err != nil {
		logger.Error("SET: user cancelled input")
		os.Exit(1)
	}

	newValue = strings.TrimSpace(newValue)
	if newValue == "" {
		logger.Debug("SET: empty input, variable left unchanged", "variable", varName)
		return
	}

	if _, exists := vars[varName]; !exists {
		logger.Error("SET: variable does not exist", "variable", varName)
		os.Exit(1)
	}
	vars[varName] = newValue

	if err := saveVars(vars); err != nil {
		logger.Error("SET: failed writing temp_vars.txt", "error", err)
		os.Exit(1)
	}

	logger.Debug("Variable updated", "variable", varName, "value", newValue)
}

func loadVars() (map[string]string, error) {
	m := make(map[string]string)
	f, err := os.Open(filepath.Join(".", "temp_vars.txt"))
	if os.IsNotExist(err) {
		return m, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m, sc.Err()
}

func saveVars(vars map[string]string) error {
	f, err := os.Create(filepath.Join(".", "temp_vars.txt"))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for k, v := range vars {
		_, _ = fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func expandPlaceholders(tmpl string, vars map[string]string) (string, error) {
	re := regexp.MustCompile(`\{\{\s*\$([A-Za-z_][A-Za-z0-9_]*)\s*\}\}`)
	return re.ReplaceAllStringFunc(tmpl, func(raw string) string {
		name := re.FindStringSubmatch(raw)[1]
		val, ok := vars[name]
		if !ok {

			panic(fmt.Errorf("variable %s not found", name))
		}
		return val
	}), nil
}
