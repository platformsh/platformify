package question

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
)

type Name struct{}

var (
	invalidChars      = regexp.MustCompile("[^a-z0-9_-]")
	consecutiveDashes = regexp.MustCompile("-+")
)

func (q *Name) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.Name != "" {
		// Skip the step
		return nil
	}

	question := &survey.Input{
		Message: "Tell us your project's application name:", Default: slugify(path.Base(answers.Cwd)),
	}

	var name string
	err := survey.AskOne(question, &name, survey.WithValidator(survey.Required), survey.WithValidator(validSlug))
	if err != nil {
		return err
	}

	answers.Name = name

	return nil
}

func slugify(s string) string {
	return consecutiveDashes.ReplaceAllLiteralString(
		invalidChars.ReplaceAllLiteralString(strings.ToLower(
			strings.TrimSpace(s),
		), "-"), "-")
}

func validSlug(val interface{}) error {
	if val.(string) != slugify(val.(string)) {
		return fmt.Errorf(
			"%s: the name can only contain lowercase alphanumeric characters, dashes, or underscores",
			val.(string),
		)
	}

	return nil
}
