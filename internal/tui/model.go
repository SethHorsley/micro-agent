package tui

import (
	"micro-agent/internal/config"
	"micro-agent/internal/providers"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type setupState int

const (
	stateSelectLargeProvider setupState = iota
	stateEnterLargeProviderKey
	stateSelectLargeModel
	stateSelectSmallProvider
	stateEnterSmallProviderKey // Only used if different from large provider
	stateSelectSmallModel
	stateDone
)

type Model struct {
	state     setupState
	config    *config.Config
	textInput textinput.Model
	err       error
	choices   []string
	cursor    int
	providers map[string][]string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return Model{
		state:     stateSelectLargeProvider,
		config:    &config.Config{Providers: make(map[string]config.Provider)},
		textInput: ti,
		choices:   []string{"OpenAI", "Anthropic"},
		providers: providers.GetAvailableModels(),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
