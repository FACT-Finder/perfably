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

func ValidID(reportIDType, reportID string) (err error) {
	switch reportIDType {
	case config.ReportIDTypeInt:
		_, err = strconv.ParseInt(reportID, 10, 64)
		return
	case config.ReportIDTypeSemver:
		_, err = semver.NewVersion(reportID)
		return
	default:
		panic("invalid idType " + reportIDType)
	}
}
