package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/question"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/platformifiers"
)

// PlatformifyCmd represents the base Platformify command when called without any subcommands
var PlatformifyCmd = &cobra.Command{
	Use:   "platformify",
	Short: "Platfomrify your application, and deploy it to the Platform.sh",
	Long: `Platformify your application, creating all the needed files
for it to be deployed to Platform.sh.

This will create the needed YAML files for both your application and your
services, choosing from a variety of stacks or simple runtimes.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := models.NewAnswers()
		ctx := models.ToContext(cmd.Context(), answers)
		q := questionnaire.New(
			&question.WorkingDirectory{},
			&question.Stack{},
			&question.Type{},
			&question.DependencyManager{},
			&question.Name{},
			&question.ApplicationRoot{},
			&question.Environment{},
			&question.BuildSteps{},
			&question.DeployCommand{},
			&question.ListenInterface{},
			&question.WebCommand{},
			&question.Services{},
		)
		err := q.AskQuestions(ctx)
		if err != nil {
			return err
		}

		pfier, err := platformifiers.GetPlatformifier(answers)
		if err != nil {
			return fmt.Errorf("creating platformifier failed: %s", err)
		}

		if err := pfier.Platformify(ctx); err != nil {
			return fmt.Errorf("could not platformify project: %w", err)
		}

		return nil
	},
}
