package main

import (
	"fmt"
	"micro-agent/internal/config"
	"micro-agent/internal/tui"
	"os"
	"path"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// Ensure we're in the correct directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	err := os.Chdir(path.Join(dir, "../.."))
	if err != nil {
		panic(err)
	}
}

func main() {
	// Check if config exists
	if _, err := config.Load(); err == nil {
		fmt.Println("Configuration already exists at ~/.config/micro-agent.yml")
		return
	}

	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
