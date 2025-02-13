package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
var FSKey contextKey = "fs"

func NewPlatformifyCmd(assets *vendorization.VendorAssets) *cobra.Command {
	var noInteraction bool
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
				noInteraction,
				assets,
			)
		},
	}

	cmd.Flags().BoolVar(&noInteraction, "no-interaction", false, "Disable interactive prompts")
	return cmd
}

func Discover(
	ctx context.Context,
	flavor string,
	noInteraction bool,
	fileSystem fs.FS,
) (*platformifier.UserInput, error) {
	answers, _ := models.FromContext(ctx)
	if answers == nil {
		answers = models.NewAnswers()
		ctx = models.ToContext(ctx, answers)
	}
	answers.Flavor = flavor
	answers.NoInteraction = noInteraction
	answers.WorkingDirectory = fileSystem
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
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return answers.ToUserInput(), nil
}

func Platformify(
	ctx context.Context,
	stdout, stderr io.Writer,
	noInteraction bool,
	assets *vendorization.VendorAssets,
) error {
	ctx = colors.ToContext(ctx, stdout, stderr)
	ctx = models.ToContext(ctx, models.NewAnswers())
	input, err := Discover(ctx, assets.ConfigFlavor, noInteraction, nil)
	if err != nil {
		return err
	}
	answers, _ := models.FromContext(ctx)
	pfier := platformifier.New(input, assets.ConfigFlavor)
	configFiles, err := pfier.Platformify(ctx)
	if err != nil {
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
