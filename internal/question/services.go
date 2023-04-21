package question

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
)

type Services struct{}

func (q *Services) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if len(answers.Services) != 0 {
		// Skip the step
		return nil
	}

	out, _, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}
	fmt.Fprintln(
		out,
		colors.Colorize(
			colors.AccentCode,
			"Last but not least, unless you’re creating a static website, your project uses services. Let’s define them:",
		),
	)
	fmt.Fprintln(out)

	question := &survey.MultiSelect{
		Message: "Which services are you using?",
		Options: models.ServiceNames.AllTitles(),
	}
	var services models.ServiceNameList
	if err := survey.AskOne(question, &services, survey.WithKeepFilter(true)); err != nil {
		return err
	}

	for _, serviceName := range services {
		versions, ok := models.ServiceTypeVersions[serviceName]
		if !ok || len(versions) == 0 {
			return nil
		}

		service := models.Service{
			Name: strings.ReplaceAll(serviceName.String(), "-", "_"),
			Type: models.ServiceType{
				Name:    serviceName.String(),
				Version: versions[0],
			},
			TypeVersions: versions,
		}
		if serviceName.IsPersistent() {
			service.Disk = models.ServiceDisks[0]
			service.DiskSizes = models.ServiceDisks
		}

		answers.Services = append(answers.Services, service)
	}

	return nil
}
