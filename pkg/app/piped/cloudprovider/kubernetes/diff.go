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

package kubernetes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pipe-cd/pipe/pkg/diff"
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

func Diff(old, new Manifest, opts ...diff.Option) (*diff.Result, error) {
	return diff.DiffUnstructureds(*old.u, *new.u, opts...)
}

func DiffList(olds, news []Manifest, opts ...diff.Option) (*DiffListResult, error) {
	adds, deletes, newChanges, oldChanges := groupManifests(olds, news)
	cr := &DiffListResult{
		Adds:    adds,
		Deletes: deletes,
		Changes: make([]DiffListChange, 0, len(newChanges)),
	}

	for i := 0; i < len(newChanges); i++ {
		result, err := Diff(oldChanges[i], newChanges[i], opts...)
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

func (r *DiffListResult) DiffString() string {
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

	const maxPrintDiffs = 3
	var prints = 0
	for _, change := range r.Changes {
		key := change.Old.Key
		opts := []diff.RenderOption{
			diff.WithLeftPadding(1),
		}
		switch {
		case key.IsSecret():
			opts = append(opts, diff.WithMaskPath("data"))
		case key.IsConfigMap():
			opts = append(opts, diff.WithMaskPath("data"))
		}
		renderer := diff.NewRenderer(opts...)

		index++
		b.WriteString(fmt.Sprintf("# %d. %s\n\n", index, key.ReadableString()))
		b.WriteString(renderer.Render(change.Diff.Nodes()))
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
