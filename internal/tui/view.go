package tui

import "fmt"

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	switch m.state {
	case stateSelectLargeProvider:
		return m.renderChoices("Select your large language model provider:")

	case stateSelectLargeModel:
		return m.renderChoices("Select the model for your large provider:")

	case stateSelectSmallProvider:
		return m.renderChoices("Select your small language model provider:")

	case stateSelectSmallModel:
		return m.renderChoices("Select the model for your small provider:")

	case stateEnterOpenAIKey:
		return "Enter your OpenAI API Key:\n\n" + m.textInput.View()

	case stateEnterAnthropicKey:
		return "Enter your Anthropic API Key:\n\n" + m.textInput.View()

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
