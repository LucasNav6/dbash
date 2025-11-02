package main

import (
	"flag"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var logger *log.Logger

func main() {
	localPath := flag.String("local", "", "Path to .dbash file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Estilos personalizados
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("🐛 DEBU").
		Padding(0, 1).
		Foreground(lipgloss.Color("240"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("✨ INFO").
		Padding(0, 1).
		Foreground(lipgloss.Color("86"))
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("⚠️  WARN").
		Padding(0, 1).
		Foreground(lipgloss.Color("214"))
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("❌ ERRO").
		Padding(0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))

	// Keys estilo
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	// Logger principal
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: *debug,
		Prefix:       "dbash ",
	})
	logger.SetStyles(styles)

	if !*debug {
		logger.SetLevel(log.InfoLevel)
	}

	if *localPath == "" {
		logger.Error("You must provide a .dbash file with --local=<path>")
		os.Exit(1)
	}

	logger.Debug("Parsing .dbash file", "path", *localPath)
	instructions, err := ParseDBashFile(*localPath)
	if err != nil {
		logger.Error("Error reading file", "err", err)
		os.Exit(1)
	}

	logger.Debug("Executing instructions", "count", len(instructions))
	err = ExecuteInstructions(instructions)
	if err != nil {
		logger.Error("Error executing instructions", "err", err)
		os.Exit(1)
	}

	logger.Info("Execution finished successfully")
}
