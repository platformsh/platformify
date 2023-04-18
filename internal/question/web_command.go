package question

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type WebCommand struct{}

func (q *WebCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.WebCommand != "" {
		// Skip the step
		return nil
	}

	question := &survey.Input{Message: "Web command:"}
	cwd, _ := os.Getwd()
	switch answers.Stack {
	case models.Django:
		prefix := ""
		pythonPath := ""
		wsgi := "app.wsgi"
		// try to find the wsgi.py file to change the default command
		if wsgiPath := utils.FindFile(path.Join(cwd, answers.ApplicationRoot), "wsgi.py"); wsgiPath != "" {
			wsgiParentDir := path.Base(path.Dir(wsgiPath))
			wsgi = fmt.Sprintf("%s.wsgi", wsgiParentDir)

			// add the pythonpath if the wsgi.py file is not in the root of the app
			wsgiRel, _ := filepath.Rel(path.Join(cwd, answers.ApplicationRoot), path.Dir(path.Dir(wsgiPath)))
			if wsgiRel != "" {
				pythonPath = "--pythonpath=" + path.Base(path.Dir(path.Dir(wsgiPath)))
			}
		}

		switch answers.DependencyManager {
		case models.Pipenv:
			prefix = "pipenv run "
		case models.Poetry:
			prefix = "poetry run "
		}
		if answers.ListenInterface == models.HTTP {
			question.Default = fmt.Sprintf("%sgunicorn %s -b 0.0.0.0:$PORT %s --log-file -", prefix, pythonPath, wsgi)
		} else {
			question.Default = fmt.Sprintf("%sgunicorn %s -b unix:$UNIX %s --log-file -", prefix, pythonPath, wsgi)
		}
	case models.NextJS:
		if answers.ListenInterface == models.HTTP {
			question.Default = "npx next start -p $PORT"
		}
	}

	var webCommand string
	err := survey.AskOne(question, &webCommand)
	if err != nil {
		return err
	}

	answers.WebCommand = webCommand

	return nil
}
