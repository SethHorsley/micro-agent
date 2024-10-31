package main

import (
	"fmt"
	"micro-agent/internal/config"
	"micro-agent/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err == nil {
		// Config exists, start chat interface
		var apiKey string
		if provider, ok := cfg.Providers["anthropic"]; ok {
			apiKey = provider.APIKey
		}

		model := cfg.LargeProvider.Model

		// Create and initialize the program
		p := tea.NewProgram(tui.NewChatModel(apiKey, model))
		tui.InitProgram(p) // Add this line

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running chat: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Config doesn't exist, run setup
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
