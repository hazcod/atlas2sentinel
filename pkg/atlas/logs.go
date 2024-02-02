package atlas

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
	"io"
	"net/http"
	"time"
)

func (p *Atlas) GetLogs(ctx context.Context, lookBackDays uint) ([][]byte, error) {
	logsBytes := make([][]byte, 0)

	projects, response, err := p.client.ProjectsApi.ListProjectsWithParams(ctx,
		&admin.ListProjectsApiParams{
			ItemsPerPage: admin.PtrInt(100),
			IncludeCount: admin.PtrBool(true),
			PageNum:      admin.PtrInt(1),
		},
	).Execute()

	if err != nil {
		return nil, fmt.Errorf("could not fetch projects: %v", err)
	}
	if response == nil {
		return nil, fmt.Errorf("nil response")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("atlas projects returned status code: %d", response.StatusCode)
	}

	if projects.GetTotalCount() == 0 {
		return nil, fmt.Errorf("no atlas projects found")
	}

	now := time.Now().UTC()
	startTime := now.AddDate(0, 0, -1*int(lookBackDays)).Unix()
	nowEpoch := now.Unix()

	for _, project := range projects.GetResults() {
		projLogger := p.Logger.WithField("module", "logs").WithField("project", project)

		hosts, _, err := p.client.MonitoringAndLogsApi.ListAtlasProcesses(ctx, project.GetId()).Execute()
		if err != nil {
			return nil, fmt.Errorf("could not fetch atlas processes: %v", err)
		}

		if hosts.GetTotalCount() == 0 {
			projLogger.Warn("no hosts found for atlas project")
			continue
		}

		for _, host := range hosts.GetResults() {
			hostLogger := p.Logger.WithField("host_hostname", host.Hostname)

			params := &admin.GetHostLogsApiParams{
				GroupId:   project.GetId(),
				HostName:  *host.Hostname,
				LogName:   "mongos",
				StartDate: &startTime,
				EndDate:   &nowEpoch,
			}

			logs, _, err := p.client.MonitoringAndLogsApi.GetHostLogsWithParams(ctx, params).Execute()
			if err != nil {
				_ = logs.Close()
				return nil, fmt.Errorf("could not fetch host '%s' logs: %v", host.Hostname, err)
			}

			logBytes, err := io.ReadAll(logs)
			if err != nil {
				_ = logs.Close()
				return nil, fmt.Errorf("could not read logs: %v", err)
			}

			logsBytes = append(logsBytes, logBytes)

			_ = logs.Close()

			hostLogger.WithField("size_bytes", len(logsBytes)).Debug("read logs")
		}
	}

	return logsBytes, nil
}
