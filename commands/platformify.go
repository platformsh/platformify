package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/platformifier"
	"github.com/platformsh/platformify/vendorization"
)

type contextKey string

var FlavorKey contextKey = "flavor"

func NewPlatformifyCmd(assets *vendorization.VendorAssets) *cobra.Command {
	cmd := &cobra.Command{
		Use:           assets.Use,
		Aliases:       []string{"ify"},
		Short:         fmt.Sprintf("Creates starter YAML files for your %s project", assets.ServiceName),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			answers := models.NewAnswers()
			answers.Flavor, _ = cmd.Context().Value(FlavorKey).(string)
			ctx := models.ToContext(cmd.Context(), answers)
			ctx = colors.ToContext(
				ctx,
				cmd.OutOrStderr(),
				cmd.ErrOrStderr(),
			)
			q := questionnaire.New(
				&question.WorkingDirectory{},
				&question.FilesOverwrite{},
				&question.Welcome{},
				&question.Stack{},
				&question.Type{},
				&question.DependencyManager{},
				&question.Locations{},
				&question.Mounts{},
				&question.Name{},
				&question.ApplicationRoot{},
				&question.Environment{},
				&question.BuildSteps{},
				&question.DeployCommand{},
				&question.SocketFamily{},
				&question.WebCommand{},
				&question.AlmostDone{},
				&question.Services{},
			)
			err := q.AskQuestions(ctx)
			if errors.Is(err, questionnaire.ErrSilent) {
				return nil
			}

			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
				return err
			}

			input := answers.ToUserInput()

			pfier := platformifier.New(input, assets.ConfigFlavor)
			err = pfier.Platformify(ctx)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
				return fmt.Errorf("could not configure project: %w", err)
			}

			done := question.Done{}
			return done.Ask(ctx)
		},
	}

	return cmd
}
