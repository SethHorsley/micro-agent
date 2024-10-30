package tui

import (
	"micro-agent/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			switch m.state {
			case stateSelectLargeProvider:
				m.config.LargeProvider.Provider = m.getSelectedProvider()
				m.choices = m.providers[m.getSelectedProvider()]
				m.cursor = 0
				m.state = stateSelectLargeModel

			case stateSelectLargeModel:
				m.config.LargeProvider.Model = m.choices[m.cursor]
				m.choices = []string{"OpenAI", "Anthropic"}
				m.cursor = 0
				m.state = stateSelectSmallProvider

			case stateSelectSmallProvider:
				m.config.SmallProvider.Provider = m.getSelectedProvider()
				m.choices = m.providers[m.getSelectedProvider()]
				m.cursor = 0
				m.state = stateSelectSmallModel

			case stateSelectSmallModel:
				m.config.SmallProvider.Model = m.choices[m.cursor]
				m.state = stateEnterOpenAIKey

			case stateEnterOpenAIKey:
				m.config.Providers["open_ai"] = config.Provider{APIKey: m.textInput.Value()}
				m.textInput.Reset()
				m.state = stateEnterAnthropicKey

			case stateEnterAnthropicKey:
				m.config.Providers["anthropic"] = config.Provider{APIKey: m.textInput.Value()}
				if err := config.Save(m.config); err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.state = stateDone
				return m, tea.Quit
			}
		}

	case tea.Msg:
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) getSelectedProvider() string {
	provider := m.choices[m.cursor]
	if provider == "OpenAI" {
		return "open_ai"
	}
	return "anthropic"
}
