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

package kubernetes

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pipe-cd/pipecd/pkg/diff"
)

const (
	diffCommand = "diff"
)

type DiffListResult struct {
	Adds    []Manifest
	Deletes []Manifest
	Changes []DiffListChange
}

func (r *DiffListResult) NoChange() bool {
	return len(r.Adds)+len(r.Deletes)+len(r.Changes) == 0
}

type DiffListChange struct {
	Old  Manifest
	New  Manifest
	Diff *diff.Result
}

func Diff(old, new Manifest, logger *zap.Logger, opts ...diff.Option) (*diff.Result, error) {
	if old.Key.IsSecret() && new.Key.IsSecret() {
		var err error
		old.u, err = normalizeNewSecret(old.u, new.u)
		if err != nil {
			return nil, err
		}
	}

	normalized, err := remarshal(new.u)
	if err != nil {
		logger.Info("compare manifests directly since it was unable to remarshal Kubernetes manifest to normalize special fields", zap.Error(err))
		return diff.DiffUnstructureds(*old.u, *new.u, opts...)
	}

	return diff.DiffUnstructureds(*old.u, *normalized, opts...)
}

func DiffList(olds, news []Manifest, logger *zap.Logger, ignoredPathsOf map[string][]string, opts ...diff.Option) (*DiffListResult, error) {
	adds, deletes, newChanges, oldChanges := groupManifests(olds, news)
	cr := &DiffListResult{
		Adds:    adds,
		Deletes: deletes,
		Changes: make([]DiffListChange, 0, len(newChanges)),
	}

	for i := 0; i < len(newChanges); i++ {
		key := oldChanges[i].Key
		opts = append(opts, diff.WithIgnoredPaths(ignoredPathsOf[key.String()]))

		result, err := Diff(oldChanges[i], newChanges[i], logger, opts...)
		if err != nil {
			return nil, err
		}
		if !result.HasDiff() {
			continue
		}
		cr.Changes = append(cr.Changes, DiffListChange{
			Old:  oldChanges[i],
			New:  newChanges[i],
			Diff: result,
		})
	}

	return cr, nil
}

func normalizeNewSecret(old, new *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var o, n v1.Secret
	runtime.DefaultUnstructuredConverter.FromUnstructured(old.Object, &o)
	runtime.DefaultUnstructuredConverter.FromUnstructured(new.Object, &n)

	// Move as much as possible fields from `o.Data` to `o.StringData` to make `o` close to `n` to minimize the diff.
	for k, v := range o.Data {
		// Skip if the field also exists in StringData.
		if _, ok := o.StringData[k]; ok {
			continue
		}

		if _, ok := n.StringData[k]; !ok {
			continue
		}

		if o.StringData == nil {
			o.StringData = make(map[string]string)
		}

		// If the field is existing in `n.StringData`, we should move that field from `o.Data` to `o.StringData`
		o.StringData[k] = string(v)
		delete(o.Data, k)
	}

	newO, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&o)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: newO}, nil
}

type DiffRenderOptions struct {
	MaskSecret    bool
	MaskConfigMap bool
	// Maximum number of changed manifests should be shown.
	// Zero means rendering all.
	MaxChangedManifests int
	// If true, use "diff" command to render.
	UseDiffCommand bool
}

func (r *DiffListResult) Render(opt DiffRenderOptions) string {
	var b strings.Builder
	index := 0
	for _, delete := range r.Deletes {
		index++
		b.WriteString(fmt.Sprintf("- %d. %s\n\n", index, delete.Key.ReadableString()))
	}
	for _, add := range r.Adds {
		index++
		b.WriteString(fmt.Sprintf("+ %d. %s\n\n", index, add.Key.ReadableString()))
	}

	maxPrintDiffs := len(r.Changes)
	if opt.MaxChangedManifests != 0 && opt.MaxChangedManifests < maxPrintDiffs {
		maxPrintDiffs = opt.MaxChangedManifests
	}

	var prints = 0
	for _, change := range r.Changes {
		key := change.Old.Key
		opts := []diff.RenderOption{
			diff.WithLeftPadding(1),
		}

		needMaskValue := false
		if opt.MaskSecret && key.IsSecret() {
			opts = append(opts, diff.WithMaskPath("data"))
			needMaskValue = true
		} else if opt.MaskConfigMap && key.IsConfigMap() {
			opts = append(opts, diff.WithMaskPath("data"))
			needMaskValue = true
		}
		renderer := diff.NewRenderer(opts...)

		index++
		b.WriteString(fmt.Sprintf("# %d. %s\n\n", index, key.ReadableString()))

		// Use our diff check in one of the following cases:
		// - not explicit set useDiffCommand option.
		// - requires masking secret or configmap value.
		if !opt.UseDiffCommand || needMaskValue {
			b.WriteString(renderer.Render(change.Diff.Nodes()))
		} else {
			// TODO: Find a way to mask values in case of using unix `diff` command.
			d, err := diffByCommand(diffCommand, change.Old, change.New)
			if err != nil {
				b.WriteString(fmt.Sprintf("An error occurred while rendering diff (%v)", err))
			} else {
				b.Write(d)
			}
		}
		b.WriteString("\n")

		prints++
		if prints >= maxPrintDiffs {
			break
		}
	}

	if prints < len(r.Changes) {
		b.WriteString(fmt.Sprintf("... (omitted %d other changed manifests\n", len(r.Changes)-prints))
	}

	return b.String()
}

func diffByCommand(command string, old, new Manifest) ([]byte, error) {
	oldBytes, err := old.YamlBytes()
	if err != nil {
		return nil, err
	}

	newBytes, err := new.YamlBytes()
	if err != nil {
		return nil, err
	}

	oldFile, err := os.CreateTemp("", "old")
	if err != nil {
		return nil, err
	}
	defer os.Remove(oldFile.Name())
	if _, err := oldFile.Write(oldBytes); err != nil {
		return nil, err
	}

	newFile, err := os.CreateTemp("", "new")
	if err != nil {
		return nil, err
	}
	defer os.Remove(newFile.Name())
	if _, err := newFile.Write(newBytes); err != nil {
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, "-u", "-N", oldFile.Name(), newFile.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if stdout.Len() > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to run diff, err = %w, %s", err, stderr.String())
	}

	// Remove two-line header from output.
	data := bytes.TrimSpace(stdout.Bytes())
	rows := bytes.SplitN(data, []byte("\n"), 3)
	if len(rows) == 3 {
		return rows[2], nil
	}
	return data, nil
}

func groupManifests(olds, news []Manifest) (adds, deletes, newChanges, oldChanges []Manifest) {
	// Sort the manifests before comparing.
	sort.Slice(news, func(i, j int) bool {
		return news[i].Key.IsLessWithIgnoringNamespace(news[j].Key)
	})
	sort.Slice(olds, func(i, j int) bool {
		return olds[i].Key.IsLessWithIgnoringNamespace(olds[j].Key)
	})

	var n, o int
	for {
		if n >= len(news) || o >= len(olds) {
			break
		}
		if news[n].Key.IsEqualWithIgnoringNamespace(olds[o].Key) {
			newChanges = append(newChanges, news[n])
			oldChanges = append(oldChanges, olds[o])
			n++
			o++
			continue
		}
		// Has in news but not in olds so this should be a added one.
		if news[n].Key.IsLessWithIgnoringNamespace(olds[o].Key) {
			adds = append(adds, news[n])
			n++
			continue
		}
		// Has in olds but not in news so this should be an deleted one.
		deletes = append(deletes, olds[o])
		o++
	}

	if len(news) > n {
		adds = append(adds, news[n:]...)
	}
	if len(olds) > o {
		deletes = append(deletes, olds[o:]...)
	}
	return
}
