package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/question"
	"github.com/platformsh/platformify/internal/questionnaire"
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
		ctx := context.TODO()
		answers := question.NewAnswers()
		q := questionnaire.New(
			&question.Stack{Answers: answers},
			&question.Type{Answers: answers},
			&question.Name{Answers: answers},
			&question.Root{Answers: answers},
			&question.Environment{Answers: answers},
			&question.BuildSteps{Answers: answers},
			&question.WebCommand{Answers: answers},
			&question.Listen{Answers: answers},
			&question.DeployCommand{Answers: answers},
			&question.DependencyManager{Answers: answers},
			&question.Services{Answers: answers},
		)
		err := q.AskQuestions(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result, err := json.MarshalIndent(answers, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(result))
	},
}
