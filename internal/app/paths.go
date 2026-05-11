package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolveOutputDir(args []string) (string, error) {
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		return filepath.Abs(args[0])
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exePath), "books"), nil
}

func (s *Service) OpenDir() error {
	absDir, err := filepath.Abs(s.OutputDir)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", absDir)
	case "darwin":
		cmd = exec.Command("open", absDir)
	default:
		cmd = exec.Command("xdg-open", absDir)
	}
	return cmd.Start()
}
