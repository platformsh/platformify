package models

// AppConfig is the data model for app configurations found in .platform.app.yaml.
type AppConfig struct {
	// Name	A unique name for the app. Must be lowercase alphanumeric characters.
	Name string `yaml:"name"`
	// Type	The base image to use with a specific app language. Format: runtime:version.
	Type AppType `yaml:"type"`
	// Size defines the amount of resources to devote to the app. Defaults to AUTO in production environments.
	Size string `yaml:"size,omitempty"`
	// The Disk space for the app in MB. Minimum value is 128. Defaults to null, meaning no disk is available.
	Disk int `yaml:"disk,omitempty"`
	// A Build flavor which determines what happens when the app is built.
	Build *BuildConfig `yaml:"build,omitempty"`
	// Relationships	A dictionary of relationships and their connections to other services and apps.
	Relationships Relationships `yaml:"relationships,omitempty"`
	// Variables	A dictionary of variables to control the environment.
	Variables map[string]string `yaml:"variables,omitempty"`
	// Hooks A dictionary of commands run at different stages in the build and deploy process.
	Hooks Hooks `yaml:"hooks,omitempty"`
	// A Dependencies dictionary of global dependencies to install before the build hook is run.
	Dependencies Dependencies `yaml:"dependencies,omitempty"`
	// A Mounts dictionary of writeable directories available after the app is built.
	// If set as a local source, disk is required.
	Mounts []Mount `yaml:"mounts,omitempty"`
	// A Web instance defining how the web application is served.
	Web *WebConfig `yaml:"web,omitempty"`
	// A Crons dictionary of scheduled tasks for the app.
	Crons *CronConfig `yaml:"cron,omitempty"`
	// Workers are alternate copies of the application to run as background processes.
	Workers *WorkerConfig `yaml:"worker,omitempty"`
	// The Timezone for crons to run. Format: a TZ database name.
	Timezone string `yaml:"timezone,omitempty"`
	// An Access dictionary	defining access control for roles accessing app environments via ssh.
	Access *AccessConfig `yaml:"access,omitempty"`
	// A Firewall dictionary	of outbound rules for the application.
	Firewall *FirewallConfig `yaml:"firewall,omitempty"`
	// A Source dictionary	with the app’s source code and operations that can be run on it.
	Source *SourceConfig `yaml:"source,omitempty"`
	// Runtime customizations to your PHP or Lisp.
	Runtime *RuntimeConfig `yaml:"runtime,omitempty"`
	// Additional hosts dictionary mapping hostnames to IPs.
	AdditionalHosts *AdditionalHosts `yaml:"additional_hosts,omitempty"`
}

type AccessConfig struct {
	SSH string `yaml:"ssh,omitempty"`
}

type AdditionalHosts map[string]string

type AppType map[string]string

func NewAppType(runtime string, version string) (AppType, error) {
	appType := AppType{}
	appType[runtime] = version
	return appType, nil
}

type BuildConfig struct {
	Flavor string `yaml:"flavor,omitempty"`
}

type CronConfig struct {
	Jobs map[string]string `yaml:"jobs,omitempty"`
}

type Dependencies map[string]string

type FirewallConfig struct {
	Rules []FirewallRule `yaml:"outbound,omitempty"`
}

type FirewallRule struct {
	Domains  []string `yaml:"domains,omitempty,inline"`
	IPs      []string `yaml:"ips,omitempty,inline"`
	Ports    []int    `yaml:"ports,omitempty,inline"`
	Protocol string   `yaml:"protocol,omitempty"`
}

type Hooks map[string]string

type Mount struct {
	Source   string `yaml:"source"`
	Target   string `yaml:"target"`
	ReadOnly bool   `yaml:"readonly,omitempty"`
}

type Relationships map[string]string

type RuntimeConfig struct {
	Extensions         []string      `yaml:"extensions,omitempty"`
	DisabledExtensions []string      `yaml:"disabled_extensions,omitempty"`
	RequestTimeout     int           `yaml:"request_terminate_timeout,omitempty"`
	SizingHints        SizingHints   `yaml:"sizing_hints,omitempty"`
	XDebug             RuntimeXDebug `yaml:"xdebug,omitempty"`
}

type RuntimeXDebug struct {
	IDEKey string `yaml:"idekey,omitempty"`
}

type SizingHints struct {
	RequestMemory  int `yaml:"request_memory,omitempty"`
	ReservedMemory int `yaml:"reserved_memory,omitempty"`
}

type SourceConfig struct {
	Operations map[string]map[string]string `yaml:"operations,omitempty"`
	Root       string                       `yaml:"root,omitempty"`
}

type WebConfig struct {
	Locations map[string]string `yaml:"locations,omitempty"`
}

type WorkerConfig struct {
	Flavor string `yaml:"flavor,omitempty"`
}
