// Copyright 2021 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type PlanPreviewResult struct {
	Applications        []ApplicationResult
	FailureApplications []FailureApplication
	FailurePipeds       []FailurePiped
}

func (r *PlanPreviewResult) HasError() bool {
	return len(r.FailureApplications)+len(r.FailurePipeds) > 0
}

func (r *PlanPreviewResult) NoChange() bool {
	return len(r.Applications)+len(r.FailureApplications)+len(r.FailurePipeds) == 0
}

type ApplicationResult struct {
	ApplicationInfo
	SyncStrategy string // QUICK_SYNC, PIPELINE
	PlanSummary  string
	PlanDetails  string
	NoChange     bool
}

type FailurePiped struct {
	PipedInfo
	Reason string
}

type FailureApplication struct {
	ApplicationInfo
	Reason      string
	PlanDetails string
}

type PipedInfo struct {
	PipedID  string
	PipedURL string
}

type ApplicationInfo struct {
	ApplicationID        string
	ApplicationName      string
	ApplicationURL       string
	EnvID                string
	EnvName              string
	EnvURL               string
	ApplicationKind      string // KUBERNETES, TERRAFORM, CLOUDRUN, LAMBDA, ECS
	ApplicationDirectory string
}

func retrievePlanPreview(
	ctx context.Context,
	remoteURL,
	baseBranch,
	headBranch,
	headCommit,
	address,
	apiKey string,
	timeout time.Duration,
) (*PlanPreviewResult, error) {

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create a temporary directory (%w)", err)
	}
	outPath := filepath.Join(dir, "result.json")

	args := []string{
		"plan-preview",
		"--repo-remote-url", remoteURL,
		"--base-branch", baseBranch,
		"--head-branch", headBranch,
		"--head-commit", headCommit,
		"--address", address,
		"--api-key", apiKey,
		"--timeout", timeout.String(),
		"--out", outPath,
	}
	cmd := exec.CommandContext(ctx, "pipectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipectl command (%w) (%s)", err, string(out))
	}

	log.Println(string(out))

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read result file (%w)", err)
	}

	var r PlanPreviewResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("failed to parse result file (%w)", err)
	}

	return &r, nil
}

const (
	successBadgeURL = `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=success&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

`
	failureBadgeURL = `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=orange&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)

`

	noChangeTitleFormat   = "Ran plan-preview against head commit %s of this pull request. PipeCD detected `0` updated application. It means no deployment will be triggered once this pull request got merged.\n"
	hasChangeTitleFormat  = "Ran plan-preview against head commit %s of this pull request. PipeCD detected `%d` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.\n"
	detailsFormat         = "<details>\n<summary>Details (Click me)</summary>\n<p>\n\n``` %s\n%s\n```\n</p>\n</details>\n"
	detailsOmittedMessage = "The details are too long to display. Please check the actions log to see full details."

	ghMessageLenLimit = 65536

	// Limit of details
	detailsLenLimit = ghMessageLenLimit - 5000  // 5000 characters could be used for other parts in the comment message.
)

func makeCommentBody(event *githubEvent, r *PlanPreviewResult) string {
	var b strings.Builder

	if !r.HasError() {
		b.WriteString(successBadgeURL)
	} else {
		b.WriteString(failureBadgeURL)
	}

	if event.IsComment {
		b.WriteString(fmt.Sprintf("@%s ", event.SenderLogin))
	}

	if r.NoChange() {
		fmt.Fprintf(&b, noChangeTitleFormat, event.HeadCommit)
		return b.String()
	}

	b.WriteString(fmt.Sprintf(hasChangeTitleFormat, event.HeadCommit, len(r.Applications)))

	changedApps, pipelineApps, quickSyncApps := groupApplicationResults(r.Applications)

	sort.SliceStable(changedApps, func(i, j int) bool {
		return len(changedApps[i].PlanDetails) < len(changedApps[j].PlanDetails)
	})

	var detailLen int64

	for _, app := range changedApps {
		fmt.Fprintf(&b, "\n## app: [%s](%s), env: [%s](%s), kind: %s\n", app.ApplicationName, app.ApplicationURL, app.EnvName, app.EnvURL, strings.ToLower(app.ApplicationKind))
		fmt.Fprintf(&b, "Sync strategy: %s\n", app.SyncStrategy)
		fmt.Fprintf(&b, "Summary: %s\n\n", app.PlanSummary)

		var lang string = "diff"
		if app.ApplicationKind == "TERRAFORM" {
			lang = "hcl"
		}

		l := utf8.RuneCountInString(app.PlanDetails)
		if detailLen+int64(l) > reservedDetailMessagesLen {
			fmt.Fprintf(&b, detailsFormat, lang, detailsOmittedMessage)
			detailLen += int64(utf8.RuneCountInString(detailsOmittedMessage))
			continue
		}

		detailLen += int64(l)
		fmt.Fprintf(&b, detailsFormat, lang, app.PlanDetails)
	}

	if len(pipelineApps)+len(quickSyncApps) > 0 {
		b.WriteString("\n## No resource changes were detected but the following apps will also be triggered\n")

		if len(pipelineApps) > 0 {
			b.WriteString("\n###### `PIPELINE`\n")
			for _, app := range pipelineApps {
				fmt.Fprintf(&b, "\n- app: [%s](%s), env: [%s](%s), kind: %s\n", app.ApplicationName, app.ApplicationURL, app.EnvName, app.EnvURL, strings.ToLower(app.ApplicationKind))
			}
		}

		if len(quickSyncApps) > 0 {
			b.WriteString("\n###### `QUICK_SYNC`\n")
			for _, app := range quickSyncApps {
				fmt.Fprintf(&b, "\n- app: [%s](%s), env: [%s](%s), kind: %s\n", app.ApplicationName, app.ApplicationURL, app.EnvName, app.EnvURL, strings.ToLower(app.ApplicationKind))
			}
		}
	}

	if !r.HasError() {
		return b.String()
	}

	fmt.Fprintf(&b, "\n---\n\n## NOTE\n\n")

	if len(r.FailureApplications) > 0 {
		fmt.Fprintf(&b, "**An error occurred while building plan-preview for the following applications**\n")

		sort.SliceStable(r.FailureApplications, func(i, j int) bool {
			return len(r.FailureApplications[i].PlanDetails) < len(r.FailureApplications[j].PlanDetails)
		})

		for _, app := range r.FailureApplications {
			fmt.Fprintf(&b, "\n## app: [%s](%s), env: [%s](%s), kind: %s\n", app.ApplicationName, app.ApplicationURL, app.EnvName, app.EnvURL, strings.ToLower(app.ApplicationKind))
			fmt.Fprintf(&b, "Reason: %s\n\n", app.Reason)

			var lang = "diff"
			if app.ApplicationKind == "TERRAFORM" {
				lang = "hcl"
			}

			fmt.Fprintf(&b, detailsFormat, lang, app.PlanDetails)
		}
	}

	if len(r.FailurePipeds) > 0 {
		fmt.Fprintf(&b, "**An error occurred while building plan-preview for applications of the following Pipeds**\n")

		for _, piped := range r.FailurePipeds {
			fmt.Fprintf(&b, "\n## piped: [%s](%s)\n", piped.PipedID, piped.PipedURL)
			fmt.Fprintf(&b, "Reason: %s\n\n", piped.Reason)
		}
	}

	return b.String()
}

func groupApplicationResults(apps []ApplicationResult) (changes, pipelines, quicks []ApplicationResult) {
	for _, app := range apps {
		if !app.NoChange {
			changes = append(changes, app)
			continue
		}
		if app.SyncStrategy == "PIPELINE" {
			pipelines = append(pipelines, app)
			continue
		}
		quicks = append(quicks, app)
	}
	return
}
