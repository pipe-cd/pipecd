#!/usr/bin/env bash

# Copyright 2024 The PipeCD Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

QUIET='false'
CLEANUP_BRANCH='false'
while getopts 'qc' opt; do
  case $opt in
    q)
      QUIET='true' ;;
    c)
      CLEANUP_BRANCH='true' ;;
    ?)
      echo "Usage: hack/cherry-pick.sh [-q] [-c] <branch> <pull-requests>" >&2
      exit 1 ;;
  esac
done
shift "$(($OPTIND -1))"

BRANCH=$1; shift 1
PULL_REQUESTS=($@)

# Check gh authentication
echo "+++ Checking gh authentication..."
gh auth status

COMMIT_HASHS=()
for pull in "${PULL_REQUESTS[@]}"; do
  hash="$(gh pr view ${pull} --json mergeCommit --jq .mergeCommit.oid)"
  COMMIT_HASHS+=($hash)
done

function join { local IFS="$1"; shift; echo "$*"; }
PULL_DASH=$(join - "${PULL_REQUESTS[@]/#/#}")
PULL_SUBJ=$(join " " "${PULL_REQUESTS[@]/#/#}")
NEWBRANCH="cherry-pick-${PULL_DASH}-to-${BRANCH}"

# Update all remote branches
echo "+++ Updating remote branches..."
git remote update
echo

# Check whether remote branch exists
echo "+++ Checking whether remote branch exists..."
BASE_TAG=$(echo ${BRANCH} | sed -e 's/.*-//' -e 's/x/0/g')
IS_EXIST=$(git ls-remote --heads origin ${BRANCH} | wc -l)
if [ $IS_EXIST -ne 1 ]; then
  git checkout -b ${BRANCH} ${BASE_TAG}
  git push origin ${BRANCH}
else
  echo "Branch ${BRANCH} already exists"
fi
echo

# Create local branch
echo "+++ Creating a local branch..."
git checkout -b ${NEWBRANCH} "origin/${BRANCH}"
echo

# Cherry-pick pull requests
COMMITS=$(join " " "${COMMIT_HASHS[@]}")
echo "+++ Cherry-picking pull requests"
git cherry-pick ${COMMITS}
echo

# Check whether to push commits and create a pull request or not
if ! ${QUIET}; then
  read -p "+++ Do you push commits and create a pull request? [y/n] " -r
  if ! [[ "${REPLY}" =~ ^[yY]$  ]]; then
    echo "Skipped." >&2
    exit 0
  fi
fi

# Push commits to remote branch
echo "+++ Pushing commits to remote branch..."
git push origin ${NEWBRANCH}
echo

# Create a pull request
echo "+++ Creating a pull request..."
pull_title="Cherry-pick ${PULL_SUBJ}"
pull_body=$(cat <<EOF
**What this PR does / why we need it**:
Cherry pick of ${PULL_SUBJ}.
EOF
)
gh pr create --title="${pull_title}" --body="${pull_body}" --head "${NEWBRANCH}" --base "${BRANCH}"

# Delete a working branch
if ${CLEANUP_BRANCH}; then
  echo "+++ Deleting a working branch..."
  git checkout -
  git branch -D ${NEWBRANCH}
fi
