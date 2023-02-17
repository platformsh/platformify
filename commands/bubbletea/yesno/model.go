package yesno

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/confirmation"
)

func Model() tea.Model {
	input := confirmation.New("Do you like PHP?", confirmation.No)

	return confirmation.NewModel(input)
}
