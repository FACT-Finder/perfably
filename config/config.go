package config

type Config struct {
	Projects map[string]*Project `yaml:"projects"`
}

type Project struct {
	Name       string       `yaml:"name"`
	IDType     string       `yaml:"id_type"`
	Layers     []string     `yaml:"pattern"`
	Dashboards []*Dashboard `yaml:"dashboards"`
}

type Dashboard struct {
	Name   string   `yaml:"name"`
	Charts []*Chart `yaml:"charts"`
}

type Chart struct {
	Name    string   `yaml:"name"`
	Metrics []string `yaml:"metrics"`
}

const (
	ReportIDTypeSemver = "semver"
	ReportIDTypeInt    = "int"
)
