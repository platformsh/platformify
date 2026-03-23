//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

const registryURL = "https://raw.githubusercontent.com/platformsh/platformsh-docs/main/shared/data/registry.json"

type Registry map[string]Service

type Service struct {
	Description string              `json:"description"`
	Disk        bool                `json:"disk"`
	Runtime     bool                `json:"runtime"`
	Type        string              `json:"type"`
	Versions    map[string][]string `json:"versions"`
}

// Mapping from registry service names to our ServiceName constants
var serviceMapping = map[string]string{
	"chrome-headless": "ChromeHeadless",
	"influxdb":        "InfluxDB",
	"kafka":           "Kafka",
	"mariadb":         "MariaDB",
	"memcached":       "Memcached",
	"mysql":           "MySQL",
	"network-storage": "NetworkStorage",
	"opensearch":      "OpenSearch",
	"oracle-mysql":    "OracleMySQL",
	"postgresql":      "PostgreSQL",
	"rabbitmq":        "RabbitMQ",
	"redis":           "Redis",
	"solr":            "Solr",
	"varnish":         "Varnish",
	"vault-kms":       "VaultKMS",
}

// Mapping from registry runtime names to our Runtime constants
var runtimeMapping = map[string]string{
	"dotnet": "DotNet",
	"elixir": "Elixir",
	"golang": "Golang",
	"java":   "Java",
	"nodejs": "NodeJS",
	"php":    "PHP",
	"python": "Python",
	"ruby":   "Ruby",
	"rust":   "Rust",
}

const versionTemplate = `//go:generate go run generate_versions.go

package models

var (
	LanguageTypeVersions = map[Runtime][]string{
{{- range .Languages }}
		{{ .Name }}: {{ .Versions }},
{{- end }}
	}

	ServiceTypeVersions = map[ServiceName][]string{
{{- range .Services }}
		{{ .Name }}: {{ .Versions }},
{{- end }}
	}
)

func DefaultVersionForRuntime(r Runtime) string {
	versions := LanguageTypeVersions[r]
	if len(versions) == 0 {
		return ""
	}
	return versions[0]
}
`

type TemplateData struct {
	Languages []TypeVersion
	Services  []TypeVersion
}

type TypeVersion struct {
	Name     string
	Versions string
}

func main() {
	// Fetch the registry data
	resp, err := http.Get(registryURL)
	if err != nil {
		log.Fatalf("Failed to fetch registry: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var registry Registry
	if err := json.Unmarshal(body, &registry); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	var languages []TypeVersion
	var services []TypeVersion

	// Process runtimes (languages)
	for serviceName, service := range registry {
		if !service.Runtime {
			continue
		}

		constantName, exists := runtimeMapping[serviceName]
		if !exists {
			continue // Skip runtimes we don't support
		}

		supported := service.Versions["supported"]
		if len(supported) == 0 {
			continue
		}

		versions := formatVersionSlice(supported)
		languages = append(languages, TypeVersion{
			Name:     constantName,
			Versions: versions,
		})
	}

	// Process services
	for serviceName, service := range registry {
		if service.Runtime {
			continue
		}

		constantName, exists := serviceMapping[serviceName]
		if !exists {
			continue // Skip services we don't support
		}

		supported := service.Versions["supported"]
		if len(supported) == 0 {
			continue
		}

		versions := formatVersionSlice(supported)
		services = append(services, TypeVersion{
			Name:     constantName,
			Versions: versions,
		})

		// Add RedisPersistent right after Redis
		if serviceName == "redis" {
			services = append(services, TypeVersion{
				Name:     "RedisPersistent",
				Versions: versions,
			})
		}
	}

	// Sort for consistent output
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Name < languages[j].Name
	})
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	// Generate the Go code
	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	data := TemplateData{
		Languages: languages,
		Services:  services,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// Format the generated Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Failed to format Go code: %v", err)
	}

	// Write to version.go file
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	outputPath := filepath.Join(dir, "version.go")
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		log.Fatalf("Failed to write version.go: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", outputPath)
}

func formatVersionSlice(versions []string) string {
	if len(versions) == 0 {
		return "{}"
	}

	// Sort versions from latest to oldest
	sortedVersions := make([]string, len(versions))
	copy(sortedVersions, versions)

	sort.Slice(sortedVersions, func(i, j int) bool {
		return compareVersions(sortedVersions[i], sortedVersions[j])
	})

	quoted := make([]string, len(sortedVersions))
	for i, v := range sortedVersions {
		quoted[i] = fmt.Sprintf(`"%s"`, v)
	}

	return fmt.Sprintf("{%s}", strings.Join(quoted, ", "))
}

// compareVersions returns true if version a is greater than version b
func compareVersions(a, b string) bool {
	// Parse version parts
	aParts := parseVersion(a)
	bParts := parseVersion(b)

	// Compare each part
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		aVal := 0
		bVal := 0

		if i < len(aParts) {
			aVal = aParts[i]
		}
		if i < len(bParts) {
			bVal = bParts[i]
		}

		if aVal > bVal {
			return true
		}
		if aVal < bVal {
			return false
		}
	}

	return false // versions are equal
}

// parseVersion parses a version string into numeric parts
func parseVersion(version string) []int {
	parts := strings.Split(version, ".")
	nums := make([]int, 0, len(parts))

	for _, part := range parts {
		// Handle cases like "25.05" or "1.0" or just "1"
		num := 0
		fmt.Sscanf(part, "%d", &num)
		nums = append(nums, num)
	}

	return nums
}
