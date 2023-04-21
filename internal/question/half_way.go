package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
)

type HalfWay struct{}

func (q *HalfWay) Ask(ctx context.Context) error {
	out, _, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "                             (\\ /)"))
	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "Hurrah! We are halfway there (-‿◦)"))
	fmt.Fprintln(out)
	return nil
}
