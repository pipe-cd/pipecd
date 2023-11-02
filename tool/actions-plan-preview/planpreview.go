// Copyright 2023 The PipeCD Authors.
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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	Env                  string
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
		"--piped-handle-timeout", timeout.String(),
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
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=success&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)`
	failureBadgeURL = `<!-- pipecd-plan-preview-->
[![PLAN_PREVIEW](https://img.shields.io/static/v1?label=PipeCD&message=Plan_Preview&color=orange&style=flat)](https://pipecd.dev/docs/user-guide/plan-preview/)`
	actionBadgeURLFormat = "[![ACTIONS](https://img.shields.io/static/v1?label=PipeCD&message=Action_Log&style=flat)](%s)"

	noChangeTitleFormat     = "Ran plan-preview against head commit %s of this pull request. PipeCD detected `0` updated application. It means no deployment will be triggered once this pull request got merged.\n"
	hasChangeTitleFormat    = "Ran plan-preview against head commit %s of this pull request. PipeCD detected `%d` updated applications and here are their plan results. Once this pull request got merged their deployments will be triggered to run as these estimations.\n"
	detailsFormat           = "<details>\n<summary>Details (Click me)</summary>\n<p>\n\n``` %s\n%s\n```\n</p>\n</details>\n\n"
	detailsOmittedMessage   = "The details are too long to display. Please check the actions log to see full details."
	appInfoWithEnvFormat    = "app: [%s](%s), env: %s, kind: %s"
	appInfoWithoutEnvFormat = "app: [%s](%s), kind: %s"

	ghMessageLenLimit = 65536

	// Limit of details
	detailsLenLimit = ghMessageLenLimit - 5000 // 5000 characters could be used for other parts in the comment message.

	githubServerURLEnv  = "GITHUB_SERVER_URL"
	githubRepositoryEnv = "GITHUB_REPOSITORY"
	githubRunIDEnv      = "GITHUB_RUN_ID"

	// Terraform plan format
	prefixTerraformPlan = "Terraform will perform the following actions:"
)

func makeCommentBody(event *githubEvent, r *PlanPreviewResult) string {
	var b strings.Builder

	if !r.HasError() {
		b.WriteString(successBadgeURL)
	} else {
		b.WriteString(failureBadgeURL)
	}

	if actionLogURL := makeActionLogURL(); actionLogURL != "" {
		fmt.Fprintf(&b, " ")
		fmt.Fprintf(&b, actionBadgeURLFormat, actionLogURL)
	}
	b.WriteString("\n\n")

	if event.IsComment {
		b.WriteString(fmt.Sprintf("@%s ", event.SenderLogin))
	}

	if r.NoChange() {
		fmt.Fprintf(&b, noChangeTitleFormat, event.HeadCommit)
		return b.String()
	}

	b.WriteString(fmt.Sprintf(hasChangeTitleFormat, event.HeadCommit, len(r.Applications)))

	changedApps, pipelineApps, quickSyncApps := groupApplicationResults(r.Applications)
	if len(changedApps)+len(pipelineApps)+len(quickSyncApps) > 0 {
		b.WriteString("\n## Plans\n\n")
	}

	var detailLen int
	for _, app := range changedApps {
		fmt.Fprintf(&b, "### %s\n", makeTitleText(&app.ApplicationInfo))
		fmt.Fprintf(&b, "Sync strategy: %s\n", app.SyncStrategy)
		fmt.Fprintf(&b, "Summary: %s\n\n", app.PlanSummary)

		var (
			lang    = "diff"
			details = app.PlanDetails
		)
		if app.ApplicationKind == "TERRAFORM" {
			lang = "hcl"
			if shortened, err := generateTerraformShortPlanDetails(details); err == nil {
				details = shortened
			}
		}

		l := utf8.RuneCountInString(details)
		if detailLen+l > detailsLenLimit {
			fmt.Fprintf(&b, detailsFormat, lang, detailsOmittedMessage)
			detailLen += utf8.RuneCountInString(detailsOmittedMessage)
			continue
		}

		if l > 0 {
			detailLen += l
			fmt.Fprintf(&b, detailsFormat, lang, details)
		}
	}

	if len(pipelineApps)+len(quickSyncApps) > 0 {
		b.WriteString("### No resource changes were detected but the following apps will also be triggered\n")

		if len(pipelineApps) > 0 {
			b.WriteString("\n###### `PIPELINE`\n")
			for _, app := range pipelineApps {
				fmt.Fprintf(&b, "\n- %s\n", makeTitleText(&app.ApplicationInfo))
			}
		}

		if len(quickSyncApps) > 0 {
			b.WriteString("\n###### `QUICK_SYNC`\n")
			for _, app := range quickSyncApps {
				fmt.Fprintf(&b, "\n- %s\n", makeTitleText(&app.ApplicationInfo))
			}
		}
	}

	if !r.HasError() {
		return b.String()
	}

	fmt.Fprintf(&b, "\n## NOTE\n\n")

	if len(r.FailureApplications) > 0 {
		fmt.Fprintf(&b, "**An error occurred while building plan-preview for the following applications**\n")

		for _, app := range r.FailureApplications {
			fmt.Fprintf(&b, "\n### %s\n", makeTitleText(&app.ApplicationInfo))
			fmt.Fprintf(&b, "Reason: %s\n\n", app.Reason)

			var lang = "diff"
			if app.ApplicationKind == "TERRAFORM" {
				lang = "hcl"
			}

			if len(app.PlanDetails) > 0 {
				fmt.Fprintf(&b, detailsFormat, lang, app.PlanDetails)
			}
		}
	}

	if len(r.FailurePipeds) > 0 {
		fmt.Fprintf(&b, "**An error occurred while building plan-preview for applications of the following Pipeds**\n")

		for _, piped := range r.FailurePipeds {
			fmt.Fprintf(&b, "\n### piped: [%s](%s)\n", piped.PipedID, piped.PipedURL)
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

func makeActionLogURL() string {
	serverURL := os.Getenv(githubServerURLEnv)
	if serverURL == "" {
		return ""
	}

	repoURL := os.Getenv(githubRepositoryEnv)
	if repoURL == "" {
		return ""
	}

	runID := os.Getenv(githubRunIDEnv)
	if runID == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s/actions/runs/%s", serverURL, repoURL, runID)
}

func makeTitleText(app *ApplicationInfo) string {
	if app.Env == "" {
		return fmt.Sprintf(appInfoWithoutEnvFormat, app.ApplicationName, app.ApplicationURL, strings.ToLower(app.ApplicationKind))
	}
	return fmt.Sprintf(appInfoWithEnvFormat, app.ApplicationName, app.ApplicationURL, app.Env, strings.ToLower(app.ApplicationKind))
}

func generateTerraformShortPlanDetails(details string) (string, error) {
	r := strings.NewReader(details)
	scanner := bufio.NewScanner(r)
	var (
		start, length int
		newLine       = len([]byte("\n"))
	)
	// NOTE: scanner.Scan() return false if the buffer size of one line exceed bufio.MaxScanTokenSize(65536 byte).
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, prefixTerraformPlan) {
			start = length
			break
		}
		length += len(scanner.Bytes())
		length += newLine
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return details[start:], nil
}
