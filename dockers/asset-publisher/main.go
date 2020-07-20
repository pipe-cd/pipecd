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
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var (
	org        = flag.String("org", "pipe-cd", "The name of GitHub organization")
	repo       = flag.String("repo", "pipe", "The name of GitHub repository")
	releaseTag = flag.String("release-tag", "", "The release tag where asset should be uploaded to")
	assetName  = flag.String("asset-name", "", "The name of the asset")
	assetFile  = flag.String("asset-file", "", "The path to the asset file")
	tokenFile  = flag.String("token-file", "", "The path to the token file")

	timeout = 15 * time.Minute
)

func main() {
	flag.Parse()

	tokenBytes, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		log.Fatalf("failed to read token file at %s: %v", *tokenFile, err)
	}
	token := string(tokenBytes)
	token = strings.TrimSpace(token)

	asset, err := os.Open(*assetFile)
	if err != nil {
		log.Fatalf("failed to open asset file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Find the release.
	release, _, err := client.Repositories.GetReleaseByTag(
		ctx,
		*org,
		*repo,
		*releaseTag,
	)
	if err != nil {
		log.Fatalf("failed to get release %s, %s",
			*releaseTag,
			strings.ReplaceAll(err.Error(), token, "TOKEN"),
		)
	}

	log.Printf("start uploading %s at %s to release %s", *assetName, *assetFile, *releaseTag)
	_, _, err = client.Repositories.UploadReleaseAsset(
		ctx,
		*org,
		*repo,
		release.GetID(),
		&github.UploadOptions{
			Name: *assetName,
		},
		asset,
	)
	if err != nil {
		log.Fatalf("failed to upload asset: %s", strings.ReplaceAll(err.Error(), token, "TOKEN"))
	}

	log.Printf("successfully uploaded asset %s", *assetName)
}
