package models

var (
	LanguageTypeVersions = map[Runtime][]string{
		DotNet: {"6.0", "3.1"},
		Elixir: {"1.13", "1.12", "1.11"},
		Golang: {"1.22", "1.21", "1.20"},
		Java:   {"19", "18", "17", "11", "8"},
		NodeJS: {"20", "18", "16"},
		PHP:    {"8.2", "8.1", "8.0"},
		Python: {"3.11", "3.10", "3.9", "3.8", "3.7"},
		Ruby:   {"3.3", "3.2", "3.1", "3.0", "2.7", "2.6", "2.5", "2.4", "2.3"},
		Rust:   {"1"},
	}

	ServiceTypeVersions = map[ServiceName][]string{
		ChromeHeadless:  {"95", "91", "86", "84", "83", "81", "80", "73"},
		InfluxDB:        {"2.3"},
		Kafka:           {"3.2"},
		MariaDB:         {"11.0", "10.11", "10.6", "10.5", "10.4", "10.3"},
		Memcached:       {"1.6", "1.5", "1.4"},
		MySQL:           {"10.6", "10.5", "10.4", "10.3"},
		NetworkStorage:  {"2.0"},
		OpenSearch:      {"2", "1.2", "1.1"},
		OracleMySQL:     {"8.0", "5.7"},
		PostgreSQL:      {"15", "14", "13", "12", "11"},
		RabbitMQ:        {"3.11", "3.10", "3.9"},
		Redis:           {"7.0", "6.2"},
		RedisPersistent: {"7.0", "6.2"},
		Solr:            {"9.1", "8.11"},
		Varnish:         {"7.2", "7.1", "6.3", "6.0"},
		VaultKMS:        {"1.12"},
	}
)

func DefaultVersionForRuntime(r Runtime) string {
	versions := LanguageTypeVersions[r]
	if len(versions) == 0 {
		return ""
	}
	return versions[0]
}
