package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/pkg/collector"
)

// PlatformifyCmd represents the base Platformify command when called without any subcommands
var PlatformifyCmd = &cobra.Command{
	Use:   "platformify",
	Short: "Platfomrify your application, and deploy it to the Platform.sh",
	Long: `Platformify your application, creating all the needed files
for it to be deployed to Platform.sh.

This will create the needed YAML files for both your application and your
services, choosing from a variety of stacks or simple runtimes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hey from Platformify!")

		// EXAMPLE: this example is made to show how to use Collector to collect user answers.
		c := collector.New()
		c.Add(collector.Runtime("python"))
		c.Add(collector.Stack("django"))
		c.Add(
			collector.Environment("POETRY_VERSION", "1.2.2"),
			collector.Environment("POETRY_VIRTUALENVS_IN_PROJECT", "true"),
		)
		c.Add(
			collector.Service("db", "mysql", "1024"),
			collector.Service("cache", "redis", "1024"),
		)
		collection := c.Collect()

		result, err := json.MarshalIndent(collection, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(result))
	},
}
