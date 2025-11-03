package instructions

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/LucasNav6/dbash/models"
	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/log"
)

func Require(instruction models.Instruction, logger *log.Logger) {
	if len(instruction.Args) == 0 {
		logger.Error("Invalid argument. Please check the syntax of REQUIRED instruction.")
		os.Exit(1)
	}

	cmd := instruction.Args[0]
	args := instruction.Args[1:]
	wantVersion := ""
	if strings.Contains(cmd, "@") {
		parts := strings.SplitN(cmd, "@", 2)
		cmd = parts[0]
		wantVersion = parts[1]
	}

	logger.Debug("Verifying tool installation", "tool", cmd, "want", wantVersion)

	if _, err := exec.LookPath(cmd); err != nil {
		logger.Debug("A required tool not found on their device", "tool", cmd)
		logger.Error("Required tool not found on your device. Please install it.", "tool", cmd)
		os.Exit(1)
	}

	var out []byte
	var err error
	runArgs := args
	if len(runArgs) == 0 {
		out, err = exec.Command(cmd, "version").Output()
	}
	if len(runArgs) == 0 && err != nil {
		out, err = exec.Command(cmd, "--version").Output()
	}
	if len(runArgs) > 0 {
		out, err = exec.Command(cmd, runArgs...).Output()
	}

	if err != nil {
		logger.Warn("REQUIRE: tool failed to execute", "tool", cmd, "error", err)
		logger.Error("Required tool not found on your device. Please install it.", "tool", cmd)
		os.Exit(1)
	}

	firstLine := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	m := re.FindStringSubmatch(firstLine)
	if len(m) < 2 {
		logger.Warn("Unable to determine the installed version of the required tool, skipping version check.",
			"tool", cmd, "output", firstLine)
		return
	}

	haveV := semver.MustParse(m[1])

	if wantVersion != "" {
		wantV, err := semver.NewVersion(wantVersion)
		if err != nil {
			logger.Error("REQUIRE: invalid version requested", "tool", cmd, "version", wantVersion)
			os.Exit(1)
		}
		if haveV.LessThan(wantV) {
			logger.Error("REQUIRE: installed version is too old",
				"tool", cmd, "have", haveV.String(), "want", wantV.String())
			os.Exit(1)
		}
	}

	logger.Debug("REQUIRE: tool is installed and meets version requirement",
		"tool", cmd, "version", haveV.String())
}
