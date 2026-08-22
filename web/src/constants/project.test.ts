import { POLICIES_STRING_REGEX } from "./project";

describe("POLICIES_STRING_REGEX", () => {
  it.each([
    "resources=*;actions=*",
    "resources=application;actions=get",
    "resources=application,piped;actions=get,list",
    "resources=application;actions=get,list,create,update,delete",
  ])("matches policy without labels: %s", (policy) => {
    expect(POLICIES_STRING_REGEX.test(policy)).toBe(true);
  });

  // Regression test for https://github.com/pipe-cd/pipecd/issues/7172
  // Label scoped resources (resources=NAME{key:value}) were rejected by the
  // form validation even though parseRBACPolicies/formalizePoliciesList
  // already supported them.
  it.each([
    "resources=application{env:prod};actions=get",
    "resources=application{env:prod,team:foo};actions=get",
    "resources=application{env:prod},piped;actions=get",
    "resources=application{env:prod},piped{env:prod,team:foo};actions=get,list",
  ])("matches policy with label scoped resources: %s", (policy) => {
    expect(POLICIES_STRING_REGEX.test(policy)).toBe(true);
  });

  it.each([
    "resources=application{env:prod;actions=get", // unclosed brace
    "resources=bogus;actions=get", // unknown resource type
    "resources=application;actions=bogus", // unknown action
    "resources=application", // missing actions
    "resources=application{a=b;c:d};actions=get", // semicolon escapes the label block
    "resources=application{};actions=get", // empty label block
  ])("does not match invalid policy: %s", (policy) => {
    expect(POLICIES_STRING_REGEX.test(policy)).toBe(false);
  });
});
