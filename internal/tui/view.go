package tui

import "fmt"

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	switch m.state {
	case stateSelectLargeProvider:
		return m.renderChoices("Select your large language model provider:")

	case stateEnterLargeProviderKey:
		return fmt.Sprintf("Enter your %s API Key:\n\n%s",
			m.config.LargeProvider.Provider, m.textInput.View())

	case stateSelectLargeModel:
		return m.renderChoices("Select the model for your large provider:")

	case stateSelectSmallProvider:
		return m.renderChoices("Select your small language model provider:")

	case stateEnterSmallProviderKey:
		return fmt.Sprintf("Enter your %s API Key:\n\n%s",
			m.config.SmallProvider.Provider, m.textInput.View())

	case stateSelectSmallModel:
		return m.renderChoices("Select the model for your small provider:")

	case stateDone:
		return "Configuration saved successfully!\n"

	default:
		return "Unknown state\n"
	}
}

func (m Model) renderChoices(title string) string {
	s := title + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s + "\nUse arrow keys to select and Enter to confirm\n"
}
