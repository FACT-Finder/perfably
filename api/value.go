package api

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/coreos/go-semver/semver"
	"github.com/labstack/echo/v4"
)

func (a *api) GetValues(ctx echo.Context, projectName string, params GetValuesParams) error {
	project, ok := a.s.Projects[projectName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("project not found: %s", projectName))
	}

	project.Lock.RLock()
	defer project.Lock.RUnlock()

	startIndex, endIndex, err := filter(project.Versions, OrDefault(params.Start, ""), OrDefault(params.End, ""))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not filter: %s", err))
	}

	result := []ReportEntry{}
	for _, v := range project.Versions[startIndex : endIndex+1] {
		vState, ok := project.Data[*v]
		if !ok {
			continue
		}

		metaValues := map[string]MetaValue{}
		for key, value := range vState.Meta {
			var url *string = nil
			if value.URL != "" {
				clone := value.URL
				url = &clone
			}

			metaValues[key] = MetaValue{Value: value.Value, Url: url}
		}

		entry := ReportEntry{
			Key:    v.String(),
			Values: Metrics{AdditionalProperties: map[string]float64{}},
			Meta:   MetaValues{AdditionalProperties: metaValues},
		}
		for _, key := range params.Key {
			if value, ok := vState.Values[key]; ok {
				entry.Values.AdditionalProperties[key] = value
			}
		}
		result = append(result, entry)
	}

	return ctx.JSON(http.StatusOK, result)
}

func OrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

func filter(versions semver.Versions, startStr, endStr string) (int, int, error) {
	start, err := semver.NewVersion(startStr)
	if err != nil && startStr != "" {
		return -1, -1, fmt.Errorf("invalid start: %s", err)
	}
	end, err := semver.NewVersion(endStr)
	if err != nil && endStr != "" {
		return -1, -1, fmt.Errorf("invalid end: %s", err)
	}

	startIndex := 0
	if start != nil {
		startIndex = sort.Search(len(versions), func(i int) bool {
			return !versions[i].LessThan(*start)
		})
	}
	endIndex := len(versions) - 1
	if end != nil {
		endIndex = sort.Search(len(versions), func(i int) bool {
			return !versions[i].LessThan(*end)
		})
		if endIndex >= len(versions) {
			endIndex = len(versions) - 1
		}
	}
	return startIndex, endIndex, nil
}
