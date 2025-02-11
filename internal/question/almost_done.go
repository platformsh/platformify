package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
)

type AlmostDone struct{}

func (q *AlmostDone) Ask(ctx context.Context) error {
	out, _, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}
	answers, ok := models.FromContext(ctx)
	if !ok || answers.NoInteraction {
		return nil
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "                       (\\_/)"))
	fmt.Fprintln(out, colors.Colorize(colors.AccentCode, "We’re almost done...  =(^.^)="))
	fmt.Fprintln(out)
	return nil
}
