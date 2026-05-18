---
title: "ExecuteStage Implementation: DIFF Stage"
weight: 17
description: >
  Implementing the DIFF stage to log and preview local file changes.
---

In the `FILE_DIFF` stage, we walk both the target directory on Piped's host and the Git repository's application directory to detect added, modified, or deleted files, printing the results to the log.

We will divide the implementation into three helper functions:

1. **`listFiles`**: Lists all relative file paths under a directory.
2. **`differenceFiles`**: Finds paths present in one file set but missing in another.
3. **`isFileContentDifferent`**: Checks if the content of a file has changed.

---

### 1. Listing Files (`listFiles`)

Let's start by writing a test. Create or append to a file named `main_test.go`:

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

To make the test pass, create the test files and directories:

```console
$ mkdir -p testdata/list_files/subdir
$ touch testdata/list_files/file1.txt testdata/list_files/file2.txt testdata/list_files/subdir/file3.txt
```

Now, implement `listFiles` using Go's standard `io/fs` package:

```go
import (
	"fmt"
	"io/fs"
)

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

Run `go test` to confirm it passes.

---

### 2. Finding File Key Differences (`differenceFiles`)

Next, let's write a test to find files that exist in one map but not in another:

```go
import "reflect"

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

Create the necessary test data:

```console
$ mkdir -p testdata/difference_files/path1 testdata/difference_files/path2
$ touch testdata/difference_files/path1/file0.txt testdata/difference_files/path1/file1.txt testdata/difference_files/path1/file2.txt testdata/difference_files/path1/file5.txt
$ touch testdata/difference_files/path2/file0.txt testdata/difference_files/path2/file3.txt testdata/difference_files/path2/file4.txt testdata/difference_files/path2/file5.txt
```

Implement `differenceFiles` in `main.go`:

```go
// differenceFiles compares map a and map b, returning keys that exist in a but not in b.
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

---

### 3. Comparing File Content (`isFileContentDifferent`)

For educational simplicity, we will check if the content has changed using direct byte comparison rather than implementing a full line-by-line diff.

Write the content comparison test:

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

Create the content test data:

```console
$ mkdir -p testdata/difference_file_content/path1 testdata/difference_file_content/path2
$ echo a > testdata/difference_file_content/path1/file1.txt
$ echo a > testdata/difference_file_content/path1/file2.txt
$ echo a > testdata/difference_file_content/path2/file1.txt
$ echo b > testdata/difference_file_content/path2/file2.txt
```

Implement `isFileContentDifferent` in `main.go`:

```go
import (
	"bytes"
	"io"
)

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

Verify that all tests pass using `go test`.

---

### 4. Implementing `executeStageDiff`

Now, let's combine these three helpers into `executeStageDiff`.

> [!IMPORTANT]
> Make sure to **exclude** PipeCD's application configuration file (`app.pipecd.yaml`) from the file comparison. Piped passes this filename under `input.Request.TargetDeploymentSource.ApplicationConfigFilename`, so we can simply use the Go builtin `delete` to remove it from our file lists.

```go
import (
	"maps"
	"slices"
)

func (plugin) executeStageDiff(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	lp := input.Client.LogPersister()

	lp.Info("Listing files in the Git repository...")
	sourceFiles, err := listFiles(os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory))
	if err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	// Exclude the application config file
	delete(sourceFiles, input.Request.TargetDeploymentSource.ApplicationConfigFilename)

	lp.Info("Listing files in the target deployment directory...")
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
			return nil, fmt.Errorf("error checking file diff for %s: %w", path, err)
		}

		if different {
			diffFiles[path] = struct{}{}
		}
	}

	// Output results to the LogPersister
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

	lp.Success("File diff completed successfully")

	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```
