package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Services Question

func (q *Services) Ask(ctx context.Context) error {
	var addService = true
	question := &survey.Confirm{
		Message: "Would you like to add a service?",
		Default: true,
	}

	err := survey.AskOne(question, &addService)
	if err != nil {
		return err
	}

	if !addService {
		return nil
	}

	serviceTypes := []string{
		"postgresql:14",
		"mariadb:11",
		"redis:2",
	}
	serviceDisks := []string{
		"1024Mb",
		"2048Mb",
		"3072Mb",
		"4096Mb",
		"5120Mb",
	}

	for addService {
		// the questions to ask
		var qs = []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "Service name"},
				Validate: survey.Required,
			},
			{
				Name: "type",
				Prompt: &survey.Select{
					Message: "Choose service type:",
					Options: serviceTypes,
					Default: nil,
				},
			},
			{
				Name: "disk",
				Prompt: &survey.Select{
					Message: "Choose service disk:",
					Options: serviceDisks,
					Default: "1024Mb",
				},
			},
		}

		service := struct {
			Name string
			Type string
			Disk string
		}{}

		err = survey.Ask(qs, &service)
		if err != nil {
			return err
		}

		q.Answers.Services = append(q.Answers.Services, Service{
			Name: service.Name,
			Type: service.Type,
			Disk: service.Disk,
		})

		question := &survey.Confirm{
			Message: "Would you like to add one more service?",
			Default: true,
		}

		err := survey.AskOne(question, &addService)
		if err != nil {
			return err
		}
	}

	return nil
}
