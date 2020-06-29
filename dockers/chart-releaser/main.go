// Copyright 2020 The PipeCD Authors.
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
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	manifestsDir    string
	bucket          string
	credentialsFile string
	indexFileName   = "index.yaml"
	timeout         = 10 * time.Minute

	charts = []string{
		"helloworld",
		"piped",
	}
)

func init() {
	flag.StringVar(&manifestsDir, "manifests-dir", "", "The path to the manifests directory")
	flag.StringVar(&bucket, "bucket", "", "The name of GCS bucket where to put charts")
	flag.StringVar(&credentialsFile, "credentials-file", "", "The path to the credentials file used while commucating with GCS")
	flag.Parse()

	if manifestsDir == "" {
		log.Fatalf("manifests-dir is required")
	}
	if bucket == "" {
		log.Fatalf("bucket is required")
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// Initialize gcs client.
	var options []option.ClientOption
	if credentialsFile != "" {
		options = append(options, option.WithCredentialsFile(credentialsFile))
	}
	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		log.Fatalf("Unable to create GCS client: %v", err)
	}

	// Make a temporary working directory.
	workingDir, err := ioutil.TempDir("", "charts")
	if err != nil {
		log.Fatalf("Unable to create a temporary working directory: %v", err)
	}
	//defer os.RemoveAll(workingDir)
	log.Printf("Successfully created a temporary working directory: %s", workingDir)

	// Download current index.yaml file from storage.
	if err := downloadIndexFile(ctx, client, filepath.Join(workingDir, indexFileName)); err != nil {
		log.Fatalf("Unable to download current index file: %v", err)
	}
	log.Printf("Successfully downloaded current index file")

	// Package new charts.
	for _, chart := range charts {
		chartPath := filepath.Join(manifestsDir, chart)
		if err := packageHelmChart(ctx, chartPath, workingDir); err != nil {
			log.Fatalf("Unable to package chart %s: %v", chart, err)
		}
		log.Printf("Successfully packaged chart %s", chart)
	}

	// Generate new index.yaml file by merging new charts.
	if err := generateNewIndex(ctx, workingDir); err != nil {
		log.Fatalf("Unable to update index file: %v", err)
	}
	log.Printf("Successfully updated index file")

	// Start uploading new packages and new index.yaml file.
	err = filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if path == workingDir {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		if info.Name() == indexFileName {
			return nil
		}
		log.Printf("Start uploading package %s...", info.Name())
		return storeFile(ctx, client, path)
	})
	if err != nil {
		log.Fatalf("Unable to store chart packages: %v", err)
	}

	if err := storeFile(ctx, client, filepath.Join(workingDir, indexFileName)); err != nil {
		log.Fatalf("Unable to store index file: %v", err)
	}

	log.Printf("Successfully stored all packages and index file")
}

func packageHelmChart(ctx context.Context, chartPath, dest string) error {
	args := []string{"package", chartPath, "--destination", dest}
	cmd := exec.CommandContext(ctx, "helm", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to package: %s (%v)", string(out), err)
	}
	return nil
}

func generateNewIndex(ctx context.Context, dir string) error {
	args := []string{"repo", "index", ".", "--merge", indexFileName}
	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update Helm index: %s (%v)", string(out), err)
	}
	return nil
}

func downloadIndexFile(ctx context.Context, client *storage.Client, dest string) error {
	rc, err := client.Bucket(bucket).Object(indexFileName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, content, 0644)
}

func storeFile(ctx context.Context, client *storage.Client, path string) error {
	name := filepath.Base(path)
	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)
	defer wc.Close()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = wc.Write(content)
	return err
}
