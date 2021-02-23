package rediskey

import "fmt"

func Metric(projectID, reportID, metricID string) string {
	return fmt.Sprintf("perfably.v1.project.%s.report.%s.metric.%s", projectID, reportID, metricID)
}

func ReportIDs(projectID string) string {
	return fmt.Sprintf("perfably.v1.project.%s.report_ids", projectID)
}

func Metrics(projectID string) string {
	return fmt.Sprintf("perfably.v1.project.%s.metric_ids", projectID)
}
