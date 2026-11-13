// Copyright 2026 The PipeCD Authors.
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

// Credits the contributors of a release. It is run by the
// release-contributors workflow through actions/github-script and does two
// things for the given release tag:
//
//   1. Appends a "Contributors" section to the release notes, listing every
//      pull request author in the release and calling out first-time ones.
//   2. Comments on every pull request in the release to say which release
//      it shipped in.
//
// Both steps are idempotent, so the workflow can safely be run again.
//
// It works for PipeCD releases (vX.Y.Z) and for plugin releases
// (<plugin path>/vX.Y.Z). The full git history and all tags must be available
// in the working directory.

const { execFileSync } = require('child_process');

const SECTION_START = '<!-- release-contributors:start -->';
const SECTION_END = '<!-- release-contributors:end -->';
const TAG_PATTERN = /^(?:([A-Za-z0-9_][A-Za-z0-9_.\/-]*)\/)?v(\d+)\.(\d+)\.(\d+)(-[0-9A-Za-z.-]+)?$/;
const BATCH_SIZE = 25;

// parseTag splits a release tag into the plugin path (empty for PipeCD
// releases) and its version parts. It returns null for any other tag.
function parseTag(tag) {
  const m = TAG_PATTERN.exec(tag);
  if (!m) {
    return null;
  }
  return {
    tag,
    path: m[1] || '',
    version: [Number(m[2]), Number(m[3]), Number(m[4])],
    prerelease: m[5] || '',
  };
}

// compareTags orders two parsed tags by semantic version.
// A pre-release sorts before the release it leads to.
function compareTags(a, b) {
  for (let i = 0; i < 3; i++) {
    if (a.version[i] !== b.version[i]) {
      return a.version[i] - b.version[i];
    }
  }
  if (a.prerelease === b.prerelease) {
    return 0;
  }
  if (a.prerelease === '') {
    return 1;
  }
  if (b.prerelease === '') {
    return -1;
  }
  return a.prerelease.localeCompare(b.prerelease, 'en', { numeric: true });
}

// findPreviousTag returns the newest tag of the same component that is older
// than the current one, or null when there is none. Pre-releases are skipped
// unless the current tag is a pre-release too, so that a final release covers
// everything since the last final release.
function findPreviousTag(current, allTags) {
  const candidates = allTags
    .map(parseTag)
    .filter((t) => t !== null)
    .filter((t) => t.path === current.path)
    .filter((t) => current.prerelease !== '' || t.prerelease === '')
    .filter((t) => compareTags(t, current) < 0)
    .sort(compareTags);
  return candidates.length > 0 ? candidates[candidates.length - 1] : null;
}

// pullRequestNumber extracts the pull request number from the subject of a
// squash commit ("... (#123)") or a merge commit ("Merge pull request #123").
function pullRequestNumber(subject) {
  let m = /\(#(\d+)\)\s*$/.exec(subject);
  if (m) {
    return Number(m[1]);
  }
  m = /^Merge pull request #(\d+) /.exec(subject);
  return m ? Number(m[1]) : null;
}

// cherryPickedNumbers returns the original pull requests listed in the body of
// a pull request opened by the cherry_pick workflow for a patch release.
function cherryPickedNumbers(body) {
  const m = /Cherry pick of ((?:#\d+[ ,]*)+)/.exec(body || '');
  return m ? m[1].match(/\d+/g).map(Number) : [];
}

function isBot(author) {
  return author.__typename === 'Bot' || author.login.endsWith('[bot]');
}

// renderSection builds the block that is added to the release notes.
// contributors and newcomers are lists of { login, number }.
function renderSection(owner, repo, contributors, newcomers) {
  const lines = [SECTION_START, '### Contributors', ''];
  lines.push('Thank you to everyone who contributed to this release:', '');
  lines.push(contributors.map((c) => `@${c.login}`).join(', '));
  if (newcomers.length > 0) {
    lines.push('', '#### New Contributors', '');
    for (const c of newcomers) {
      lines.push(`  - @${c.login} made their first contribution in https://github.com/${owner}/${repo}/pull/${c.number}`);
    }
  }
  lines.push(SECTION_END);
  return lines.join('\n');
}

// mergeSection replaces the existing contributors block of a release body,
// or appends the block when there is none yet.
function mergeSection(body, section) {
  const text = (body || '').replace(/\r\n/g, '\n');
  const start = text.indexOf(SECTION_START);
  const end = text.indexOf(SECTION_END);
  if (start !== -1 && end > start) {
    return text.slice(0, start) + section + text.slice(end + SECTION_END.length);
  }
  const head = text.replace(/\s+$/, '');
  return head === '' ? section + '\n' : `${head}\n\n${section}\n`;
}

function commentMarker(tag) {
  return `<!-- released-in: ${tag} -->`;
}

function git(args) {
  return execFileSync('git', args, { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 }).trim();
}

function chunk(list, size) {
  const out = [];
  for (let i = 0; i < list.length; i += size) {
    out.push(list.slice(i, i + size));
  }
  return out;
}

// fetchPullRequests loads the author and the latest comments of the given
// pull requests, a batch per request to stay well inside the rate limit.
async function fetchPullRequests(github, owner, repo, numbers) {
  const out = new Map();
  for (const batch of chunk(numbers, BATCH_SIZE)) {
    const fields = batch.map((n) => `
      pr${n}: pullRequest(number: ${n}) {
        number
        merged
        body
        author { __typename login }
        comments(last: 50) { nodes { body } }
      }`);
    const query = `query($owner: String!, $repo: String!) {
      repository(owner: $owner, name: $repo) {${fields.join('')}
      }
    }`;
    let data;
    try {
      data = await github.graphql(query, { owner, repo });
    } catch (err) {
      // A number that is not a pull request fails only its own field.
      if (!err.data) {
        throw err;
      }
      data = err.data;
    }
    for (const pr of Object.values((data && data.repository) || {})) {
      if (pr) {
        out.set(pr.number, pr);
      }
    }
  }
  return out;
}

// findNewcomers returns the logins that had no pull request merged before
// the given time, which makes this release their first one.
async function findNewcomers(github, owner, repo, logins, before) {
  const out = new Set();
  for (const batch of chunk(logins, BATCH_SIZE)) {
    const vars = {};
    const decls = [];
    const fields = [];
    batch.forEach((login, i) => {
      vars[`q${i}`] = `repo:${owner}/${repo} is:pr is:merged author:${login} merged:<${before}`;
      decls.push(`$q${i}: String!`);
      fields.push(`u${i}: search(query: $q${i}, type: ISSUE, first: 1) { issueCount }`);
    });
    const data = await github.graphql(`query(${decls.join(', ')}) { ${fields.join('\n')} }`, vars);
    batch.forEach((login, i) => {
      if (data[`u${i}`].issueCount === 0) {
        out.add(login);
      }
    });
  }
  return out;
}

async function run({ github, context, core, tag, dryRun }) {
  const { owner, repo } = context.repo;

  const current = parseTag(tag || '');
  if (!current) {
    core.setFailed(`"${tag}" is not a release tag (expected vX.Y.Z or <plugin path>/vX.Y.Z)`);
    return;
  }

  const previous = findPreviousTag(current, git(['tag', '--list']).split('\n'));
  core.info(`Release ${current.tag}, previous release ${previous ? previous.tag : '(none)'}`);
  if (!previous && current.path === '') {
    core.setFailed(`Could not find the release before ${current.tag}`);
    return;
  }

  // Collect the pull requests of the release, oldest first.
  const logArgs = ['log', '--reverse', '--format=%s', previous ? `${previous.tag}..${current.tag}` : current.tag];
  if (current.path !== '') {
    logArgs.push('--', current.path);
  }
  const numbers = [];
  for (const subject of git(logArgs).split('\n')) {
    const n = pullRequestNumber(subject);
    if (n !== null && !numbers.includes(n)) {
      numbers.push(n);
    }
  }
  core.info(`Found ${numbers.length} pull requests in the release`);
  if (numbers.length === 0) {
    return;
  }

  const pulls = await fetchPullRequests(github, owner, repo, numbers);

  // A patch release is made of cherry-pick pull requests opened by a bot.
  // Credit the original pull requests they carry instead.
  const expanded = [];
  for (const n of numbers) {
    const pr = pulls.get(n);
    const originals = pr && pr.author && isBot(pr.author) ? cherryPickedNumbers(pr.body) : [];
    for (const o of originals.length > 0 ? originals : [n]) {
      if (!expanded.includes(o)) {
        expanded.push(o);
      }
    }
  }
  const missing = expanded.filter((n) => !pulls.has(n));
  for (const [n, pr] of await fetchPullRequests(github, owner, repo, missing)) {
    pulls.set(n, pr);
  }

  const merged = expanded
    .map((n) => pulls.get(n))
    .filter((pr) => pr && pr.merged && pr.author && !isBot(pr.author));

  // The first pull request of each author in this release.
  const firstPull = new Map();
  for (const pr of merged) {
    if (!firstPull.has(pr.author.login)) {
      firstPull.set(pr.author.login, pr.number);
    }
  }
  const contributors = [...firstPull.entries()]
    .map(([login, number]) => ({ login, number }))
    .sort((a, b) => a.login.localeCompare(b.login, 'en', { sensitivity: 'base' }));
  if (contributors.length === 0) {
    core.info('No contributors to credit');
    return;
  }

  // Without a previous release there is no line to measure "first" against.
  let newcomers = [];
  if (previous) {
    const before = new Date(git(['log', '-1', '--format=%cI', previous.tag])).toISOString().replace(/\.\d+Z$/, 'Z');
    const logins = await findNewcomers(github, owner, repo, contributors.map((c) => c.login), before);
    newcomers = contributors.filter((c) => logins.has(c.login));
  }
  core.info(`${contributors.length} contributors, ${newcomers.length} of them new`);

  const { data: release } = await github.request('GET /repos/{owner}/{repo}/releases/tags/{tag}', { owner, repo, tag: current.tag });

  // 1. Add the contributors section to the release notes.
  const section = renderSection(owner, repo, contributors, newcomers);
  const body = mergeSection(release.body, section);
  if (body === (release.body || '').replace(/\r\n/g, '\n')) {
    core.info('Release notes already list the contributors');
  } else if (dryRun) {
    core.info(`[dry run] Would add this section to the release notes:\n${section}`);
  } else {
    await github.request('PATCH /repos/{owner}/{repo}/releases/{release_id}', { owner, repo, release_id: release.id, body });
    core.info('Added the contributors section to the release notes');
  }

  // 2. Tell every pull request which release it shipped in.
  // Pre-releases are skipped so a change is announced once, by the final release.
  if (release.prerelease || current.prerelease !== '') {
    core.info('Pre-release: not commenting on pull requests');
    return;
  }
  const marker = commentMarker(current.tag);
  const pending = merged.filter((pr) => !pr.comments.nodes.some((c) => c.body.includes(marker)));
  core.info(`Commenting on ${pending.length} pull requests (${merged.length - pending.length} already done)`);
  let failed = 0;
  for (const pr of pending) {
    const message = `${marker}\nThis pull request was released in [${current.tag}](${release.html_url}). Thank you for the contribution, @${pr.author.login}!`;
    if (dryRun) {
      core.info(`[dry run] Would comment on #${pr.number}: ${message.split('\n')[1]}`);
      continue;
    }
    try {
      await github.request('POST /repos/{owner}/{repo}/issues/{issue_number}/comments', { owner, repo, issue_number: pr.number, body: message });
    } catch (err) {
      failed++;
      core.warning(`Failed to comment on #${pr.number}: ${err.message}`);
    }
    // Stay below the secondary rate limit for creating content.
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  if (failed > 0) {
    core.setFailed(`Failed to comment on ${failed} pull requests; run the workflow again for ${current.tag} to retry them`);
  }
}

module.exports = run;
Object.assign(module.exports, {
  parseTag,
  compareTags,
  findPreviousTag,
  pullRequestNumber,
  cherryPickedNumbers,
  renderSection,
  mergeSection,
});
