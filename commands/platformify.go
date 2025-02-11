package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

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
var NoInteractionKey contextKey = "no-interaction"

func NewPlatformifyCmd(assets *vendorization.VendorAssets) *cobra.Command {
	cmd := &cobra.Command{
		Use:           assets.Use,
		Aliases:       []string{"ify"},
		Short:         fmt.Sprintf("Creates starter YAML files for your %s project", assets.ServiceName),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Platformify(
				cmd.Context(),
				cmd.OutOrStderr(),
				cmd.ErrOrStderr(),
				assets,
			)
		},
	}

	return cmd
}

func Platformify(ctx context.Context, stdout, stderr io.Writer, assets *vendorization.VendorAssets) error {
	answers := models.NewAnswers()
	answers.Flavor, _ = ctx.Value(FlavorKey).(string)
	answers.NoInteraction, _ = ctx.Value(NoInteractionKey).(bool)
	ctx = models.ToContext(ctx, answers)
	ctx = colors.ToContext(
		ctx,
		stdout,
		stderr,
	)
	q := questionnaire.New(
		&question.WorkingDirectory{},
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
		fmt.Fprintln(stderr, colors.Colorize(colors.ErrorCode, err.Error()))
		return err
	}

	input := answers.ToUserInput()

	pfier := platformifier.New(input, assets.ConfigFlavor)
	configFiles, err := pfier.Platformify(ctx)
	if err != nil {
		fmt.Fprintln(stderr, colors.Colorize(colors.ErrorCode, err.Error()))
		return fmt.Errorf("could not configure project: %w", err)
	}

	filesToCreateUpdate := make([]string, 0, len(configFiles))
	for file := range configFiles {
		filesToCreateUpdate = append(filesToCreateUpdate, file)
	}

	filesOverwrite := question.FilesOverwrite{FilesToCreateUpdate: filesToCreateUpdate}
	if err := filesOverwrite.Ask(ctx); err != nil {
		return err
	}

	for file, contents := range configFiles {
		filePath := path.Join(answers.Cwd, file)
		if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|os.ModePerm); err != nil {
			fmt.Fprintf(stderr, "Could not create parent directories of file %s: %s\n", file, err)
			continue
		}

		if err := os.WriteFile(filePath, contents, 0o664); err != nil {
			fmt.Fprintf(stderr, "Could not write file %s: %s\n", file, err)
			continue
		}
	}

	done := question.Done{}
	return done.Ask(ctx)
}
