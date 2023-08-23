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
		Message: "Select all the services you are using",
		Options: models.ServiceNames.AllTitles(),
	}

	var services models.ServiceNameList
	for {
		if err := survey.AskOne(question, &services, survey.WithKeepFilter(true)); err != nil {
			return err
		}

		if len(services) > 0 {
			break
		}

		confirmQuestion := &survey.Confirm{
			Message: "You have not selected any service, would you like to proceed anyway?",
			Default: false,
		}
		proceed := false
		if err := survey.AskOne(confirmQuestion, &proceed); err != nil {
			return err
		}

		if proceed {
			break
		}
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
