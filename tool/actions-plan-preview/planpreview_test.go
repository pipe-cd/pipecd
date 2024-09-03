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

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeCommentBody(t *testing.T) {
	// NOTE: Skip this temporarily since this test result is diiffrent on CI and Local env.
	if os.Getenv("CI") != "" {
		t.Skip()
	}

	testcases := []struct {
		name     string
		event    githubEvent
		result   PlanPreviewResult
		expected string
	}{
		{
			name: "no changes",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result:   PlanPreviewResult{},
			expected: "testdata/comment-no-changes.txt",
		},
		{
			name: "only changed app",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							Env:                  "env-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
			},
			expected: "testdata/comment-only-changed-app.txt",
		},
		{
			name: "has no diff apps",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-2",
							ApplicationName:      "app-name-2",
							ApplicationURL:       "app-url-2",
							Env:                  "env-2",
							ApplicationKind:      "app-kind-2",
							ApplicationDirectory: "app-dir-2",
						},
						SyncStrategy: "QUICK_SYNC",
						PlanSummary:  "",
						PlanDetails:  "",
						NoChange:     true,
					},
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							Env:                  "env-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-3",
							ApplicationName:      "app-name-3",
							ApplicationURL:       "app-url-3",
							Env:                  "env-3",
							ApplicationKind:      "app-kind-3",
							ApplicationDirectory: "app-dir-3",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "",
						PlanDetails:  "",
						NoChange:     true,
					},
				},
			},
			expected: "testdata/comment-has-no-diff-apps.txt",
		},
		{
			name: "no env",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
			},
			expected: "testdata/comment-no-env.txt",
		},
		{
			name: "has failed app",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							Env:                  "env-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
				FailureApplications: []FailureApplication{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-2",
							ApplicationName:      "app-name-2",
							ApplicationURL:       "app-url-2",
							Env:                  "env-2",
							ApplicationKind:      "app-kind-2",
							ApplicationDirectory: "app-dir-2",
						},
						Reason:      "reason-2",
						PlanDetails: "",
					},
				},
			},
			expected: "testdata/comment-has-failed-app.txt",
		},
		{
			name: "has failed piped",
			event: githubEvent{
				HeadCommit: "abc",
			},
			result: PlanPreviewResult{
				Applications: []ApplicationResult{
					{
						ApplicationInfo: ApplicationInfo{
							ApplicationID:        "app-id-1",
							ApplicationName:      "app-name-1",
							ApplicationURL:       "app-url-1",
							Env:                  "env-1",
							ApplicationKind:      "app-kind-1",
							ApplicationDirectory: "app-dir-1",
						},
						SyncStrategy: "PIPELINE",
						PlanSummary:  "plan-summary-1",
						PlanDetails:  "plan-details-1",
						NoChange:     false,
					},
				},
				FailurePipeds: []FailurePiped{
					{
						PipedInfo: PipedInfo{
							PipedID:  "piped-id-1",
							PipedURL: "piped-url-1",
						},
						Reason: "piped-reason-1",
					},
				},
			},
			expected: "testdata/comment-has-failed-piped.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := testdata.ReadFile(tc.expected)
			require.NoError(t, err)

			got := makeCommentBody(&tc.event, &tc.result)
			assert.Equal(t, string(expected), got)
		})
	}
}

func TestGenerateTerraformShortPlanDetails(t *testing.T) {
	testcases := []struct {
		name        string
		planDetails string
		want        string
		wantErr     bool
	}{
		{
			name: "simple",
			planDetails: `terraform init -no-color
Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "x.xx.x"...
- Installing hashicorp/google vx.xx.x...
- Installed hashicorp/google vx.xx.x (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
terraform plan -lock=false -detailed-exitcode -no-colorgoogle_compute_global_address.xxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxx]
google_service_account.xxxxx: Refreshing state... [id=projects/xxxx/serviceAccounts/xxxxx@xxxxx.iam.gserviceaccount.com]
google_compute_global_address.xxxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxxxx]
google_dns_record_set.xxxxx: Refreshing state... [id=xxxxx/A]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # google_dns_record_set.xxx will be created
  + resource "google_dns_record_set" "xxxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxx"
      + name         = "xxxx.xxxx.xxxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.x",
        ]
      + ttl          = xxx
      + type         = "A"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
1 to add, 0 to change, 0 to destroy`,
			want: `Terraform will perform the following actions:

  # google_dns_record_set.xxx will be created
  + resource "google_dns_record_set" "xxxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxx"
      + name         = "xxxx.xxxx.xxxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.x",
        ]
      + ttl          = xxx
      + type         = "A"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
1 to add, 0 to change, 0 to destroy`,
			wantErr: false,
		},
		{
			name: "warning",
			planDetails: `terraform init -no-color
Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "x.xx.x"...
- Installing hashicorp/google vx.xx.x...
- Installed hashicorp/google vx.xx.x (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.


Warning: Version constraints inside provider configuration blocks are deprecated

  on xxx.tf line 14, in provider "google":
  14:   version     = "x.xx.x"

Terraform 0.13 and earlier allowed provider version constraints inside the
provider configuration block, but that is now deprecated and will be removed
in a future version of Terraform. To silence this warning, move the provider
version constraint into the required_providers block.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
terraform plan -lock=false -detailed-exitcode -no-colorgoogle_compute_global_address.xxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxx]
google_service_account.xxxxx: Refreshing state... [id=projects/xxxx/serviceAccounts/xxxxx@xxxxx.iam.gserviceaccount.com]
google_compute_global_address.xxxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxxxx]
google_dns_record_set.xxxxx: Refreshing state... [id=xxxxx/A]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # google_dns_record_set.xxx will be created
  + resource "google_dns_record_set" "xxxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxx"
      + name         = "xxxx.xxxx.xxxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.x",
        ]
      + ttl          = xxx
      + type         = "A"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Warning: Version constraints inside provider configuration blocks are deprecated

  on xxx.tf line 14, in provider "google":
  14:   version     = "x.xx.x"

Terraform 0.13 and earlier allowed provider version constraints inside the
provider configuration block, but that is now deprecated and will be removed
in a future version of Terraform. To silence this warning, move the provider
version constraint into the required_providers block.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
1 to add, 0 to change, 0 to destroy`,
			want: `Terraform will perform the following actions:

  # google_dns_record_set.xxx will be created
  + resource "google_dns_record_set" "xxxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxx"
      + name         = "xxxx.xxxx.xxxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.x",
        ]
      + ttl          = xxx
      + type         = "A"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Warning: Version constraints inside provider configuration blocks are deprecated

  on xxx.tf line 14, in provider "google":
  14:   version     = "x.xx.x"

Terraform 0.13 and earlier allowed provider version constraints inside the
provider configuration block, but that is now deprecated and will be removed
in a future version of Terraform. To silence this warning, move the provider
version constraint into the required_providers block.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
1 to add, 0 to change, 0 to destroy`,
			wantErr: false,
		},
		{
			name: "chages made outside of Terraform",
			planDetails: `terraform init -no-color
Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "x.xx.x"...
- Installing hashicorp/google vx.xx.x...
- Installed hashicorp/google vx.xx.x (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
terraform plan -lock=false -detailed-exitcode -no-colorgoogle_compute_global_address.xxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxx]
google_service_account.xxxxx: Refreshing state... [id=projects/xxxx/serviceAccounts/xxxxx@xxxxx.iam.gserviceaccount.com]
google_compute_global_address.xxxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxxxx]
google_dns_record_set.xxxxx: Refreshing state... [id=xxxxx/A]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the
last "terraform apply":

  # google_project_iam_binding.xxxx has been changed
  ~ resource "google_project_iam_binding" "xxxx" {
      ~ etag    = "xxxxx" -> "yyyyy"
        id      = "xxxx-xxxx/roles/cloudsql.client"
      ~ members = [
          + "serviceAccount:xxxxxx@xxxxxx.iam.gserviceaccount.com",
            # (1 unchanged element hidden)
        ]
        # (2 unchanged attributes hidden)
    }

Unless you have made equivalent changes to your configuration, or ignored the
relevant attributes using ignore_changes, the following plan may include
actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create
  ~ update in-place

Terraform will perform the following actions:

  # google_dns_managed_zone.xxxxx will be created
  + resource "google_dns_managed_zone" "xxxxx" {
      + description   = "xxxxxxxxx"
      + dns_name      = "xxx.xxxx.xxx."
      + force_destroy = xxxxxx
      + id            = (known after apply)
      + name          = "xxxxxxxx"
      + name_servers  = (known after apply)
      + project       = (known after apply)
      + visibility    = "xxxxx"
    }

  # google_dns_record_set.xxxx will be created
  + resource "google_dns_record_set" "xxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxxxxx"
      + name         = "xxx.xxxx.xxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.xx",
        ]
      + ttl          = xxx
      + type         = "A"
    }

  # google_project_iam_binding.xxx will be updated in-place
  ~ resource "google_project_iam_binding" "xxx" {
        id      = "xxxxxx/roles/cloudsql.client"
      ~ members = [
          - "serviceAccount:xxxxxxxx@xxxxx.iam.gserviceaccount.com",
            # (1 unchanged element hidden)
        ]
        # (3 unchanged attributes hidden)
    }

Plan: 2 to add, 1 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
2 to add, 1 to change, 0 to destroy`,
			want: `Terraform will perform the following actions:

  # google_dns_managed_zone.xxxxx will be created
  + resource "google_dns_managed_zone" "xxxxx" {
      + description   = "xxxxxxxxx"
      + dns_name      = "xxx.xxxx.xxx."
      + force_destroy = xxxxxx
      + id            = (known after apply)
      + name          = "xxxxxxxx"
      + name_servers  = (known after apply)
      + project       = (known after apply)
      + visibility    = "xxxxx"
    }

  # google_dns_record_set.xxxx will be created
  + resource "google_dns_record_set" "xxxx" {
      + id           = (known after apply)
      + managed_zone = "xxxxxxx"
      + name         = "xxx.xxxx.xxx."
      + project      = (known after apply)
      + rrdatas      = [
          + "xx.xxx.xx.xx",
        ]
      + ttl          = xxx
      + type         = "A"
    }

  # google_project_iam_binding.xxx will be updated in-place
  ~ resource "google_project_iam_binding" "xxx" {
        id      = "xxxxxx/roles/cloudsql.client"
      ~ members = [
          - "serviceAccount:xxxxxxxx@xxxxx.iam.gserviceaccount.com",
            # (1 unchanged element hidden)
        ]
        # (3 unchanged attributes hidden)
    }

Plan: 2 to add, 1 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
2 to add, 1 to change, 0 to destroy`,
			wantErr: false,
		},
		{
			name: "moved resources",
			planDetails: `terraform init -no-color
Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "x.xx.x"...
- Installing hashicorp/google vx.xx.x...
- Installed hashicorp/google vx.xx.x (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
terraform plan -lock=false -detailed-exitcode -no-colorgoogle_compute_global_address.xxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxx]
google_service_account.xxxxx: Refreshing state... [id=projects/xxxx/serviceAccounts/xxxxx@xxxxx.iam.gserviceaccount.com]
google_compute_global_address.xxxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxxxx]
google_dns_record_set.xxxxx: Refreshing state... [id=xxxxx/A]

Terraform will perform the following actions:

  # google_dns_record_set.xxx has moved to google_dns_record_set.xxx
    resource "google_compute_global_forwarding_rule" "xxx" {
        id           = "xxxx"
        managed_zone = "xxxx"
        name         = "xxxx.xxxx.xxxx."
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.`,
			want: `Terraform will perform the following actions:

  # google_dns_record_set.xxx has moved to google_dns_record_set.xxx
    resource "google_compute_global_forwarding_rule" "xxx" {
        id           = "xxxx"
        managed_zone = "xxxx"
        name         = "xxxx.xxxx.xxxx."
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.`,
			wantErr: false,
		},
		{
			name: "Only new outputs",
			planDetails: `terraform init -no-color
Initializing the backend...

Successfully configured the backend "gcs"! Terraform will automatically
use this backend unless the backend configuration changes.

Initializing provider plugins...
- Finding hashicorp/google versions matching "x.xx.x"...
- Installing hashicorp/google vx.xx.x...
- Installed hashicorp/google vx.xx.x (signed by HashiCorp)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
terraform plan -lock=false -detailed-exitcode -no-colorgoogle_compute_global_address.xxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxx]
google_service_account.xxxxx: Refreshing state... [id=projects/xxxx/serviceAccounts/xxxxx@xxxxx.iam.gserviceaccount.com]
google_compute_global_address.xxxx: Refreshing state... [id=projects/xxxx/global/addresses/xxxxxx]
google_dns_record_set.xxxxx: Refreshing state... [id=xxxxx/A]

Changes to Outputs:
  + global_address = xxxx

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.`,
			want: `Changes to Outputs:
  + global_address = xxxx

You can apply this plan to save these new output values to the Terraform
state, without changing any real infrastructure.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.`,
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := generateTerraformShortPlanDetails(tc.planDetails)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
