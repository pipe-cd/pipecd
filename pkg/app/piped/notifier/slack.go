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

package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	slackUsername     = "PipeCD"
	slackInfoColor    = "#222429"
	slackSuccessColor = "#629650"
	slackErrorColor   = "#9C3C31"
	slackWarnColor    = "#C1A337"
)

type slack struct {
	name       string
	config     config.NotificationReceiverSlack
	webURL     string
	httpClient *http.Client
	eventCh    chan model.NotificationEvent
	logger     *zap.Logger
}

func newSlackSender(name string, cfg config.NotificationReceiverSlack, webURL string, logger *zap.Logger) *slack {
	return &slack{
		name:   name,
		config: cfg,
		webURL: strings.TrimRight(webURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		eventCh: make(chan model.NotificationEvent, 100),
		logger:  logger.Named("slack"),
	}
}

func (s *slack) Run(ctx context.Context) error {
	for {
		select {
		case event, ok := <-s.eventCh:
			if ok {
				s.sendEvent(ctx, event)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *slack) Notify(event model.NotificationEvent) {
	s.eventCh <- event
}

func (s *slack) Close(ctx context.Context) {
	close(s.eventCh)

	// Send all remaining events.
	for {
		select {
		case event, ok := <-s.eventCh:
			if !ok {
				return
			}
			s.sendEvent(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}

func (s *slack) sendEvent(ctx context.Context, event model.NotificationEvent) {
	msg, ok := s.buildSlackMessage(event, s.webURL)
	if !ok {
		s.logger.Info(fmt.Sprintf("ignore event %s", event.Type.String()))
		return
	}
	if err := s.sendMessage(ctx, msg); err != nil {
		s.logger.Error(fmt.Sprintf("unable to send notification to slack: %v", err))
	}
}

func (s *slack) sendMessage(ctx context.Context, msg slackMessage) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.HookURL, buf)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
		return fmt.Errorf("%s from Slack: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return nil
}

func (s *slack) buildSlackMessage(event model.NotificationEvent, webURL string) (slackMessage, bool) {
	var (
		title, link, text string
		color             = slackInfoColor
		timestamp         = time.Now().Unix()
		fields            []slackField
	)

	generateDeploymentEventData := func(d *model.Deployment, accounts string) {
		link = fmt.Sprintf("%s/deployments/%s?project=%s", webURL, d.Id, d.ProjectId)
		fields = []slackField{
			{"Project", truncateText(d.ProjectId, 8), true},
			{"Application", makeSlackLink(d.ApplicationName, fmt.Sprintf("%s/applications/%s?project=%s", webURL, d.ApplicationId, d.ProjectId)), true},
			{"Kind", strings.ToLower(d.Kind.String()), true},
			{"Deployment", makeSlackLink(truncateText(d.Id, 8), link), true},
			{"Triggered By", d.TriggeredBy(), true},
			{"Mention To", accounts, true},
			{"Started At", makeSlackDate(d.CreatedAt), true},
		}
	}
	generateDeploymentEventDataForTriggerFailed := func(app *model.Application, hash, msg string) {
		link = fmt.Sprintf("%s/applications/%s?project=%s", webURL, app.Id, app.ProjectId)
		commitURL, err := git.MakeCommitURL(app.GitPath.Repo.Remote, hash)
		if err != nil {
			s.logger.Error(fmt.Sprintf("failed to get the URL for the specified commit: %v", err))
		}
		fields = []slackField{
			{"Project", truncateText(app.ProjectId, 8), true},
			{"Application", makeSlackLink(app.Name, link), true},
			{"Kind", strings.ToLower(app.Kind.String()), true},
		}
		if commitURL != "" {
			fields = append(fields, slackField{"Commit", makeSlackLink(truncateText(msg, 8), commitURL), true})
		}

	}
	generatePipedEventData := func(id, name, version, project string) {
		link = fmt.Sprintf("%s/settings/piped?project=%s", webURL, project)
		fields = []slackField{
			{"Name", name, true},
			{"Version", version, true},
			{"Project", truncateText(project, 8), true},
			{"Id", id, true},
		}
	}

	switch event.Type {
	case model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGERED:
		md := event.Metadata.(*model.NotificationEventDeploymentTriggered)
		title = fmt.Sprintf("Triggered a new deployment for %q", md.Deployment.ApplicationName)
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_PLANNED:
		md := event.Metadata.(*model.NotificationEventDeploymentPlanned)
		title = fmt.Sprintf("Deployment for %q was planned", md.Deployment.ApplicationName)
		text = md.Summary
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_WAIT_APPROVAL:
		md := event.Metadata.(*model.NotificationEventDeploymentWaitApproval)
		title = fmt.Sprintf("Deployment for %q is waiting for an approval", md.Deployment.ApplicationName)
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_APPROVED:
		md := event.Metadata.(*model.NotificationEventDeploymentApproved)
		title = fmt.Sprintf("Deployment for %q was approved", md.Deployment.ApplicationName)
		text = fmt.Sprintf("Approved by %s", md.Approver)
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_SUCCEEDED:
		md := event.Metadata.(*model.NotificationEventDeploymentSucceeded)
		title = fmt.Sprintf("Deployment for %q was completed successfully", md.Deployment.ApplicationName)
		color = slackSuccessColor
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_FAILED:
		md := event.Metadata.(*model.NotificationEventDeploymentFailed)
		title = fmt.Sprintf("Deployment for %q was failed", md.Deployment.ApplicationName)
		text = md.Reason
		color = slackErrorColor
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_CANCELLED:
		md := event.Metadata.(*model.NotificationEventDeploymentCancelled)
		title = fmt.Sprintf("Deployment for %q was cancelled", md.Deployment.ApplicationName)
		text = fmt.Sprintf("Cancelled by %s", md.Commander)
		color = slackWarnColor
		generateDeploymentEventData(md.Deployment, getAccountsAsString(md.MentionedAccounts))

	case model.NotificationEventType_EVENT_DEPLOYMENT_TRIGGER_FAILED:
		md := event.Metadata.(*model.NotificationEventDeploymentTriggerFailed)
		title = fmt.Sprintf("Failed to trigger a new deployment for %s", md.Application.Name)
		text = md.Reason
		generateDeploymentEventDataForTriggerFailed(md.Application, md.CommitHash, md.CommitMessage)

	case model.NotificationEventType_EVENT_PIPED_STARTED:
		md := event.Metadata.(*model.NotificationEventPipedStarted)
		title = "A piped has been started"
		generatePipedEventData(md.Id, md.Name, md.Version, md.ProjectId)

	case model.NotificationEventType_EVENT_PIPED_STOPPED:
		md := event.Metadata.(*model.NotificationEventPipedStopped)
		title = "A piped has been stopped"
		generatePipedEventData(md.Id, md.Name, md.Version, md.ProjectId)

	// TODO: Support application type of notification event.
	default:
		return slackMessage{}, false
	}

	return makeSlackMessage(title, link, text, color, timestamp, fields...), true
}

type slackMessage struct {
	Username    string            `json:"username"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Title     string       `json:"title"`
	TitleLink string       `json:"title_link"`
	Text      string       `json:"text"`
	Fields    []slackField `json:"fields"`
	Color     string       `json:"color,omitempty"`
	Markdown  []string     `json:"mrkdwn_in,omitempty"`
	Timestamp int64        `json:"ts,omitempty"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func makeSlackLink(title, url string) string {
	return fmt.Sprintf("<%s|%s>", url, title)
}

func makeSlackDate(unix int64) string {
	return fmt.Sprintf("<!date^%d^{date_num} {time_secs}|date>", unix)
}

func truncateText(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}

func makeSlackMessage(title, titleLink, text, color string, timestamp int64, fields ...slackField) slackMessage {
	return slackMessage{
		Username: slackUsername,
		Attachments: []slackAttachment{{
			Title:     title,
			TitleLink: titleLink,
			Text:      text,
			Fields:    fields,
			Color:     color,
			Markdown:  []string{"text"},
			Timestamp: timestamp,
		}},
	}
}

func getAccountsAsString(accounts []string) string {
	if len(accounts) == 0 {
		return ""
	}
	formattedAccounts := make([]string, 0, len(accounts))
	for _, a := range accounts {
		formattedAccounts = append(formattedAccounts, fmt.Sprintf("<@%s>", a))
	}
	return strings.Join(formattedAccounts, " ")
}
