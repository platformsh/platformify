package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/commands/bubbletea/singleselect"
	"github.com/platformsh/platformify/commands/bubbletea/singleselectsearch"
	"github.com/platformsh/platformify/commands/bubbletea/text"
	"github.com/platformsh/platformify/commands/bubbletea/yesno"
)

var bubbleteaCmd = &cobra.Command{
	Use: "tea",
	Run: func(cmd *cobra.Command, args []string) {
		m := singleselect.Model()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		m = singleselectsearch.Model()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		m = yesno.Model()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		m = text.Model()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}
