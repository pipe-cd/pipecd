---
title: "Running stages: ExecuteStage and the DIFF stage"
linkTitle: "Running stages: ExecuteStage and the DIFF stage"
weight: 6
description: >
  Implement the ExecuteStage dispatcher and the FILE_DIFF stage.
---

`ExecuteStage` is where the plugin does the real deployment work. Every stage the plugin declared in `FetchDefinedStages` runs through this one method. The file plugin has three: `FILE_DIFF`, `FILE_SYNC`, and `FILE_ROLLBACK`. Implementing all three at once is a lot, so this chapter wires up the dispatcher and implements the `FILE_DIFF` stage. The next chapter covers `FILE_SYNC` and `FILE_ROLLBACK`.

## Route each stage with ExecuteStage

`piped` calls `ExecuteStage` once per stage and passes the stage name in `input.Request.StageName`. Switch on that name and delegate to a method for each stage. Start with the dispatcher and three empty methods to fill in:

```go
func (p *plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[deployTargetConfig], input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case stageDiff:
		return p.executeStageDiff(ctx, input)
	case stageSync:
		return p.executeStageSync(ctx, input)
	case stageRollback:
		return p.executeStageRollback(ctx, input)
	default:
		return nil, fmt.Errorf("unknown stage: %s", input.Request.StageName)
	}
}

func (p *plugin) executeStageDiff(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}

func (p *plugin) executeStageSync(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}

func (p *plugin) executeStageRollback(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	panic("unimplemented")
}
```

Two things matter before you implement a stage:

- The logs a user sees on the web UI come from the log persister, which you get with `input.Client.LogPersister()`, not from `input.Logger`. `input.Logger` writes to the plugin's own logs, which the user does not see.
- To fail a stage, do not return an error. Set `Status` on the response to `sdk.StageStatusFailure` and return it. Returning an error is for unexpected problems, not for a stage that ran and failed.

## Implement the DIFF stage

The `FILE_DIFF` stage compares the files in the deployment source, which come from the Git repository, against the target directory, and prints the differences to the log. Build it from three small helpers. Unlike the plumbing in the earlier chapters, these carry real logic, so write a test for each one first.

### List the files in a directory

`listFiles` walks a filesystem and returns the set of file paths it contains. Save the following as `main_test.go`, starting with the imports this first test needs:

```go
package main

import (
	"os"
	"testing"
)

func TestListFiles(t *testing.T) {
	path := "./testdata/list_files"
	expectedFiles := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}

	files, err := listFiles(os.DirFS(path))
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	if len(files) != len(expectedFiles) {
		t.Fatalf("expected %d files, got %d", len(expectedFiles), len(files))
	}

	for _, expectedFile := range expectedFiles {
		if _, found := files[expectedFile]; !found {
			t.Errorf("expected file %s not found in the list", expectedFile)
		}
	}
}
```

Create the files the test needs:

```bash
mkdir -p testdata/list_files/subdir
touch testdata/list_files/file1.txt testdata/list_files/file2.txt testdata/list_files/subdir/file3.txt
```

The test fails because `listFiles` does not exist yet. Add it to `plugin.go`:

```go
func listFiles(f fs.FS) (map[string]struct{}, error) {
	files := make(map[string]struct{})

	if err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files[path] = struct{}{}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error walking through files: %w", err)
	}

	return files, nil
}
```

Run `go test` and confirm it passes.

### Find files present in only one directory

`differenceFiles` compares two sets of paths and returns the paths in the first that are not in the second. Add its test to `main_test.go`:

```go
func TestDifferenceFiles(t *testing.T) {
	path1 := "./testdata/difference_files/path1"
	path2 := "./testdata/difference_files/path2"

	expectedDifferences1 := map[string]struct{}{
		"file1.txt": {},
		"file2.txt": {},
	}
	expectedDifferences2 := map[string]struct{}{
		"file3.txt": {},
		"file4.txt": {},
	}

	files1, err := listFiles(os.DirFS(path1))
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	files2, err := listFiles(os.DirFS(path2))
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	differences1 := differenceFiles(files1, files2)
	if !reflect.DeepEqual(differences1, expectedDifferences1) {
		t.Fatalf("expected %v differences, got %v", expectedDifferences1, differences1)
	}

	differences2 := differenceFiles(files2, files1)
	if !reflect.DeepEqual(differences2, expectedDifferences2) {
		t.Fatalf("expected %v differences, got %v", expectedDifferences2, differences2)
	}
}
```

The test uses `reflect`, so add it to the test file's imports. Create the files:

```bash
mkdir -p testdata/difference_files/path1 testdata/difference_files/path2
(cd testdata/difference_files/path1 && touch file0.txt file1.txt file2.txt file5.txt)
(cd testdata/difference_files/path2 && touch file0.txt file3.txt file4.txt file5.txt)
```

Add `differenceFiles` to `plugin.go`:

```go
func differenceFiles(a, b map[string]struct{}) map[string]struct{} {
	differences := make(map[string]struct{})

	for path := range a {
		if _, ok := b[path]; !ok {
			differences[path] = struct{}{}
		}
	}

	return differences
}
```

Run `go test` again and confirm it passes.

### Detect changed file content

Computing a line-by-line or character-by-character diff is involved, and the diffing algorithm is not the point of this tutorial. So `isFileContentDifferent` only reports whether a file's content matches, not what changed. Add its test:

```go
func TestIsFileContentDifferent(t *testing.T) {
	fs1 := os.DirFS("./testdata/difference_file_content/path1")
	fs2 := os.DirFS("./testdata/difference_file_content/path2")

	testCases := []struct {
		name          string
		fsA           fs.FS
		fsB           fs.FS
		path          string
		wantDifferent bool
		wantErr       bool
	}{
		{
			name:          "same content",
			fsA:           fs1,
			fsB:           fs2,
			path:          "file1.txt",
			wantDifferent: false,
			wantErr:       false,
		},
		{
			name:          "different content",
			fsA:           fs1,
			fsB:           fs2,
			path:          "file2.txt",
			wantDifferent: true,
			wantErr:       false,
		},
		{
			name:          "file not found",
			fsA:           fs1,
			fsB:           fs2,
			path:          "file3.txt",
			wantDifferent: false,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotDifferent, err := isFileContentDifferent(tc.fsA, tc.fsB, tc.path)
			if (err != nil) != tc.wantErr {
				t.Fatalf("isFileContentDifferent() error = %v, wantErr %v", err, tc.wantErr)
			}
			if gotDifferent != tc.wantDifferent {
				t.Errorf("isFileContentDifferent() = %v, want %v", gotDifferent, tc.wantDifferent)
			}
		})
	}
}
```

The test's cases hold `fs.FS` values, so add `io/fs` to the test file's imports. Create the files:

```bash
mkdir -p testdata/difference_file_content/path1 testdata/difference_file_content/path2
echo a > testdata/difference_file_content/path1/file1.txt
echo a > testdata/difference_file_content/path1/file2.txt
echo a > testdata/difference_file_content/path2/file1.txt
echo b > testdata/difference_file_content/path2/file2.txt
```

Add `isFileContentDifferent` to `plugin.go`:

```go
func isFileContentDifferent(a, b fs.FS, path string) (bool, error) {
	aFile, err := a.Open(path)
	if err != nil {
		return false, fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer aFile.Close()

	bFile, err := b.Open(path)
	if err != nil {
		return false, fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer bFile.Close()

	aContent, err := io.ReadAll(aFile)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %w", path, err)
	}

	bContent, err := io.ReadAll(bFile)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %w", path, err)
	}

	return !bytes.Equal(aContent, bContent), nil
}
```

Run `go test` and confirm all three tests pass.

### Assemble executeStageDiff

Now combine the three helpers. `executeStageDiff` lists the source and target files, sorts them into added, removed, and changed, and prints each group to the log. The log output is not unit-tested here.

One detail: keep the PipeCD application configuration file out of the diff. Its name is passed as `input.Request.TargetDeploymentSource.ApplicationConfigFilename`, so `delete` it from the source file list.

Replace the empty `executeStageDiff` with the following:

```go
func (p *plugin) executeStageDiff(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	lp := input.Client.LogPersister()

	lp.Info("Listing files in the git repository...")
	sourceFiles, err := listFiles(os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory))
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	delete(sourceFiles, input.Request.TargetDeploymentSource.ApplicationConfigFilename)

	lp.Info("Listing files in the target directory...")
	targetFiles, err := listFiles(os.DirFS(input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Path))
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	addedFiles := differenceFiles(sourceFiles, targetFiles)
	removedFiles := differenceFiles(targetFiles, sourceFiles)

	mergedFiles := maps.Clone(sourceFiles)
	maps.Copy(mergedFiles, targetFiles)

	diffFiles := make(map[string]struct{})
	for path := range mergedFiles {
		if _, ok := addedFiles[path]; ok {
			continue
		}
		if _, ok := removedFiles[path]; ok {
			continue
		}

		different, err := isFileContentDifferent(
			os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory),
			os.DirFS(input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Path),
			path,
		)
		if err != nil {
			return nil, fmt.Errorf("error checking if file content is different: %w", err)
		}

		if different {
			diffFiles[path] = struct{}{}
		}
	}

	lp.Info("Summary of the file diff:")
	lp.Info("--------------------------------")
	lp.Info("Added files:")
	for _, path := range slices.Sorted(maps.Keys(addedFiles)) {
		lp.Info(path)
	}

	lp.Info("--------------------------------")
	lp.Info("Removed files:")
	for _, path := range slices.Sorted(maps.Keys(removedFiles)) {
		lp.Info(path)
	}

	lp.Info("--------------------------------")
	lp.Info("Changed files:")
	for _, path := range slices.Sorted(maps.Keys(diffFiles)) {
		lp.Info(path)
	}

	lp.Info("--------------------------------")
	lp.Success("File diff completed")

	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```

An added file is one that exists in the source but not the target; a removed file is the reverse. A changed file exists in both, so it is neither added nor removed, and its content differs. Sorting each group with `slices.Sorted` keeps the log output stable.

These helpers use several standard-library packages. Make sure `plugin.go` imports them:

```go
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)
```

Build the project and run the tests one more time:

```bash
go build ./...
go test ./...
```

The plugin can now show what a deployment would change. In the next chapter you implement the `FILE_SYNC` stage that applies those changes and the `FILE_ROLLBACK` stage that undoes them.
