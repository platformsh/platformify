package models

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const (
	ChromeHeadless  ServiceName = "chrome-headless"
	InfluxDB        ServiceName = "influxdb"
	Kafka           ServiceName = "kafka"
	MariaDB         ServiceName = "mariadb"
	Memcached       ServiceName = "memcached"
	MySQL           ServiceName = "mysql"
	NetworkStorage  ServiceName = "network-storage"
	OpenSearch      ServiceName = "opensearch"
	OracleMySQL     ServiceName = "oracle-mysql"
	PostgreSQL      ServiceName = "postgresql"
	RabbitMQ        ServiceName = "rabbitmq"
	Redis           ServiceName = "redis"
	RedisPersistent ServiceName = "redis-persistent"
	Solr            ServiceName = "solr"
	Varnish         ServiceName = "varnish"
	VaultKMS        ServiceName = "vault-kms"
)

var (
	ServiceNames = ServiceNameList{
		ChromeHeadless,
		InfluxDB,
		Kafka,
		MariaDB,
		Memcached,
		MySQL,
		NetworkStorage,
		OpenSearch,
		OracleMySQL,
		PostgreSQL,
		RabbitMQ,
		Redis,
		RedisPersistent,
		Solr,
		Varnish,
		VaultKMS,
	}
)

type ServiceName string

func (s ServiceName) String() string {
	return string(s)
}

func (s ServiceName) Title() string {
	switch s {
	case ChromeHeadless:
		return "Chrome Headless"
	case InfluxDB:
		return "InfluxDB"
	case Kafka:
		return "Kafka"
	case MariaDB:
		return "MariaDB"
	case Memcached:
		return "Memcached"
	case MySQL:
		return "MySQL"
	case NetworkStorage:
		return "Network Storage"
	case OpenSearch:
		return "OpenSearch"
	case OracleMySQL:
		return "Oracle MySQL"
	case PostgreSQL:
		return "PostgreSQL"
	case RabbitMQ:
		return "RabbitMQ"
	case Redis:
		return "Redis"
	case RedisPersistent:
		return "Redis Persistent"
	case Solr:
		return "Solr"
	case Varnish:
		return "Varnish"
	case VaultKMS:
		return "Vault KMS"
	default:
		return ""
	}
}

func (s ServiceName) IsPersistent() bool {
	switch s {
	case ChromeHeadless, Memcached, Redis, RedisPersistent:
		return false
	default:
		return true
	}
}

type ServiceNameList []ServiceName

func (s *ServiceNameList) WriteAnswer(_ string, value interface{}) error {
	switch answer := value.(type) {
	case []survey.OptionAnswer: // MultiSelect
		for _, item := range answer {
			service, err := ServiceNames.serviceByTitle(item.Value)
			if err != nil {
				return err
			}
			*s = append(*s, service)
		}
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (s *ServiceNameList) AllTitles() []string {
	titles := make([]string, 0, len(*s))
	for _, service := range *s {
		titles = append(titles, service.Title())
	}
	return titles
}

func (s *ServiceNameList) serviceByTitle(title string) (ServiceName, error) {
	for _, service := range *s {
		if service.Title() == title {
			return service, nil
		}
	}
	return "", fmt.Errorf("service name by title is not found")
}
