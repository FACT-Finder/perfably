package config

type Config struct {
	Projects map[string]*Project `yaml:"projects" json:"projects"`
}

type Project struct {
	Name       string       `yaml:"name" json:"name"`
	Layers     []string     `yaml:"pattern" json:"layers"`
	Dashboards []*Dashboard `yaml:"dashboards" json:"dashboards"`
}

type Dashboard struct {
	Name   string   `yaml:"name" json:"name"`
	Charts []*Chart `yaml:"charts" json:"charts"`
}

type Chart struct {
	Name    string   `yaml:"name" json:"name"`
	Unit    string   `yaml:"unit,omitempty" json:"unit,omitempty"`
	Metrics []string `yaml:"metrics" json:"metrics"`
}

const (
	UnitNanosecond  = "ns"
	UnitMicrosecond = "us"
	UnitMillisecond = "ms"
	UnitSecond      = "s"
	UnitMinute      = "m"
	UnitHour        = "h"
)
