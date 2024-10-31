package tui

import (
	"fmt"
	"micro-agent/internal/anthropic"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Global program reference
var program *tea.Program

// InitProgram sets up the global program reference
func InitProgram(p *tea.Program) {
	program = p
}

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			Foreground(lipgloss.Color("240"))
)

// Message types for our update loop
type (
	streamMsg      string
	streamComplete struct{}
	errorMsg       error
)

type ChatModel struct {
	client    *anthropic.Client
	model     string
	viewport  viewport.Model
	textarea  textarea.Model
	messages  []string
	streaming bool
	err       error
	height    int
	width     int
}

func NewChatModel(apiKey, model string) ChatModel {
	ta := textarea.New()
	ta.Focus()
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.Placeholder = "Type your message..."

	vp := viewport.New(80, 20)
	vp.SetContent("")

	return ChatModel{
		client:   anthropic.NewClient(apiKey),
		model:    model,
		textarea: ta,
		viewport: vp,
		messages: make([]string, 0),
		height:   24,
		width:    80,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

// Stream command that sends chunks back to our update loop
func streamResponse(client *anthropic.Client, model, message string) tea.Cmd {
	return func() tea.Msg {
		// Create a channel for the streaming chunks
		streamChan := make(chan tea.Msg)

		// Start streaming in a goroutine
		go func() {
			err := client.SendMessage(model, message, func(chunk string) {
				streamChan <- streamMsg(chunk)
			})
			if err != nil {
				streamChan <- errorMsg(err)
			}

			streamChan <- streamComplete{}
			close(streamChan)
		}()

		// Return the first message
		msg := <-streamChan

		// Create a command for subsequent messages
		go func() {
			for msg := range streamChan {
				switch msg.(type) {
				case streamMsg, streamComplete, errorMsg:
					program.Send(msg)
				}
			}
		}()

		return msg
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC:
			return m, tea.Quit

		case msg.String() == "shift+enter":
			m.textarea, tiCmd = m.textarea.Update(msg)
			return m, tiCmd

		case msg.Type == tea.KeyEnter:
			if !m.streaming && strings.TrimSpace(m.textarea.Value()) != "" {
				userMsg := m.textarea.Value()
				m.messages = append(m.messages, "You: "+userMsg)
				m.messages = append(m.messages, "Assistant: ")
				m.textarea.Reset()
				m.streaming = true
				m.updateViewport()

				return m, streamResponse(m.client, m.model, userMsg)
			}
		}

	case streamMsg:
		// Append the new chunk to the last message
		m.messages[len(m.messages)-1] = m.messages[len(m.messages)-1] + string(msg)
		m.updateViewport()
		return m, nil

	case streamComplete:
		m.streaming = false
		return m, nil

	case errorMsg:
		m.err = msg
		m.streaming = false
		return m, nil

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		headerHeight := 0
		footerHeight := 2
		verticalMarginHeight := 2

		m.textarea.SetWidth(msg.Width)

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight - footerHeight - verticalMarginHeight - m.textarea.Height()
	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	cmds = append(cmds, tiCmd, vpCmd)
	return m, tea.Batch(cmds...)
}

func (m *ChatModel) updateViewport() {
	m.viewport.SetContent(strings.Join(m.messages, "\n\n"))
	m.viewport.GotoBottom()
}

func (m ChatModel) helpView() string {
	return infoStyle.Render(fmt.Sprintf(
		"↑/↓: scroll • shift+enter: new line • enter: send message • ctrl+c: quit",
	))
}

func (m ChatModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n%s",
		m.viewport.View(),
		m.textarea.View(),
		m.helpView(),
	)
}
