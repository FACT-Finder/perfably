package model

import (
	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
)

type ReportEntry struct {
	Key    semver.Version  `json:"key"`
	Values state.DataPoint `json:"values"`
	Meta   state.MetaPoint `json:"meta"`
}
type MetaValue struct {
	Value string `json:"value"`
	URL   string `json:"url,omitempty"`
}
