//go:generate go run generate_versions.go

package models

var (
	LanguageTypeVersions = map[Runtime][]string{
		DotNet: {"8.0", "7.0", "6.0"},
		Elixir: {"1.18", "1.15", "1.14"},
		Golang: {"1.25", "1.24", "1.23", "1.22", "1.21", "1.20"},
		Java:   {"21", "19", "18", "17", "11", "8"},
		NodeJS: {"24", "22", "20"},
		PHP:    {"8.4", "8.3", "8.2", "8.1"},
		Python: {"3.13", "3.12", "3.11", "3.10", "3.9", "3.8"},
		Ruby:   {"3.4", "3.3", "3.2", "3.1", "3.0"},
		Rust:   {"1"},
	}

	ServiceTypeVersions = map[ServiceName][]string{
		ChromeHeadless:  {"120", "113", "95", "91"},
		ClickHouse:      {"25.3", "24.3", "23.8"},
		InfluxDB:         {"2.7", "2.3"},
		Kafka:           {"3.7", "3.6", "3.4", "3.2"},
		MariaDB:         {"11.8", "11.4", "10.11", "10.6"},
		Memcached:       {"1.6", "1.5", "1.4"},
		MySQL:           {"11.8", "11.4", "10.11", "10.6"},
		NetworkStorage:  {"1.0"},
		OpenSearch:      {"3", "2"},
		OracleMySQL:     {"8.0", "5.7"},
		PostgreSQL:      {"18", "17", "16", "15", "14", "13", "12"},
		RabbitMQ:        {"4.1", "4.0", "3.13", "3.12"},
		Redis:           {"8.0", "7.2"},
		RedisPersistent: {"8.0", "7.2"},
		Solr:            {"9.9", "9.6", "9.4", "9.2", "9.1", "8.11"},
		Varnish:         {"7.6", "7.3", "7.2", "6.0"},
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
