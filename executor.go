package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"

// 	"github.com/charmbracelet/bubbles/spinner"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/huh"
// 	"github.com/charmbracelet/lipgloss"
// )

// // ---------- VARS FILE ----------
// var varsFilePath string

// func initVarsFile(scriptName string) error {
// 	dir := "vars"
// 	if err := os.MkdirAll(dir, 0755); err != nil {
// 		return err
// 	}
// 	varsFilePath = filepath.Join(dir, scriptName+".vars")
// 	return nil
// }

// func appendVarToFile(key, value string) {
// 	if varsFilePath == "" {
// 		logger.Debug("appendVarToFile: varsFilePath no inicializado")
// 		return
// 	}
// 	value = strings.ReplaceAll(value, `"`, `\"`)

// 	var lines []string
// 	file, err := os.Open(varsFilePath)
// 	if err == nil {
// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			line := scanner.Text()
// 			if strings.HasPrefix(line, fmt.Sprintf(`$%s=`, key)) {
// 				continue
// 			}
// 			lines = append(lines, line)
// 		}
// 		file.Close()
// 	}
// 	lines = append(lines, fmt.Sprintf(`$%s="%s"`, key, value))
// 	content := strings.Join(lines, "\n") + "\n"
// 	if err := os.WriteFile(varsFilePath, []byte(content), 0644); err != nil {
// 		logger.Error("appendVarToFile: failed to write vars file", "file", varsFilePath, "error", err)
// 	}
// }

// // ---------- CLONE MODEL ----------
// type cloneModel struct {
// 	spinner spinner.Model
// 	url     string
// 	err     error
// 	done    bool
// }

// func newCloneModel(url string) *cloneModel {
// 	s := spinner.New()
// 	s.Spinner = spinner.Dot
// 	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
// 	return &cloneModel{spinner: s, url: url}
// }

// func (m cloneModel) Init() tea.Cmd {
// 	return tea.Batch(
// 		m.spinner.Tick,
// 		func() tea.Msg {
// 			cmd := exec.Command("git", "clone", "--depth=1", m.url)
// 			if err := cmd.Run(); err != nil {
// 				return cloneErrMsg{err}
// 			}
// 			return cloneDoneMsg{}
// 		},
// 	)
// }

// type cloneDoneMsg struct{}
// type cloneErrMsg struct{ error }

// func (m cloneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case cloneDoneMsg:
// 		m.done = true
// 		return m, tea.Quit
// 	case cloneErrMsg:
// 		m.err = msg.error
// 		return m, tea.Quit
// 	default:
// 		var cmd tea.Cmd
// 		m.spinner, cmd = m.spinner.Update(msg)
// 		return m, cmd
// 	}
// }

// func (m cloneModel) View() string {
// 	if m.err != nil {
// 		return fmt.Sprintf("❌ Failed to clone %s\n", m.url)
// 	}
// 	if m.done {
// 		return fmt.Sprintf("✅ Cloned %s\n", m.url)
// 	}
// 	return m.spinner.View() + " Cloning " + m.url + "..."
// }

// // ---------- EXECUTORS ----------
// var executors = map[string]func(Instruction) error{
// 	"REQUIRE_SUDO": func(i Instruction) error {
// 		logger.Debug("REQUIRE_SUDO: checking if sudo is needed")
// 		var confirm bool
// 		err := huh.NewConfirm().
// 			Title("Do you want to continue with sudo privileges?").
// 			Value(&confirm).
// 			Run()
// 		if err != nil {
// 			logger.Debug("REQUIRE_SUDO: user failed to confirm", "error", err)
// 			return fmt.Errorf("sudo confirmation failed: %w", err)
// 		}
// 		if !confirm {
// 			logger.Debug("REQUIRE_SUDO: user declined sudo")
// 			os.Exit(0)
// 		}
// 		logger.Debug("REQUIRE_SUDO: user confirmed sudo")
// 		return nil
// 	},

// 	"REQUIRE": func(i Instruction) error {
// 		cmd := i.Args[0]
// 		version := ""
// 		if strings.Contains(cmd, "@") {
// 			parts := strings.Split(cmd, "@")
// 			cmd = parts[0]
// 			version = parts[1]
// 		}
// 		logger.Debug("REQUIRE: checking tool", "tool", cmd, "version", version)
// 		if _, err := exec.LookPath(cmd); err != nil {
// 			logger.Debug("REQUIRE: tool not found", "tool", cmd)
// 			return fmt.Errorf("'%s' not found", cmd)
// 		}
// 		logger.Debug("REQUIRE: tool found", "tool", cmd)
// 		return nil
// 	},

// 	"ARGS": func(i Instruction) error {
// 		varName := strings.TrimPrefix(i.Args[0], "$")
// 		value := strings.Trim(i.Args[1], `"`)
// 		Vars[varName] = value
// 		appendVarToFile(varName, value)
// 		logger.Debug("ARGS: set variable", "var", varName, "value", value)
// 		return nil
// 	},

// 	"ASK": func(i Instruction) error {
// 		if len(i.Args) < 2 {
// 			return fmt.Errorf("ASK: need prompt and $var")
// 		}
// 		rawVarName := i.Args[len(i.Args)-1]
// 		prompt := strings.Join(i.Args[:len(i.Args)-1], " ")
// 		if !strings.HasPrefix(rawVarName, "$") {
// 			return fmt.Errorf("variable must start with $")
// 		}
// 		varName := strings.TrimPrefix(rawVarName, "$")
// 		if _, ok := Vars[varName]; !ok {
// 			return fmt.Errorf("variable '%s' not declared", varName)
// 		}
// 		for k, v := range Vars {
// 			prompt = strings.ReplaceAll(prompt, fmt.Sprintf("{{$%s}}", k), v)
// 		}
// 		if prompt == "" {
// 			prompt = fmt.Sprintf("Value for '%s': ", varName)
// 		}
// 		logger.Debug("ASK: prompting user", "var", varName, "prompt", prompt)
// 		var input string
// 		err := huh.NewInput().Title(prompt).Value(&input).Run()
// 		if err != nil {
// 			logger.Debug("ASK: input failed", "var", varName, "error", err)
// 			return fmt.Errorf("input failed: %w", err)
// 		}
// 		if input != "" {
// 			Vars[varName] = input
// 			appendVarToFile(varName, input)
// 			logger.Debug("ASK: variable updated", "var", varName, "newValue", input)
// 		} else {
// 			logger.Debug("ASK: no input, keeping current value", "var", varName, "value", Vars[varName])
// 		}
// 		return nil
// 	},

// 	"CLONE": func(i Instruction) error {
// 		url := i.Args[0]
// 		p := tea.NewProgram(newCloneModel(url), tea.WithOutput(os.Stderr))
// 		finalModel, err := p.Run()
// 		if err != nil {
// 			return err
// 		}
// 		m := finalModel.(*cloneModel)
// 		if m.err != nil {
// 			return fmt.Errorf("git clone failed: %w", m.err)
// 		}
// 		return nil
// 	},

// 	"WORKDIR": func(i Instruction) error {
// 		dir := i.Args[0]
// 		logger.Debug("WORKDIR: changing directory", "dir", dir)
// 		if err := os.Chdir(dir); err != nil {
// 			logger.Debug("WORKDIR: failed", "dir", dir, "error", err)
// 			return fmt.Errorf("chdir failed: %w", err)
// 		}
// 		logger.Debug("WORKDIR: success", "dir", dir)
// 		return nil
// 	},

// 	"CHECKOUT": func(i Instruction) error {
// 		tag := i.Args[0]
// 		for k, v := range Vars {
// 			tag = strings.ReplaceAll(tag, "${"+k+"}", v)
// 		}
// 		logger.Debug("CHECKOUT: checking out", "tag", tag)
// 		cmd := exec.Command("git", "checkout", tag)
// 		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			logger.Debug("CHECKOUT: failed", "tag", tag, "error", err)
// 			return fmt.Errorf("git checkout failed: %w", err)
// 		}
// 		logger.Debug("CHECKOUT: success", "tag", tag)
// 		return nil
// 	},

// 	"RUN": func(i Instruction) error {
// 		cmdLine := strings.Join(i.Args, " ")
// 		for k, v := range Vars {
// 			cmdLine = strings.ReplaceAll(cmdLine, "${"+k+"}", v)
// 		}
// 		logger.Debug("RUN: executing command", "cmd", cmdLine)
// 		cmd := exec.Command("sh", "-c", cmdLine)
// 		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			logger.Debug("RUN: failed", "cmd", cmdLine, "error", err)
// 			return fmt.Errorf("command failed: %w", err)
// 		}
// 		logger.Debug("RUN: success", "cmd", cmdLine)
// 		return nil
// 	},
// }

// var Vars = map[string]string{}

// func ExecuteInstructions(instructions []Instruction) error {
// 	for _, inst := range instructions {
// 		execFunc, ok := executors[inst.Command]
// 		if !ok {
// 			fmt.Printf("⚠️ Unknown command: %s\n", inst.Raw)
// 			continue
// 		}
// 		if err := execFunc(inst); err != nil {
// 			return fmt.Errorf("Error executing %s: %v", inst.Raw, err)
// 		}
// 	}
// 	return nil
// }
