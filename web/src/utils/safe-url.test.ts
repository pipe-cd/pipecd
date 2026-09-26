import { getSafeExternalUrl, getSafeRedirectPath } from "./safe-url";

describe("getSafeRedirectPath", () => {
  const fallback = "/applications";

  it.each([
    ["/deployments", "/deployments"],
    [
      "/deployments/abc?project=quickstart",
      "/deployments/abc?project=quickstart",
    ],
    ["/", "/"],
  ])("given %p, returns %p", (path, expected) => {
    expect(getSafeRedirectPath(path, fallback)).toBe(expected);
  });

  it.each([
    [null],
    [undefined],
    [""],
    ["deployments"],
    ["//evil.example.com"],
    ["/\\evil.example.com"],
    ["\\\\evil.example.com"],
    ["/deployments\\..\\"],
    ["https://evil.example.com"],
    ["javascript:alert(1)"],
  ])("given unsafe path %p, returns the fallback", (path) => {
    expect(getSafeRedirectPath(path, fallback)).toBe(fallback);
  });
});

describe("getSafeExternalUrl", () => {
  it.each([
    ["https://github.com/pipe-cd/pipecd/commit/abc"],
    ["http://git.example.com/commit/abc"],
  ])("given %p, returns it", (url) => {
    expect(getSafeExternalUrl(url)).toBe(url);
  });

  it.each([
    [null],
    [undefined],
    [""],
    ["commit-url"],
    ["/relative/path"],
    ["//evil.example.com"],
    ["javascript:alert(1)"],
    ["JavaScript:alert(1)"],
    ["data:text/html,<script>alert(1)</script>"],
  ])("given unsafe url %p, returns undefined", (url) => {
    expect(getSafeExternalUrl(url)).toBeUndefined();
  });
});
