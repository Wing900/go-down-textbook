package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenyb-go/go-down-textbook/internal/app"
	"github.com/chenyb-go/go-down-textbook/internal/tui"
)

func main() {
	outputDir, err := app.ResolveOutputDir(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析输出目录失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	program := tea.NewProgram(
		tui.NewModel(app.NewService(outputDir)),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI 运行失败: %v\n", err)
		os.Exit(1)
	}
}
