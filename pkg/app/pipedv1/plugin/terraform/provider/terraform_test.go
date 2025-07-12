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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanHasChangeRegex(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "older than v1.5.0",
			input:    "Plan: 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 1 to add, 2 to change, 3 to destroy.", "", "1", "2", "3"},
		},
		{
			name:     "later than v1.5.0",
			input:    "Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.", "0", "1", "2", "3"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, planHasChangeRegex.FindStringSubmatch(tc.input))
		})
	}
}

func TestParsePlanResult(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       string
		expected    PlanResult
		expectedErr bool
	}{
		{
			name:        "older than v1.5.0",
			input:       `Plan: 1 to add, 2 to change, 3 to destroy.`,
			expected:    PlanResult{Adds: 1, Changes: 2, Destroys: 3, HasStateChanges: true},
			expectedErr: false,
		},
		{
			name:        "later than v1.5.0",
			input:       `Plan: 1 to import, 1 to add, 2 to change, 3 to destroy.`,
			expected:    PlanResult{Imports: 1, Adds: 1, Changes: 2, Destroys: 3, HasStateChanges: true},
			expectedErr: false,
		},
		{
			name:        "Invalid number of changes",
			input:       `Plan: a to add, 2 to change, 3 to destroy.`,
			expectedErr: true,
		},
		{
			name:        "Invalid plan result output",
			input:       `Plan: 1 to add, 2 to change.`,
			expectedErr: true,
		},
		{
			name: "Changes to outputs",
			input: `terraform init -no-color
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
			expected:    PlanResult{HasStateChanges: true},
			expectedErr: false,
		},
		{
			name: "Refactor", // when using moved blocks or removed blocks
			input: `terraform init -no-color
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

 # google_compute_global_address.xxx will no longer be managed by Terraform, but will not be destroyed
 # (destroy = false is set in the configuration)
 . resource "google_compute_global_address" "xxx" {
        id                 = "xxxx"
        name               = "xxxx"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.`,
			expected:    PlanResult{HasStateChanges: true},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := parsePlanResult(tc.input, false)
			assert.Equal(t, tc.expectedErr, err != nil)
			result.PlanOutput = ""
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRender(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		expected    string
		planResult  *PlanResult
		expectedErr bool
	}{
		{
			name: "success",
			planResult: &PlanResult{
				Imports:  1,
				Adds:     2,
				Changes:  3,
				Destroys: 4,
				PlanOutput: `
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:
  + resource "test-add" "test" {
      + id    = (known after apply)
    }
  - resource "test-del" "test" {
      + id    = "foo"
    }

Plan: 1 to import, 2 to add, 3 to change, 4 to destroy.
`,
			},
			expected: `    resource "test-add" "test" {
+       id    = (known after apply)
    }
    resource "test-del" "test" {
+       id    = "foo"
    }
Plan: 1 to import, 2 to add, 3 to change, 4 to destroy.
`,
			expectedErr: false,
		},
		{
			name: "New outputs",
			planResult: &PlanResult{
				HasStateChanges: true,
				PlanOutput: `terraform init -no-color
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
			},
			expected: "",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tc.planResult.Render()
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
