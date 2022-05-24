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

// Package eventwatcher provides facilities to update config files when new
// event found. It can be done by periodically comparing the latest value user
// registered and the value in the files placed at Git repositories.
package eventwatcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp/syntax"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/backoff"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/regexpool"
	"github.com/pipe-cd/pipecd/pkg/yamlprocessor"
)

const (
	// The latest value and Event name are supposed.
	defaultCommitMessageFormat = "Replace values with %q set by Event %q"
	defaultCheckInterval       = time.Minute

	retryPushNum      = 3
	retryPushInterval = 5 * time.Second

	// NOTE: We don't support batch update, that's why we have a constant number of updates to make.
	numToMakeOutdated = 10
)

type Watcher interface {
	Run(context.Context) error
}

type eventLister interface {
	// Gives back not-handled event list which is sorted by newest.
	ListNotHandled(name string, labels map[string]string, minCreatedAt int64, limit int) []*model.Event
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type apiClient interface {
	GetLatestEvent(ctx context.Context, req *pipedservice.GetLatestEventRequest, opts ...grpc.CallOption) (*pipedservice.GetLatestEventResponse, error)
	ReportEventStatuses(ctx context.Context, req *pipedservice.ReportEventStatusesRequest, opts ...grpc.CallOption) (*pipedservice.ReportEventStatusesResponse, error)
}

type watcher struct {
	config      *config.PipedSpec
	eventLister eventLister
	gitClient   gitClient
	apiClient   apiClient
	logger      *zap.Logger
	wg          sync.WaitGroup

	// All cloned repository will be placed under this.
	workingDir string
	// Maximum timestamp of the last Event read.
	// A map from repo-id to the UNIX timestamp that has been read.
	milestoneMap sync.Map
}

func NewWatcher(cfg *config.PipedSpec, eventLister eventLister, gitClient gitClient, apiClient apiClient, logger *zap.Logger) Watcher {
	return &watcher{
		config:      cfg,
		eventLister: eventLister,
		gitClient:   gitClient,
		apiClient:   apiClient,
		logger:      logger.Named("event-watcher"),
	}
}

// Run spawns goroutines for each git repository. They periodically fetch the latest Event
// from the control-plane to compare the value with one in the git repository.
func (w *watcher) Run(ctx context.Context) error {
	w.logger.Info("start running event watcher")

	workingDir, err := os.MkdirTemp("", "event-watcher")
	if err != nil {
		return fmt.Errorf("failed to create the working directory: %w", err)
	}
	defer os.RemoveAll(workingDir)
	w.workingDir = workingDir

	repoCfgs := w.config.GetRepositoryMap()
	for _, r := range w.config.EventWatcher.GitRepos {
		cfg, ok := repoCfgs[r.RepoID]
		if !ok {
			return fmt.Errorf("repo id %q doesn't exist in the Piped repositories list", r.RepoID)
		}
		repo, err := w.cloneRepo(ctx, cfg)
		if err != nil {
			return err
		}
		defer repo.Clean()

		w.wg.Add(1)
		go w.run(ctx, repo, cfg)
	}

	w.wg.Wait()
	return nil
}

// run works against a single git repo. It periodically compares the value in the given
// git repository and one in the control-plane. And then pushes those with differences.
func (w *watcher) run(ctx context.Context, repo git.Repo, repoCfg config.PipedRepository) {
	defer w.wg.Done()

	var (
		commitMsg                  string
		includedCfgs, excludedCfgs []string
	)
	// Use user-defined settings if there is.
	for _, r := range w.config.EventWatcher.GitRepos {
		if r.RepoID != repoCfg.RepoID {
			continue
		}
		commitMsg = r.CommitMessage
		includedCfgs = r.Includes
		excludedCfgs = r.Excludes
		break
	}
	checkInterval := time.Duration(w.config.EventWatcher.CheckInterval)
	if checkInterval == 0 {
		checkInterval = defaultCheckInterval
	}

	w.logger.Info("start watching events", zap.String("repo", repoCfg.RepoID))
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := repo.Pull(ctx, repo.GetClonedBranch())
			if err != nil {
				w.logger.Error("failed to perform git pull",
					zap.String("repo-id", repoCfg.RepoID),
					zap.String("branch", repo.GetClonedBranch()),
					zap.Error(err),
				)
				if err := repo.Clean(); err != nil {
					w.logger.Error("failed to remove repo directory",
						zap.String("path", repo.GetPath()),
						zap.Error(err),
					)
				}
				w.logger.Info("Try to re-clone because it's more likely to be unable to pull the next time too",
					zap.String("repo-id", repoCfg.RepoID),
				)
				repo, err = w.cloneRepo(ctx, repoCfg)
				if err != nil {
					w.logger.Error("failed to re-clone repository",
						zap.String("repo-id", repoCfg.RepoID),
						zap.Error(err),
					)
				}
				continue
			}
			cfg, err := config.LoadEventWatcher(repo.GetPath(), includedCfgs, excludedCfgs)
			if errors.Is(err, config.ErrNotFound) {
				w.logger.Info("configuration file for Event Watcher not found",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if err != nil {
				w.logger.Error("failed to load configuration file for Event Watcher",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
				continue
			}
			if err := w.updateValues(ctx, repo, repoCfg.RepoID, cfg.Events, commitMsg); err != nil {
				w.logger.Error("failed to update the values",
					zap.String("repo-id", repoCfg.RepoID),
					zap.Error(err),
				)
			}
		}
	}
}

// cloneRepo clones the git repository under the working directory.
func (w *watcher) cloneRepo(ctx context.Context, repoCfg config.PipedRepository) (git.Repo, error) {
	dst, err := os.MkdirTemp(w.workingDir, repoCfg.RepoID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new temporary directory: %w", err)
	}
	repo, err := w.gitClient.Clone(ctx, repoCfg.RepoID, repoCfg.Remote, repoCfg.Branch, dst)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository %s: %w", repoCfg.RepoID, err)
	}
	return repo, nil
}

// updateValues inspects all Event-definition and pushes the changes to git repo if there is.
func (w *watcher) updateValues(ctx context.Context, repo git.Repo, repoID string, eventCfgs []config.EventWatcherEvent, commitMsg string) error {
	// Copy the repo to another directory to modify local file to avoid reverting previous changes.
	tmpDir, err := os.MkdirTemp(w.workingDir, "repo")
	if err != nil {
		return fmt.Errorf("failed to create a new temporary directory: %w", err)
	}
	tmpRepo, err := repo.Copy(filepath.Join(tmpDir, "tmp-repo"))
	if err != nil {
		return fmt.Errorf("failed to copy the repository to the temporary directory: %w", err)
	}
	defer tmpRepo.Clean()

	var milestone int64
	firstRead := true
	if v, ok := w.milestoneMap.Load(repoID); ok {
		milestone = v.(int64)
		firstRead = false
	}
	var (
		handledEvents    = make([]*pipedservice.ReportEventStatusesRequest_Event, 0, len(eventCfgs))
		outDatedEvents   = make([]*pipedservice.ReportEventStatusesRequest_Event, 0)
		maxTimestamp     int64
		outDatedDuration = time.Hour
	)
	for _, e := range eventCfgs {
		notHandledEvents := w.eventLister.ListNotHandled(e.Name, e.Labels, milestone+1, numToMakeOutdated)
		if len(notHandledEvents) == 0 {
			continue
		}
		if len(notHandledEvents) > 1 {
			// Events other than the latest will be OUTDATED.
			for _, e := range notHandledEvents[1:] {
				outDatedEvents = append(outDatedEvents, &pipedservice.ReportEventStatusesRequest_Event{
					Id:                e.Id,
					Status:            model.EventStatus_EVENT_OUTDATED,
					StatusDescription: fmt.Sprintf("The new event %q has been created", notHandledEvents[0].Id),
				})
			}
		}

		latestEvent := notHandledEvents[0]
		if firstRead {
			resp, err := w.apiClient.GetLatestEvent(ctx, &pipedservice.GetLatestEventRequest{
				Name:   e.Name,
				Labels: e.Labels,
			})
			if err != nil {
				return fmt.Errorf("failed to get the latest event: %w", err)
			}
			// The case where the latest event has already been handled.
			if resp.Event.CreatedAt > latestEvent.CreatedAt {
				outDatedEvents = append(outDatedEvents, &pipedservice.ReportEventStatusesRequest_Event{
					Id:                notHandledEvents[0].Id,
					Status:            model.EventStatus_EVENT_OUTDATED,
					StatusDescription: fmt.Sprintf("The new event %q has been created", resp.Event.Id),
				})
				continue
			}
		}
		if time.Since(time.Unix(latestEvent.CreatedAt, 0)) > outDatedDuration {
			outDatedEvents = append(outDatedEvents, &pipedservice.ReportEventStatusesRequest_Event{
				Id:                latestEvent.Id,
				Status:            model.EventStatus_EVENT_OUTDATED,
				StatusDescription: fmt.Sprintf("Too much time has passed since the event %q was created", latestEvent.Id),
			})
			continue
		}
		if err := w.commitFiles(ctx, latestEvent.Data, e, tmpRepo, commitMsg); err != nil {
			w.logger.Error("failed to commit outdated files", zap.Error(err))
			handledEvents = append(handledEvents, &pipedservice.ReportEventStatusesRequest_Event{
				Id:                latestEvent.Id,
				Status:            model.EventStatus_EVENT_FAILURE,
				StatusDescription: fmt.Sprintf("Failed to change files: %v", err),
			})
			continue
		}
		handledEvents = append(handledEvents, &pipedservice.ReportEventStatusesRequest_Event{
			Id:                latestEvent.Id,
			Status:            model.EventStatus_EVENT_SUCCESS,
			StatusDescription: fmt.Sprintf("Successfully updated %d files in the %q repository", len(e.Replacements), repoID),
		})
		if latestEvent.CreatedAt > maxTimestamp {
			maxTimestamp = latestEvent.CreatedAt
		}
	}
	if len(outDatedEvents) > 0 {
		if _, err := w.apiClient.ReportEventStatuses(ctx, &pipedservice.ReportEventStatusesRequest{Events: outDatedEvents}); err != nil {
			return fmt.Errorf("failed to report event statuses: %w", err)
		}
		w.logger.Info(fmt.Sprintf("successfully made %d events OUTDATED", len(outDatedEvents)))
	}
	if len(handledEvents) == 0 {
		return nil
	}

	retry := backoff.NewRetry(retryPushNum, backoff.NewConstant(retryPushInterval))
	_, err = retry.Do(ctx, func() (interface{}, error) {
		err := tmpRepo.Push(ctx, tmpRepo.GetClonedBranch())
		return nil, err
	})
	if err == nil {
		if _, err := w.apiClient.ReportEventStatuses(ctx, &pipedservice.ReportEventStatusesRequest{Events: handledEvents}); err != nil {
			return fmt.Errorf("failed to report event statuses: %w", err)
		}
		w.milestoneMap.Store(repoID, maxTimestamp)
		return nil
	}

	// If push fails because the local branch was not fresh, exit to retry again in the next interval.
	if err == git.ErrBranchNotFresh {
		w.logger.Warn("failed to push commits", zap.Error(err))
		return nil
	}

	// If push fails because of the other reason, re-set all statuses to FAILURE.
	for i := range handledEvents {
		if handledEvents[i].Status == model.EventStatus_EVENT_FAILURE {
			continue
		}
		handledEvents[i].Status = model.EventStatus_EVENT_FAILURE
		handledEvents[i].StatusDescription = fmt.Sprintf("Failed to push changed files: %v", err)
	}
	if _, err := w.apiClient.ReportEventStatuses(ctx, &pipedservice.ReportEventStatusesRequest{Events: handledEvents}); err != nil {
		return fmt.Errorf("failed to report event statuses: %w", err)
	}
	w.milestoneMap.Store(repoID, maxTimestamp)
	return fmt.Errorf("failed to push commits: %w", err)
}

// commitFiles commits changes if the data in Git is different from the latest event.
func (w *watcher) commitFiles(ctx context.Context, latestData string, eventCfg config.EventWatcherEvent, repo git.Repo, commitMsg string) error {
	// Determine files to be changed by comparing with the latest event.
	changes := make(map[string][]byte, len(eventCfg.Replacements))
	for _, r := range eventCfg.Replacements {
		var (
			path       = filepath.Join(repo.GetPath(), r.File)
			newContent []byte
			upToDate   bool
			err        error
		)
		switch {
		case r.YAMLField != "":
			newContent, upToDate, err = modifyYAML(path, r.YAMLField, latestData)
		case r.JSONField != "":
			// TODO: Empower Event watcher to parse JSON format
		case r.HCLField != "":
			// TODO: Empower Event watcher to parse HCL format
		case r.Regex != "":
			newContent, upToDate, err = modifyText(path, r.Regex, latestData)
		}
		if err != nil {
			return err
		}
		if upToDate {
			continue
		}

		if err := os.WriteFile(path, newContent, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		changes[r.File] = newContent
	}
	if len(changes) == 0 {
		return nil
	}

	if commitMsg == "" {
		commitMsg = fmt.Sprintf(defaultCommitMessageFormat, latestData, eventCfg.Name)
	}
	if err := repo.CommitChanges(ctx, repo.GetClonedBranch(), commitMsg, false, changes); err != nil {
		return fmt.Errorf("failed to perform git commit: %w", err)
	}
	w.logger.Info(fmt.Sprintf("event watcher will update values of Event %q", eventCfg.Name))
	return nil
}

// modifyYAML returns a new YAML content as a first returned value if the value of given
// field was outdated. True as a second returned value means it's already up-to-date.
func modifyYAML(path, field, newValue string) ([]byte, bool, error) {
	yml, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read file: %w", err)
	}

	processor, err := yamlprocessor.NewProcessor(yml)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse yaml file: %w", err)
	}

	v, err := processor.GetValue(field)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get value at %s in %s: %w", field, path, err)
	}
	value, err := convertStr(v)
	if err != nil {
		return nil, false, fmt.Errorf("a value of unknown type is defined at %s in %s: %w", field, path, err)
	}
	if newValue == value {
		// Already up-to-date.
		return nil, true, nil
	}

	// Modify the local file and put it into the change list.
	if err := processor.ReplaceString(field, newValue); err != nil {
		return nil, false, fmt.Errorf("failed to replace value at %s with %s: %w", field, newValue, err)
	}

	return processor.Bytes(), false, nil
}

// convertStr converts a given value into a string.
func convertStr(value interface{}) (out string, err error) {
	switch v := value.(type) {
	case string:
		out = v
	case int:
		out = strconv.Itoa(v)
	case int64:
		out = strconv.FormatInt(v, 10)
	case uint64:
		out = strconv.FormatUint(v, 10)
	case float64:
		out = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		out = strconv.FormatBool(v)
	default:
		err = fmt.Errorf("failed to convert %T into string", v)
	}
	return
}

// modifyText returns a new text replacing all matches of the given regex with the newValue.
// The only first capturing group enclosed by `()` will be replaced.
// True as a second returned value means it's already up-to-date.
func modifyText(path, regexText, newValue string) ([]byte, bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read file: %w", err)
	}

	pool := regexpool.DefaultPool()
	regex, err := pool.Get(regexText)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compile regex text (%s): %w", regexText, err)
	}

	// Extract the first capturing group.
	firstGroup := ""
	re, err := syntax.Parse(regexText, syntax.Perl)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse the first capturing group regex: %w", err)
	}
	for _, s := range re.Sub {
		if s.Op == syntax.OpCapture {
			firstGroup = s.String()
			break
		}
	}
	if firstGroup == "" {
		return nil, false, fmt.Errorf("capturing group not found in the given regex")
	}
	subRegex, err := pool.Get(firstGroup)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compile the first capturing group: %w", err)
	}

	var touched, outDated bool
	newText := regex.ReplaceAllFunc(content, func(match []byte) []byte {
		touched = true
		outDated = string(subRegex.Find(match)) != newValue
		// Return text replacing the only first capturing group with the newValue.
		return subRegex.ReplaceAll(match, []byte(newValue))
	})
	if !touched {
		return nil, false, fmt.Errorf("the content of %s doesn't match %s", path, regexText)
	}
	if !outDated {
		return nil, true, nil
	}

	return newText, false, nil
}
