---
title: "ExecuteStage Implementation: SYNC Stage"
weight: 18
description: >
  Implementing the SYNC stage to copy and synchronize local files.
---

In the `FILE_SYNC` stage, we execute the actual deployment operations:

1. Copying all files from the Git repository's application directory to the target deployment path.
2. Deleting any orphaned files in the target path that are no longer tracked in Git.

---

### 1. Copying Files (`copyFiles`)

First, write a test in `main_test.go`. Since this involves disk manipulation, we will isolate it using `t.TempDir()`:

```go
import "path/filepath"

func TestCopyFiles(t *testing.T) {
	srcDir := "testdata/list_files"
	dstDir := t.TempDir()

	// Exclude file2.txt for testing
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

	delete(srcFiles, "file2.txt") // file2.txt was excluded

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

Now, implement `copyFiles` in `main.go`. We will create parent directories dynamically as files are written:

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

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
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

---

### 2. Deleting Orphaned Files (`removeFiles`)

Write the deletion test in `main_test.go`:

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

Create the required test data:

```console
$ mkdir -p testdata/remove_files/src/subdir testdata/remove_files/dst_before/subdir
$ touch testdata/remove_files/src/file1.txt testdata/remove_files/src/subdir/file2.txt testdata/remove_files/src/excluded_file.txt
$ touch testdata/remove_files/dst_before/file1.txt testdata/remove_files/dst_before/subdir/file2.txt testdata/remove_files/dst_before/file_to_remove.txt testdata/remove_files/dst_before/excluded_file.txt
```

Implement `removeFiles` in `main.go`:

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

Run `go test` to confirm everything works properly.

---

### 3. Implementing `executeStageSync`

Now, let's orchestrate these in `executeStageSync`:

```go
func (plugin) executeStageSync(ctx context.Context, input *sdk.ExecuteStageInput[applicationConfig]) (*sdk.ExecuteStageResponse, error) {
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

	lp.Info("Removing orphaned files from target directory...")
	if err := removeFiles(
		input.Request.TargetDeploymentSource.ApplicationConfig.Spec.Path,
		os.DirFS(input.Request.TargetDeploymentSource.ApplicationDirectory),
		map[string]struct{}{
			input.Request.TargetDeploymentSource.ApplicationConfigFilename: {},
		},
	); err != nil {
		return nil, fmt.Errorf("error removing files: %w", err)
	}

	lp.Success("File synchronization completed successfully")
	return &sdk.ExecuteStageResponse{
		Status: sdk.StageStatusSuccess,
	}, nil
}
```
