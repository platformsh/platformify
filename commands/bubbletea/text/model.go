package text

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	choice    string
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.choice = m.textInput.Value()
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.choice != "" {
		return fmt.Sprintf(
			"What is your name? %s\n",
			m.choice,
		) + "\n"
	}
	return fmt.Sprintf(
		"What is your name?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func Model() tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Username"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &model{
		textInput: ti,
		err:       nil,
	}
}
