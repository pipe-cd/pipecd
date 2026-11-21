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

// Run with: node --test hack/release-contributors.test.js

const assert = require('node:assert');
const { test } = require('node:test');

const {
  parseTag,
  findPreviousTag,
  pullRequestNumber,
  cherryPickedNumbers,
  renderSection,
  mergeSection,
} = require('./release-contributors.js');

test('parseTag', () => {
  assert.deepStrictEqual(parseTag('v0.58.0'), { tag: 'v0.58.0', path: '', version: [0, 58, 0], prerelease: '' });
  assert.strictEqual(parseTag('pkg/app/pipedv1/plugin/kubernetes/v0.4.0').path, 'pkg/app/pipedv1/plugin/kubernetes');
  assert.strictEqual(parseTag('v1.0.0-rc1').prerelease, '-rc1');
  assert.strictEqual(parseTag('0.58.0'), null);
  assert.strictEqual(parseTag('v0.58'), null);
  assert.strictEqual(parseTag('--upload-pack=x/v1.0.0'), null);
});

test('findPreviousTag', () => {
  const tags = [
    'v0.9.0',
    'v0.10.0',
    'v0.10.1',
    'v0.11.0-rc1',
    'v0.11.0',
    'pkg/app/pipedv1/plugin/kubernetes/v0.1.0',
    'pkg/app/pipedv1/plugin/kubernetes/v0.2.0',
    'pkg/app/pipedv1/plugin/kubernetes_multicluster/v0.1.0',
    'not-a-release',
  ];
  const prev = (tag) => {
    const p = findPreviousTag(parseTag(tag), tags);
    return p ? p.tag : null;
  };
  assert.strictEqual(prev('v0.10.0'), 'v0.9.0');
  assert.strictEqual(prev('v0.11.0'), 'v0.10.1');
  assert.strictEqual(prev('v0.11.0-rc1'), 'v0.10.1');
  assert.strictEqual(prev('v0.11.0-rc2'), 'v0.11.0-rc1');
  assert.strictEqual(prev('v0.12.0'), 'v0.11.0');
  assert.strictEqual(prev('v0.9.0'), null);
  assert.strictEqual(prev('pkg/app/pipedv1/plugin/kubernetes/v0.2.0'), 'pkg/app/pipedv1/plugin/kubernetes/v0.1.0');
  assert.strictEqual(prev('pkg/app/pipedv1/plugin/kubernetes_multicluster/v0.1.0'), null);
});

test('pullRequestNumber', () => {
  assert.strictEqual(pullRequestNumber('Fix git worktree race condition issue (#7201)'), 7201);
  assert.strictEqual(pullRequestNumber('Revert "Fix (#10)" (#12)'), 12);
  assert.strictEqual(pullRequestNumber('Merge pull request #38 from someone/branch'), 38);
  assert.strictEqual(pullRequestNumber('Fix #12 without a pull request'), null);
});

test('cherryPickedNumbers', () => {
  const body = '**What this PR does / why we need it**:\nCherry pick of #6240 #6242 #6267.\n\n**Note:** see #1';
  assert.deepStrictEqual(cherryPickedNumbers(body), [6240, 6242, 6267]);
  assert.deepStrictEqual(cherryPickedNumbers('Fixes #12'), []);
  assert.deepStrictEqual(cherryPickedNumbers(null), []);
});

test('renderSection and mergeSection', () => {
  const contributors = [{ login: 'alice', number: 3 }, { login: 'bob', number: 5 }];
  const section = renderSection('pipe-cd', 'pipecd', contributors, [contributors[1]]);
  assert.match(section, /@alice, @bob/);
  assert.match(section, /@bob made their first contribution in https:\/\/github.com\/pipe-cd\/pipecd\/pull\/5/);
  assert.doesNotMatch(renderSection('pipe-cd', 'pipecd', contributors, []), /New Contributors/);

  const merged = mergeSection('## Release v1.0.0\r\n\r\n- A change\r\n', section);
  assert.strictEqual(merged, `## Release v1.0.0\n\n- A change\n\n${section}\n`);

  // Running again replaces the section instead of adding a second one.
  const updated = renderSection('pipe-cd', 'pipecd', [contributors[0]], []);
  assert.strictEqual(mergeSection(merged, updated), `## Release v1.0.0\n\n- A change\n\n${updated}\n`);
  assert.strictEqual(mergeSection(merged, section), merged);
  assert.strictEqual(mergeSection('', section), `${section}\n`);
});
