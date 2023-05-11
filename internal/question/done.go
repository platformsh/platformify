package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/models"
)

type Done struct{}

func (q *Done) Ask(ctx context.Context) error {
	out, _, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "┌───────────────────────────────────────────────────┐")
	fmt.Fprintln(out, "│   CONGRATULATIONS!                                │")
	fmt.Fprintln(out, "│                                                   │")
	fmt.Fprintln(out, "│   We have created the following files for your:   │")
	fmt.Fprintln(out, "│     - .environment                                │")
	fmt.Fprintln(out, "│     - .platform.app.yaml                          │")
	fmt.Fprintln(out, "│     - .platform/services.yaml                     │")
	fmt.Fprintln(out, "│     - .platform/routes.yaml                       │")
	fmt.Fprintln(out, "│                                                   │")
	fmt.Fprintln(out, "│   We’re jumping for joy! ⍢                        │")
	fmt.Fprintln(out, "└───────────────────────────────────────────────────┘")
	fmt.Fprintln(out, "         │ /")
	fmt.Fprintln(out, "         │/")
	fmt.Fprintln(out, "         │")
	fmt.Fprintln(out, "  (\\ /)")
	fmt.Fprintln(out, "  ( . .)")
	fmt.Fprintln(out, "  o (_(“)(“)")
	fmt.Fprintln(out)
	if answers.HasGit {
		fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "You can now deploy your application to Platform.sh!"))
		fmt.Fprintln(
			out,
			colors.Colorize(
				colors.AccentCode,
				"To do so, commit your files and deploy your application using the Platform.sh CLI:",
			),
		)
		fmt.Fprintln(out, "  $ git add .")
		fmt.Fprintln(out, "  $ git commit -m 'Add Platform.sh configuration files'")
		fmt.Fprintln(out, "  $ platform project:set-remote")
		fmt.Fprintln(out, "  $ platform push")
		fmt.Fprintln(out)
		return nil
	}

	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "You can now deploy your application to Platform.sh!"))
	fmt.Fprintln(
		out,
		colors.Colorize(
			colors.AccentCode,
			//nolint:lll
			"To do so, you need to create a Git repository, commit your files and deploy your application using the Platform.sh CLI:",
		),
	)
	fmt.Fprintf(out, "  $ git init %s\n", answers.WorkingDirectory)
	fmt.Fprintln(out, "  $ git add .")
	fmt.Fprintln(out, "  $ git commit -m 'Add Platform.sh configuration files'")
	fmt.Fprintln(out, "  $ platform project:set-remote")
	fmt.Fprintln(out, "  $ platform push")
	fmt.Fprintln(out)

	return nil
}
