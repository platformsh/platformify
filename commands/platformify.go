package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cfgFile string

// PlatformifyCmd represents the base Platformify command when called without any subcommands
var PlatformifyCmd = &cobra.Command{
	Use:   "platformify",
	Short: "Platfomrify your application, and deploy it to the Platform.sh",
	Long: `Platformify your application, creating all the needed files
for it to be deployed to Platform.sh.

This will create the needed YAML files for both your application and your
servcices, choosing from a variety of stacks or simple runtimes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hey from Platformify!")
	},
}
