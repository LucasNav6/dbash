package loggers

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func Init(isDebugFlagActived bool) *log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})

	styles := log.DefaultStyles()

	styles.Timestamp = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("[✓] OKEY").
		Foreground(lipgloss.Color("245"))

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("[!] WARN").
		Foreground(lipgloss.Color("229"))

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("[X] ERRO").
		Foreground(lipgloss.Color("210"))

	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("[?] DEBU").
		Foreground(lipgloss.Color("225"))

	logger.SetStyles(styles)

	if isDebugFlagActived {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	return logger
}
