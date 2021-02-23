package model

import (
	"strconv"

	"github.com/FACT-Finder/perfably/config"
	"github.com/coreos/go-semver/semver"
)

type DataPoint map[string]float64

type ReportEntry struct {
	Key    string    `json:"key"`
	Values DataPoint `json:"values"`
}

func ValidId(reportIDType, reportId string) (err error) {
	switch reportIDType {
	case config.ReportIDTypeInt:
		_, err = strconv.ParseInt(reportId, 10, 64)
		return
	case config.ReportIDTypeSemver:
		_, err = semver.NewVersion(reportId)
		return
	default:
		panic("invalid idType " + reportIDType)
	}
}
