package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/vendorization"
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

	assets, _ := vendorization.FromContext(ctx)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "┌───────────────────────────────────────────────────┐")
	fmt.Fprintln(out, "│   CONGRATULATIONS!                                │")
	fmt.Fprintln(out, "│                                                   │")
	fmt.Fprintln(out, "│   We have created the following files for your:   │")
	for _, f := range assets.ProprietaryFiles {
		fmt.Fprintf(out, "│     - %-44s│\n", f)
	}
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
		fmt.Fprintln(out, colors.Colorize(colors.AccentCode, fmt.Sprintf(
			"You can now deploy your application to %s!",
			assets.ServiceName,
		)))
		fmt.Fprintln(
			out,
			colors.Colorize(
				colors.AccentCode,
				fmt.Sprintf(
					"To do so, commit your files and deploy your application using the %s CLI:",
					assets.ServiceName,
				),
			),
		)
		fmt.Fprintln(out, "  $ git add .")
		fmt.Fprintf(out, "  $ git commit -m 'Add %s configuration files'\n", assets.ServiceName)
		fmt.Fprintf(out, "  $ %s project:set-remote\n", assets.Binary)
		fmt.Fprintf(out, "  $ %s push\n", assets.Binary)
		fmt.Fprintln(out)
		return nil
	}

	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, fmt.Sprintf(
		"You can now deploy your application to %s!",
		assets.ServiceName,
	)))
	fmt.Fprintln(
		out,
		colors.Colorize(
			colors.AccentCode,
			fmt.Sprintf(
				"To do so, you need to create a Git repository, commit your files and deploy your application using the %s CLI:",
				assets.ServiceName,
			),
		),
	)
	fmt.Fprintf(out, "  $ git init %s\n", answers.WorkingDirectory)
	fmt.Fprintln(out, "  $ git add .")
	fmt.Fprintf(out, "  $ git commit -m 'Add %s configuration files'", assets.ServiceName)
	fmt.Fprintf(out, "  $ %s project:set-remote", assets.Binary)
	fmt.Fprintf(out, "  $ %s push", assets.Binary)
	fmt.Fprintln(out)

	return nil
}
