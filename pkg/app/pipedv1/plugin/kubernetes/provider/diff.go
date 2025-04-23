// Copyright 2024 The PipeCD Authors.
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

package provider

import (
	"fmt"
	"sort"
	"strings"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pipe-cd/pipecd/pkg/plugin/diff"
)

func Diff(old, new Manifest, logger *zap.Logger, opts ...diff.Option) (*diff.Result, error) {
	if old.IsSecret() && new.IsSecret() {
		var err error
		old.body, err = normalizeNewSecret(old.body, new.body)
		if err != nil {
			return nil, err
		}
	}

	key := old.Key().String()

	normalizedOld, err := remarshal(old.body)
	if err != nil {
		logger.Info("compare manifests directly since it was unable to remarshal old Kubernetes manifest to normalize special fields", zap.Error(err))
		return diff.DiffUnstructureds(*old.body, *new.body, key, opts...)
	}

	normalizedNew, err := remarshal(new.body)
	if err != nil {
		logger.Info("compare manifests directly since it was unable to remarshal new Kubernetes manifest to normalize special fields", zap.Error(err))
		return diff.DiffUnstructureds(*old.body, *new.body, key, opts...)
	}

	return diff.DiffUnstructureds(*normalizedOld, *normalizedNew, key, opts...)
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

type DiffListResult struct {
	Adds    []Manifest
	Deletes []Manifest
	Changes []DiffListChange
}

type DiffListChange struct {
	Old  Manifest
	New  Manifest
	Diff *diff.Result
}

func (r *DiffListResult) NoChanges() bool {
	return len(r.Adds)+len(r.Deletes)+len(r.Changes) == 0
}

func (r *DiffListResult) TotalOutOfSync() int {
	return len(r.Adds) + len(r.Deletes) + len(r.Changes)
}

func DiffList(liveManifests, desiredManifests []Manifest, logger *zap.Logger, opts ...diff.Option) (*DiffListResult, error) {
	adds, deletes, newChanges, oldChanges := groupManifests(liveManifests, desiredManifests)
	result := &DiffListResult{
		Adds:    adds,
		Deletes: deletes,
		Changes: make([]DiffListChange, 0, len(newChanges)),
	}

	for i := 0; i < len(newChanges); i++ {
		diffResult, err := Diff(oldChanges[i], newChanges[i], logger, opts...)
		if err != nil {
			logger.Error("Failed to diff manifests", zap.Error(err))
			continue
		}
		if !diffResult.HasDiff() {
			continue
		}
		result.Changes = append(result.Changes, DiffListChange{
			Old:  oldChanges[i],
			New:  newChanges[i],
			Diff: diffResult,
		})
	}

	return result, nil
}

func groupManifests(olds, news []Manifest) (adds, deletes, newChanges, oldChanges []Manifest) {
	// Sort the manifests before comparing.
	sort.Slice(news, func(i, j int) bool {
		return news[i].Key().normalize().String() < news[j].Key().normalize().String()
	})
	sort.Slice(olds, func(i, j int) bool {
		return olds[i].Key().normalize().String() < olds[j].Key().normalize().String()
	})

	var n, o int
	for {
		if n >= len(news) || o >= len(olds) {
			break
		}
		if news[n].Key().normalize().String() == olds[o].Key().normalize().String() {
			newChanges = append(newChanges, news[n])
			oldChanges = append(oldChanges, olds[o])
			n++
			o++
			continue
		}
		// Has in news but not in olds so this should be a added one.
		if news[n].Key().normalize().String() < olds[o].Key().normalize().String() {
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
	return adds, deletes, newChanges, oldChanges
}

type DiffRenderOptions struct {
	MaskSecret    bool
	MaskConfigMap bool
	// Maximum number of changed manifests should be shown.
	// Zero means rendering all.
	MaxChangedManifests int
}

func (r *DiffListResult) Render(opt DiffRenderOptions) string {
	var b strings.Builder
	index := 0
	for _, delete := range r.Deletes {
		index++
		b.WriteString(fmt.Sprintf("- %d. %s\n\n", index, delete.Key().ReadableString()))
	}
	for _, add := range r.Adds {
		index++
		b.WriteString(fmt.Sprintf("+ %d. %s\n\n", index, add.Key().ReadableString()))
	}

	maxPrintDiffs := len(r.Changes)
	if opt.MaxChangedManifests != 0 && opt.MaxChangedManifests < maxPrintDiffs {
		maxPrintDiffs = opt.MaxChangedManifests
	}

	for _, change := range r.Changes[:maxPrintDiffs] {
		key := change.Old.Key()
		opts := []diff.RenderOption{
			diff.WithLeftPadding(1),
		}

		if opt.MaskSecret && change.Old.IsSecret() {
			opts = append(opts, diff.WithMaskPath("data"))
		} else if opt.MaskConfigMap && change.Old.IsConfigMap() {
			opts = append(opts, diff.WithMaskPath("data"))
		}
		renderer := diff.NewRenderer(opts...)

		index++
		b.WriteString(fmt.Sprintf("# %d. %s\n\n", index, key.ReadableString()))

		b.WriteString(renderer.Render(change.Diff.Nodes()))
		b.WriteString("\n")
	}

	if maxPrintDiffs < len(r.Changes) {
		b.WriteString(fmt.Sprintf("... (omitted %d other changed manifests)\n", len(r.Changes)-maxPrintDiffs))
	}

	return b.String()
}
