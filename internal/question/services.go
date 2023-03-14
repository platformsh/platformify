package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Services struct{}

func (q *Services) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	serviceTypes := []string{
		"chrome-headless",
		"elasticsearch",
		"influxdb",
		"kafka",
		"mariadb",
		"memcached",
		"mongodb",
		"mongodb-enterprise",
		"mysql",
		"network-storage",
		"opensearch",
		"oracle-mysql",
		"postgresql",
		"rabbitmq",
		"redis",
		"redis-persistent",
		"solr",
		"varnish",
		"vault-kms",
	}
	serviceDisks := []string{
		"1024",
		"2048",
		"3072",
		"4096",
		"5120",
	}

	for {
		var addService = true
		var question survey.Prompt
		question = &survey.Confirm{
			Message: "Would you like to add a service?",
			Default: true,
		}

		err := survey.AskOne(question, &addService)
		if err != nil {
			return err
		}

		if !addService {
			break
		}

		question = &survey.Input{Message: "Service name:"}

		var serviceName string
		err = survey.AskOne(question, &serviceName)
		if err != nil {
			return err
		}

		question = &survey.Select{
			Message: "Choose service type:",
			Options: serviceTypes,
		}

		var serviceTypeName string
		err = survey.AskOne(question, &serviceTypeName)
		if err != nil {
			return err
		}

		question = nil
		switch serviceTypeName {
		case "chrome-headless":
			question = &survey.Select{
				Message: "Choose Headless Chrome version:",
				Options: []string{
					"95", "91", "86", "84", "83", "81", "80", "73",
				},
				Default: "95",
			}
		case "elasticsearch": // no versions
		case "influxdb": // only one version
			question = &survey.Select{
				Message: "Choose InfluxDB version:",
				Options: []string{
					"2.3",
				},
				Default: "2.3",
			}
		case "kafka": // only one version
			question = &survey.Select{
				Message: "Choose Kafka version:",
				Options: []string{
					"3.2",
				},
				Default: "3.2",
			}
		case "mariadb":
			question = &survey.Select{
				Message: "Choose MariaDB/MySQL version:",
				Options: []string{
					"10.6", "10.5", "10.4", "10.3",
				},
				Default: "10.6",
			}
		case "memcached":
			question = &survey.Select{
				Message: "Memcached version:",
				Options: []string{
					"1.6", "1.5", "1.4",
				},
				Default: "1.6",
			}
		case "mongodb": // no versions
		case "mongodb-enterprise":
			question = &survey.Select{
				Message: "Choose MongoDB version:",
				Options: []string{
					"5.0", "4.4", "4.2",
				},
				Default: "5.0",
			}
		case "mysql":
			question = &survey.Select{
				Message: "Choose MariaDB/MySQL version:",
				Options: []string{
					"10.6", "10.5", "10.4", "10.3",
				},
				Default: "10.6",
			}
		case "network-storage": // only one version
			question = &survey.Select{
				Message: "Choose Network Storage version:",
				Options: []string{
					"2.0",
				},
				Default: "2.0",
			}
		case "opensearch":
			question = &survey.Select{
				Message: "Choose OpenSearch version:",
				Options: []string{
					"2", "1.2", "1.1",
				},
				Default: "2",
			}
		case "oracle-mysql":
			question = &survey.Select{
				Message: "Choose Oracle MySQL version:",
				Options: []string{
					"8.0", "5.7",
				},
				Default: "8.0",
			}
		case "postgresql":
			question = &survey.Select{
				Message: "Choose PostgreSQL version:",
				Options: []string{
					"15", "14", "13", "12", "11",
				},
				Default: "15",
			}
		case "rabbitmq":
			question = &survey.Select{
				Message: "Choose RabbitMQ version:",
				Options: []string{
					"3.11", "3.10", "3.9",
				},
				Default: "3.11",
			}
		case "redis":
			question = &survey.Select{
				Message: "Choose Redis version:",
				Options: []string{
					"7.0", "6.2",
				},
				Default: "7.0",
			}
		case "redis-persistent":
			question = &survey.Select{
				Message: "Choose Persistent Redis version:",
				Options: []string{
					"7.0", "6.2",
				},
				Default: "7.0",
			}
		case "solr":
			question = &survey.Select{
				Message: "Choose Solr version:",
				Options: []string{
					"9.1", "8.11",
				},
				Default: "9.1",
			}
		case "varnish":
			question = &survey.Select{
				Message: "Choose Varnish version:",
				Options: []string{
					"7.2", "7.1", "6.3", "6.0",
				},
				Default: "7.2",
			}
		case "vault-kms": // only one version
			question = &survey.Select{
				Message: "Choose Vault KMS version:",
				Options: []string{
					"1.12",
				},
				Default: "1.12",
			}
		}

		var serviceTypeVersion string
		if question != nil {
			err = survey.AskOne(question, &serviceTypeVersion)
			if err != nil {
				return err
			}
		}

		question = &survey.Select{
			Message: "Choose service disk space:",
			Options: serviceDisks,
			Default: "1024",
		}

		var serviceDisk string
		err = survey.AskOne(question, &serviceDisk)
		if err != nil {
			return err
		}

		answers.Services = append(answers.Services, models.Service{
			Name: serviceName,
			Type: models.ServiceType{
				Name:    serviceTypeName,
				Version: serviceTypeVersion,
			},
			Disk: serviceDisk,
		})
	}

	return nil
}
