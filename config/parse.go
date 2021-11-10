package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var units = map[string]bool{
	UnitNanosecond:  true,
	UnitMicrosecond: true,
	UnitMillisecond: true,
	UnitSecond:      true,
	UnitMinute:      true,
	UnitHour:        true,
}

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
		for _, dashboard := range project.Dashboards {
			for _, chart := range dashboard.Charts {
				_, valid := units[chart.Unit]
				if chart.Unit != "" && !valid {
					return nil, fmt.Errorf("invalid unit '%s' for chart %s.%s", chart.Unit, id, chart.Name)
				}
			}
		}
	}

	return config, nil
}
