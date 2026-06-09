package question

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/vendorization"
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
	if answers.Type.Runtime.Type == "php" {
		return nil
	}

	assets, _ := vendorization.FromContext(ctx)

	//nolint:lll
	answers.WebCommand = fmt.Sprintf(
		"echo 'Put your web server command in here! You need to listen to \"$UNIX\" unix socket. Read more about it here: %s#web-commands'; sleep 60",
		assets.Docs().AppReference,
	)

	if answers.SocketFamily == models.TCP {
		//nolint:lll
		answers.WebCommand = fmt.Sprintf(
			"echo 'Put your web server command in here! You need to listen to \"$PORT\" port. Read more about it here: %s#web-commands'; sleep 60",

			assets.Docs().AppReference,
		)
	}

	switch answers.Stack {
	case models.Django:
		prefix := ""
		pythonPath := ""
		wsgi := "app.wsgi"
		// try to find the wsgi.py file to change the default command
		wsgiPath := utils.FindFile(answers.WorkingDirectory, answers.ApplicationRoot, "wsgi.py")
		if wsgiPath != "" {
			wsgiParentDir := path.Base(path.Dir(wsgiPath))
			wsgi = fmt.Sprintf("%s.wsgi", wsgiParentDir)

			// add the pythonpath if the wsgi.py file is not in the root of the app
			wsgiRel, _ := filepath.Rel(
				answers.ApplicationRoot,
				path.Dir(path.Dir(wsgiPath)),
			)
			if wsgiRel != "." {
				pythonPath = "--pythonpath=" + path.Base(path.Dir(path.Dir(wsgiPath)))
			}
		}

		if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			prefix = "pipenv run "
		} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			prefix = "poetry run "
		}
		if answers.SocketFamily == models.TCP {
			answers.WebCommand = fmt.Sprintf("%sgunicorn %s -b 0.0.0.0:$PORT %s --log-file -", prefix, pythonPath, wsgi)
			return nil
		}

		answers.WebCommand = fmt.Sprintf("%sgunicorn %s -b unix:$SOCKET %s --log-file -", prefix, pythonPath, wsgi)
		return nil
	case models.Rails:
		if answers.SocketFamily == models.TCP {
			answers.WebCommand = "bundle exec rails server"
			return nil
		}

		answers.WebCommand = "bundle exec puma -b unix://$SOCKET"
		return nil
	case models.NextJS:
		answers.WebCommand = "npx next start -p $PORT"
		return nil
	case models.Strapi:
		if _, ok := utils.GetJSONValue(
			answers.WorkingDirectory,
			[]string{"scripts", "start"},
			"package.json",
			true,
		); ok {
			if slices.Contains(answers.DependencyManagers, models.Yarn) {
				answers.WebCommand = "NODE_ENV=production yarn start"
			} else {
				answers.WebCommand = "NODE_ENV=production npm start"
			}
		}
	case models.Express:
		if _, ok := utils.GetJSONValue(
			answers.WorkingDirectory,
			[]string{"scripts", "start"},
			"package.json",
			true,
		); ok {
			if slices.Contains(answers.DependencyManagers, models.Yarn) {
				answers.WebCommand = "NODE_ENV=production yarn start"
			} else {
				answers.WebCommand = "NODE_ENV=production npm start"
			}
			return nil
		}

		if mainPath, ok := utils.GetJSONValue(
			answers.WorkingDirectory,
			[]string{"main"},
			"package.json",
			true,
		); ok {
			answers.WebCommand = fmt.Sprintf("node %s", mainPath.(string))
			return nil
		}

		if indexFile := utils.FindFile(answers.WorkingDirectory, "", "index.js"); indexFile != "" {
			answers.WebCommand = fmt.Sprintf("node %s", indexFile)
			return nil
		}
	case models.Flask:
		appPath := ""
		// try to find the app.py, api.py or server.py files
		for _, name := range []string{"app.py", "server.py", "api.py"} {
			if _, err := fs.Stat(answers.WorkingDirectory, name); err == nil {
				appPath = fmt.Sprintf("'%s:app'", strings.TrimSuffix(name, ".py"))
				break
			}
		}
		if appPath == "" {
			return nil
		}

		prefix := ""
		if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			prefix = "pipenv run "
		} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			prefix = "poetry run "
		}

		if answers.SocketFamily == models.TCP {
			answers.WebCommand = fmt.Sprintf("%sgunicorn -b 0.0.0.0:$PORT %s --log-file -", prefix, appPath)
			return nil
		}

		answers.WebCommand = fmt.Sprintf("%sgunicorn -b unix:$SOCKET %s --log-file -", prefix, appPath)
		return nil
	}

	return nil
}
