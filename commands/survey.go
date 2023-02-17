package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var surveyCmd = &cobra.Command{
	Use: "survey",
	Run: func(cmd *cobra.Command, args []string) {
		var languages = []string{
			"HTML and CSS",
			"Python",
			"Java",
			"JavaScript",
			"Swift",
			"C++",
			"C#",
			"R",
			"Golang (Go)",
			"PHP",
			"TypeScript",
			"Scala",
			"Shell",
			"PowerShell",
			"Perl",
			"Haskell",
			"Kotlin",
			"Visual Basic .NET",
			"SQL",
			"Delphi",
			"MATLAB",
			"Groovy",
			"Lua",
			"Rust",
			"Ruby",
			"C",
			"Dart",
			"DM",
		}

		// the questions to ask
		var qs = []*survey.Question{
			{
				Name: "language",
				Prompt: &survey.Select{
					Message: "Choose a programming language:",
					Options: languages,
					Default: nil,
				},
			},
			{
				Name: "php",
				Prompt: &survey.Confirm{
					Message: "Do you like PHP?",
					Default: false,
				},
			},
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: "What is your name?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
		}

		answers := struct {
			Language string
			PHP      bool
			Name     string
		}{}

		// perform the questions
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		switch answers.Language {
		case "PHP":
			if answers.PHP {
				fmt.Printf("Hey %s! You chose %s because you love PHP.\n", answers.Name, answers.Language)
			} else {
				fmt.Printf("Hey %s! Why did you choose %s if you don't like it?!\n", answers.Name, answers.Language)
			}
		default:
			if answers.PHP {
				fmt.Printf("Hey %s! Why did you choose %s if you love PHP?!\n", answers.Name, answers.Language)
			} else {
				fmt.Printf("Hey %s! You chose %s because you obviously don't like PHP.\n", answers.Name, answers.Language)
			}
		}
	},
}
