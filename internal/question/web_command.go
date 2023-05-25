package question

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

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

	// Do not ask the command for PHP applications
	if answers.Type.Runtime == models.PHP {
		return nil
	}

	switch answers.Stack {
	case models.Django:
		prefix := ""
		pythonPath := ""
		wsgi := "app.wsgi"
		// try to find the wsgi.py file to change the default command
		wsgiPath := utils.FindFile(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), "wsgi.py")
		if wsgiPath != "" {
			wsgiParentDir := path.Base(path.Dir(wsgiPath))
			wsgi = fmt.Sprintf("%s.wsgi", wsgiParentDir)

			// add the pythonpath if the wsgi.py file is not in the root of the app
			wsgiRel, _ := filepath.Rel(
				path.Join(answers.WorkingDirectory, answers.ApplicationRoot),
				path.Dir(path.Dir(wsgiPath)),
			)
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
		if answers.SocketFamily == models.TCP {
			answers.WebCommand = fmt.Sprintf("%sgunicorn %s -b 0.0.0.0:$PORT %s --log-file -", prefix, pythonPath, wsgi)
			return nil
		}

		answers.WebCommand = fmt.Sprintf("%sgunicorn %s -b unix:$UNIX %s --log-file -", prefix, pythonPath, wsgi)
		return nil
	case models.NextJS:
		answers.WebCommand = "npx next start -p $PORT"
		return nil
	case models.Strapi:
		switch answers.DependencyManager {
		case models.Yarn:
			answers.WebCommand = "NODE_ENV=production yarn start"
		default:
			answers.WebCommand = "NODE_ENV=production npm start"
		}
		return nil
	default:
		//nolint:lll
		answers.WebCommand = "echo 'Put your web server command in here! You need to listen to \"$UNIX\" unix socket. Read more about it here: https://docs.platform.sh/create-apps/app-reference.html#web-commands'; sleep 60"
		if answers.SocketFamily == models.TCP {
			//nolint:lll
			answers.WebCommand = "echo 'Put your web server command in here! You need to listen to \"$PORT\" port. Read more about it here: https://docs.platform.sh/create-apps/app-reference.html#web-commands'; sleep 60"
		}
		return nil
	}
}
