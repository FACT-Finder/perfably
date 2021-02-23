package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func New(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	for id, project := range config.Projects {
		if project.IDType != ReportIDTypeSemver && project.IDType != ReportIDTypeInt {
			return nil, fmt.Errorf("invalid id_type %s for project %s", project.IDType, id)
		}
	}

	return config, nil
}
