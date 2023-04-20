package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
)

type Welcome struct{}

func (q *Welcome) Ask(ctx context.Context) error {
	out, _, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	fmt.Fprintln(out, colors.Colorize(colors.BrandCode, "Welcome to Platform.sh!"))
	fmt.Fprintln(out, colors.Colorize(colors.BrandCode, "Let's get started with a few questions."))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "We need to know a bit more about your project. This will only take a minute!")
	fmt.Fprintln(out)
	return nil
}
