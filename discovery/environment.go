package discovery

// Returns the application environment, either from memory or by discovering it on the spot
func (d *Discoverer) Environment() (map[string]string, error) {
	if environment, ok := d.memory["environment"]; ok {
		return environment.(map[string]string), nil
	}

	environment, err := d.discoverEnvironment()
	if err != nil {
		return nil, err
	}

	d.memory["environment"] = environment
	return environment, nil
}

func (d *Discoverer) discoverEnvironment() (map[string]string, error) {
	depManagers, err := d.DependencyManagers()
	if err != nil {
		return nil, err
	}

	typ, err := d.Type()
	if err != nil {
		return nil, err
	}

	environment := make(map[string]string)

	for _, dm := range depManagers {
		switch dm {
		case "poetry":
			environment["POETRY_VERSION"] = "1.4.0"
			environment["POETRY_VIRTUALENVS_IN_PROJECT"] = "true"
		case "pipenv":
			environment["PIPENV_TOOL_VERSION"] = "2023.2.18"
			environment["PIPENV_VENV_IN_PROJECT"] = "1"
		case "npm", "yarn":
			if typ != "nodejs" {
				environment["N_PREFIX"] = "/app/.global"
			}
		}
	}

	return environment, nil
}
