package singleselectsearch

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice *item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = &i
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != nil {
		return lipgloss.NewStyle().Margin(0, 0, 1, 0).Render(fmt.Sprintf("Choose a programming language: %s", m.choice.Title()))
	}
	return docStyle.Render(m.list.View())
}

func Model() tea.Model {
	items := []list.Item{
		item{title: "HTML and CSS"},
		item{title: "Python"},
		item{title: "Java"},
		item{title: "JavaScript"},
		item{title: "Swift"},
		item{title: "C++"},
		item{title: "C#"},
		item{title: "R"},
		item{title: "Golang (Go)"},
		item{title: "PHP"},
		item{title: "TypeScript"},
		item{title: "Scala"},
		item{title: "Shell"},
		item{title: "PowerShell"},
		item{title: "Perl"},
		item{title: "Haskell"},
		item{title: "Kotlin"},
		item{title: "Visual Basic .NET"},
		item{title: "SQL"},
		item{title: "Delphi"},
		item{title: "MATLAB"},
		item{title: "Groovy"},
		item{title: "Lua"},
		item{title: "Rust"},
		item{title: "Ruby"},
		item{title: "C"},
		item{title: "Dart"},
		item{title: "DM"},
	}

	l := list.NewDefaultDelegate()
	l.ShowDescription = false
	l.SetSpacing(0)
	m := model{list: list.New(items, l, 0, 0)}
	m.list.Title = "One more time but you can use a filter:"

	return &m
}
