---
title: "Running stages: the SYNC and ROLLBACK stages"
linkTitle: "Running stages: the SYNC and ROLLBACK stages"
weight: 7
description: >
  Implement the FILE_SYNC and FILE_ROLLBACK stages.
---

The previous chapter showed what a deployment would change. This chapter makes the change. The `FILE_SYNC` stage applies the files from the deployment source to the target directory, and the `FILE_ROLLBACK` stage puts the target directory back the way it was. Both build on the same two operations: copy the files that belong there, and remove the files that no longer do.

## Implement the SYNC stage

`FILE_SYNC` synchronizes the target directory with the Git repository in two steps. First it copies every file from the deployment source into the target directory. Then it removes any file in the target directory that no longer exists in the source. Build one helper for each step, test-first, then combine them.

### Copy files

`copyFiles` copies every file from a source filesystem into a target directory. It also takes a set of paths to exclude, which is how the plugin keeps the PipeCD application configuration file out of the target directory. Add the test to `main_test.go`. It writes into `t.TempDir()` so it does not touch your working tree, and it reuses the `testdata/list_files` directory from the previous chapter:

```go
func TestCopyFiles(t *testing.T) {
	srcDir := "testdata/list_files"
	dstDir := t.TempDir()

	if err := copyFiles(dstDir, os.DirFS(srcDir), map[string]struct{}{"file2.txt": {}}); err != nil {
		t.Fatalf("copyFiles() error = %v", err)
	}

	srcFiles, err := listFiles(os.DirFS(srcDir))
	if err != nil {
		t.Fatalf("listFiles() on source dir failed: %v", err)
	}

	dstFiles, err := listFiles(os.DirFS(dstDir))
	if err != nil {
		t.Fatalf("listFiles() on dest dir failed: %v", err)
	}

	delete(srcFiles, "file2.txt") // file2.txt is excluded

	if !reflect.DeepEqual(srcFiles, dstFiles) {
		t.Errorf("copied files list differs. got %v, want %v", dstFiles, srcFiles)
	}

	for path := range srcFiles {
		srcContent, err := os.ReadFile(filepath.Join(srcDir, path))
		if err != nil {
			t.Fatalf("failed to read source file %s: %v", path, err)
		}

		dstContent, err := os.ReadFile(filepath.Join(dstDir, path))
		if err != nil {
			t.Fatalf("failed to read destination file %s: %v", path, err)
		}

		if !bytes.Equal(srcContent, dstContent) {
			t.Errorf("content of %s is different", path)
		}
	}
}
```

The test uses `bytes` and `path/filepath`, so add them to the test file's imports. Now add `copyFiles` to `plugin.go`:

```go
func copyFiles(dstDir string, files fs.FS, exclude map[string]struct{}) error {
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if _, ok := exclude[path]; ok {
			return nil
		}

		dstPath := filepath.Join(dstDir, path)

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}

		srcFile, err := files.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return nil
	}

	if err := fs.WalkDir(files, ".", walkDirFunc); err != nil {
		return fmt.Errorf("walking through files: %w", err)
	}

	return nil
}
```

The walk skips directories rather than creating them. Each file's parent directory is created with `os.MkdirAll` right before the file is written, so empty directories are never copied, and no file is written without its parent existing.

Run `go test` and confirm it passes.

### Remove files

`removeFiles` deletes the files in the target directory that are no longer in the source. It reuses `listFiles` and `differenceFiles` from the previous chapter and takes the same exclude set. Add the test to `main_test.go`. It copies the prepared `dst_before` directory into `t.TempDir()` first, because the test modifies the target:

```go
func TestRemoveFiles(t *testing.T) {
	srcDir := "testdata/remove_files/src"
	dstDir := t.TempDir()
	if err := copyFiles(dstDir, os.DirFS("testdata/remove_files/dst_before"), nil); err != nil {
		t.Fatalf("failed to copy dst_before: %v", err)
	}

	if err := removeFiles(dstDir, os.DirFS(srcDir), map[string]struct{}{"excluded_file.txt": {}}); err != nil {
		t.Fatalf("removeFiles() error = %v", err)
	}

	srcFS := os.DirFS(srcDir)
	expectedFiles, err := listFiles(srcFS)
	if err != nil {
		t.Fatalf("failed to list files in src dir: %v", err)
	}

	delete(expectedFiles, "excluded_file.txt")

	dstFiles, err := listFiles(os.DirFS(dstDir))
	if err != nil {
		t.Fatalf("failed to list files in dst dir: %v", err)
	}

	if !reflect.DeepEqual(dstFiles, expectedFiles) {
		t.Errorf("file list differs. got %v, want %v", dstFiles, expectedFiles)
	}

	if _, err := os.Stat(filepath.Join(dstDir, "file_to_remove.txt")); !os.IsNotExist(err) {
		t.Errorf("file_to_remove.txt was not removed")
	}
}
```

Create the test data:

```bash
mkdir -p testdata/remove_files/src/subdir testdata/remove_files/dst_before/subdir
touch testdata/remove_files/src/file1.txt testdata/remove_files/src/subdir/file2.txt testdata/remove_files/src/excluded_file.txt
touch testdata/remove_files/dst_before/file1.txt testdata/remove_files/dst_before/subdir/file2.txt testdata/remove_files/dst_before/file_to_remove.txt testdata/remove_files/dst_before/excluded_file.txt
```

Add `removeFiles` to `plugin.go`:

```go
func removeFiles(dstDir string, files fs.FS, exclude map[string]struct{}) error {
	sourceFiles, err := listFiles(files)
	if err != nil {
		return fmt.Errorf("listing files: %w", err)
	}

	for path := range exclude {
		delete(sourceFiles, path)
	}

	dstFiles, err := listFiles(os.DirFS(dstDir))
	if err != nil {
		return fmt.Errorf("listing files: %w", err)
	}

	removedFiles := differenceFiles(dstFiles, sourceFiles)

	for path := range removedFiles {
		if err := os.Remove(filepath.Join(dstDir, path)); err != nil {
			return fmt.Errorf("removing file %s: %w", path, err)
		}
	}

	return nil
}
```

Run `go test` again and confirm it passes.

### Assemble executeStageSync

Combine the two helpers. `executeStageSync` copies the deployment source files into the target directory and then removes anything left behind. Both steps exclude the application configuration file, whose name comes from `input.Request.TargetDeploymentSource.ApplicationConfigFilename`. Replace the empty `executeStageSync` with the following:

```go
func (p *plugin) executeStageSync(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	lp := input.Client.LogPersister()

	lp.Info("Copying files to the target directory...")
	if err := copyFiles(
		input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.TargetDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error copying files: %w", err)
	}

	lp.Info("Removing files which are not in the git repository from the target directory...")
	if err := removeFiles(
		input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.TargetDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error removing files: %w", err)
	}

	lp.Success("File sync completed")
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```

## Implement the ROLLBACK stage

`FILE_ROLLBACK` runs when a deployment fails and the target directory has to go back to its previous state. It does exactly what `executeStageSync` does, with one change: it reads from `input.Request.RunningDeploymentSource` instead of `input.Request.TargetDeploymentSource`.

The two sources are the difference between the two stages. `TargetDeploymentSource` is the deployment being rolled out, the new state. `RunningDeploymentSource` is the deployment that was already running, the state to restore. Copying the running source back over the target directory, then removing anything not in it, returns the directory to how it was before the failed deployment.

Replace the empty `executeStageRollback` with the following:

```go
func (p *plugin) executeStageRollback(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
	lp := input.Client.LogPersister()

	lp.Info("Copying files to the target directory...")
	if err := copyFiles(
		input.Request.RunningDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.RunningDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.RunningDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error copying files: %w", err)
	}

	lp.Info("Removing files which are not in the git repository from the target directory...")
	if err := removeFiles(
		input.Request.RunningDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.RunningDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.RunningDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error removing files: %w", err)
	}

	lp.Success("File rollback completed")
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```

These helpers add one more standard-library package. Make sure `plugin.go` imports `path/filepath`:

```go
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)
```

Build the project and run the tests:

```bash
go build ./...
go test ./...
```

All six methods of the deployment plugin now have a real implementation. The plugin can report versions, plan stages for either sync strategy, and run the diff, sync, and rollback stages. In the next chapter you build the plugin into a binary and run it against a local `piped`.
