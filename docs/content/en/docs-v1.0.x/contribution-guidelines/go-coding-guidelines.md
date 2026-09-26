---
title: "Go Coding Guidelines"
linkTitle: "Go Coding Guidelines"
weight: 2
description: >
  PipeCD-specific guidance for writing and reviewing Go code.
---

These guidelines aim to make everyday coding decisions explicit so
that code review can focus on whether a change does what it is
intended to do.

Apply the guidance to new and changed Go code. Do not mix unrelated
cleanups or repository-wide migrations into a functional change.
Existing exceptions do not need to be rewritten as part of adopting
this document.

General Go guidance is linked rather than restated. References apply
to the cited topics, not to an external guide in its entirety. The
text below records PipeCD-specific choices and exceptions; these are
not claims about the only correct way to write Go.

## Validate at boundaries

For general error handling, see [Go Code Review Comments: Don't Panic](https://go.dev/wiki/CodeReviewComments#dont-panic).

Concentrate validation at architectural boundaries, such as gRPC
request handling, configuration decoding, and datastore writes. Within
a boundary, rely on the guarantees established there rather than
repeating the same checks in every function.

Do not add routine nil checks to internal functions for arguments that
callers are required to supply. This does not remove the need to
handle nil values that are valid under an API's contract.

For PipeCD production code, do not use `panic` to report invalid
state; return an error instead. Do not introduce project-defined
`Must...` helpers in production code; such helpers may be used in
tests.

### Validate plugin responses at the SDK boundary

Document the meaning of a nil response and the conditions for a
successful response on each SDK interface method. When a plugin
returns successfully, the SDK should check whether a response is
present when the contract requires one.

Return `codes.Internal` for a nil response that the contract does not
explicitly allow. Do not treat that contract violation as success or
replace it with a default response. Keep this check at the SDK
boundary rather than duplicating it in internal conversion functions.

Preserve documented uses of nil. For example, `DetermineStrategy`
permits `(nil, nil)` to indicate that the plugin has no custom
strategy selection logic; the SDK then uses PipelineSync.

This response-presence check does not require general panic recovery
or comprehensive response-content validation.

### Express required capabilities in types

For the choice between interfaces and type parameters, see [The Go
Blog: When To Use Generics](https://go.dev/blog/when-generics).

For PipeCD configuration and plugin contracts:

- Require a configuration type's `Validate()` method through an
  interface or type constraint.
- Use runtime type assertions for optional capabilities, such as a
  plugin's implementation of `Initializer`. Do not make optional
  capabilities mandatory.
- Continue validating required fields and value ranges at boundaries.

Use this guidance when designing new or changed code; it does not call
for a repository-wide replacement of runtime type assertions.

## Error handling

For error messages, wrapping, and matching, see:

- [Go Code Review Comments: Error Strings](https://go.dev/wiki/CodeReviewComments#error-strings).
- [The Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors), including "Whether to
  Wrap".

For PipeCD, use sentinel errors and wrapping rather than introducing
custom error types. Use the standard library's error facilities rather
than `github.com/pkg/errors`.

Create gRPC status errors at the gRPC boundary. Stores, providers, and
plugin implementations should return ordinary Go errors; the boundary
layer is responsible for translating them into transport-level errors.

### Preserve diagnostic information from plugins

Plugin errors or stage logs should provide enough information to
identify the failed operation and investigate its cause. Record the
cause even when a normal execution failure is represented by a failed
stage result rather than an RPC error.

Keep conversion to gRPC status errors at the SDK boundary.
Plugin-specific sentinel errors do not need to retain their identity
across RPC calls.

Do not add a shared error-classification API or a new error wrapper
without a concrete diagnostic need or a caller that acts on an error
category. These guidelines do not prescribe `codes.Internal` for every
failure.

## Concurrency and context

For general concurrency guidance and API contracts, see:

- [context](https://pkg.go.dev/context) and [http.Request.Context](https://pkg.go.dev/net/http#Request.Context) for propagation and cancellation.
- [Go Code Review Comments: Goroutine Lifetimes](https://go.dev/wiki/CodeReviewComments#goroutine-lifetimes).
- [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup), including [Group.SetLimit](https://pkg.go.dev/golang.org/x/sync/errgroup#Group.SetLimit), for coordinated work and limits.
- [Effective Go: Channels](https://go.dev/doc/effective_go#channels) for buffering and synchronization.
- [Uber Go Style Guide: Zero-value Mutexes are Valid](https://github.com/uber-go/guide/blob/master/style.md#zero-value-mutexes-are-valid) for mutex fields.
- [time.NewTimer](https://pkg.go.dev/time#NewTimer) and [time.NewTicker](https://pkg.go.dev/time#NewTicker) for timer and ticker lifecycles.

As a PipeCD-specific exception to the general context guidance,
objects that explicitly own a lifecycle context may retain it. Make
that ownership clear; their lifecycle methods must cancel it.

In PipeCD production code, reserve `context.Background()` for
process-root contexts and work deliberately detached from parent
cancellation for graceful shutdown. Explain the latter case in a
comment.

Use `go.uber.org/atomic` rather than `sync/atomic`.

Name mutex fields `mu` by default; use purpose-specific prefixes when
a struct has more than one mutex.

### Preserve fan-out failure semantics

Choose synchronization tools according to the required behavior, not
just for visual consistency. Decide whether one failure cancels the
other operations, whether all operations must finish, and how their
results determine the overall result. Reuse an existing helper when it
provides those same semantics.

For collecting one result per input, see [errgroup's Parallel example](https://pkg.go.dev/golang.org/x/sync/errgroup#example-Group-Parallel).

Choose `errgroup.WithContext` or `sync.WaitGroup` to suit the
behavior. Do not replace one with the other indiscriminately. Preserve
rollback's existing partial-success semantics when changing
synchronization tools. A new generic fan-out helper is not required.

## Logging

For structured logging with zap, see the [zap package overview](https://pkg.go.dev/go.uber.org/zap#pkg-overview).

For component-owned loggers, accept a zap logger through the
constructor and store it in an unexported field. Do not pass loggers
through context.

Use the per-call logger supplied by an executor or SDK input when that
is the interface's existing model; plugin authors do not need to keep
another logger solely to follow the constructor pattern.

Use `logger.Named("controller")` to name a component, rather than
adding a `zap.String("component", ...)` field.

Use kebab-case for structured log field names and the same key for the
same concept. Use `app-id` for an application ID.

Use zap's structured fields for variable values.

Use the standard `log` package only in `main` before zap has been
initialized.

Logging details and returning an error are not inherently duplicate
reporting. At a gRPC boundary, logging detailed diagnostics while
returning a more general client-facing error can serve separate
responsibilities.

## Interfaces and dependency injection

For interface design and ownership, see [Go Code Review Comments:
Interfaces](https://go.dev/wiki/CodeReviewComments#interfaces).

For PipeCD constructors, use `New` followed by the type name, rather
than `Make`, `Create`, or `Build` as constructor prefixes.

Keep functional options in infrastructure packages, such as file
storage, datastore backends, RPC, Git, and Redis. Use ordinary
constructor parameters for domain logic. Do not use an option to
inject a logger.

### Choose test doubles by maintenance cost

Default to a concise handwritten test double for a small interface
needed by the code under test. Generated mocks are also acceptable
when an existing mock fits the test or a handwritten implementation
would be costly to maintain.

Do not mandate one approach per module. Do not introduce dependencies
that cross the established module boundaries just to reuse test
doubles.

Follow local naming conventions. Do not require every handwritten
double to be named a "fake" or reserve "mock" for generated code. A
double's role in the test matters more than how it was produced.

No bulk migration of existing tests is required.

## Tests

For general test structure and helper usage, see:

- [Go Test Comments: Table-Driven Tests vs Multiple Test Functions](https://go.dev/wiki/TestComments#table-driven-tests-vs-multiple-test-functions).
- [Go Test Comments: Mark Test Helpers](https://go.dev/wiki/TestComments#mark-test-helpers).

For PipeCD, write tests in the same package as the code under test. In
table-driven tests, name the case-name field `name`.

Use `t.Parallel()` in both top-level tests and subtests when they can
run independently. Do not parallelize tests that mutate process-wide
state or depend on resources that are not isolated. See [testing:
Subtests and Sub-benchmarks](https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks) for parallel execution semantics.

For tests that require an external environment configured through
environment variables, use `t.Skip` when that environment is not
configured.

### Use consistent names in new test tables

In new table-driven tests, use `want` for an expected value and
`wantErr` for an expected error. When there are multiple expected
values, add meaning to the name, such as `wantStatus`.

When adding cases to an existing table, follow that table's naming.
This is a default for new code, not a claim that other names are
wrong. Do not add unrelated renames, a bulk naming migration, or a
dedicated linter.

### Choose assertions by whether the test can continue

For failure continuation, see [Go Test Comments: Keep Going](https://go.dev/wiki/TestComments#keep-going). For
assertion behavior and goroutine restrictions, see [testify: assert](https://github.com/stretchr/testify#assert-package),
[testify: require](https://github.com/stretchr/testify#require-package), and [testing.T.FailNow](https://pkg.go.dev/testing#T.FailNow).

Within PipeCD's existing testify tests, use `require` for
prerequisites and `assert` for independent expectations. This is a
project-specific choice, not adoption of Go Test Comments'
recommendation to avoid assertion libraries.

Not every error assertion requires `require`; choose according to
whether subsequent checks depend on it. Do not enforce this choice
indiscriminately through `testifylint`'s `require-error` rule.

## Module boundaries

Independent plugin modules under `pkg/app/pipedv1/plugin/` must not
import the PipeCD main module, `github.com/pipe-cd/pipecd`. Use
`github.com/pipe-cd/piped-plugin-sdk-go` for integration with PipeCD.
This does not prohibit dependencies on the standard library or
third-party libraries.

In pipedv1 code, use `pkg/configv1`, not the legacy `pkg/config`.

Do not commit `replace` directives in Go module files.

For the deprecated `io/ioutil` package and its replacements, see the
[io/ioutil documentation](https://pkg.go.dev/io/ioutil).

## Naming

For general Go naming, see:

- [Go Code Review Comments: Initialisms](https://go.dev/wiki/CodeReviewComments#initialisms) and [Package Names](https://go.dev/wiki/CodeReviewComments#package-names).
- [Effective Go: Getters](https://go.dev/doc/effective_go#Getters).

Use snake_case for file names.

For new JSON field names, use forms such as `projectID` and `repoURL`.
Do not rename an existing serialized field solely for consistency;
compatibility changes require separate consideration.

Accessors that follow protobuf-generated naming are an exception to
the getter naming guidance. Operations that retrieve data are not
ordinary field accessors and may use `Get` in their names.

## Commits and pull requests

Follow [CONTRIBUTING.md](https://github.com/pipe-cd/pipecd/blob/master/CONTRIBUTING.md) for commit messages, pull requests, DCO
sign-off, licensing, and contribution checks.

## Automated checks

Use the repository's existing formatting, linting, and test commands.
See [.golangci.yml](https://github.com/pipe-cd/pipecd/blob/master/.golangci.yml) for the current linter and formatter configuration.

This document does not imply that every guideline is already enforced
by a linter. New checks, fixes to existing checks, and migrations of
existing code should be proposed separately. Any added check should
enforce the intended behavior rather than a simpler but different
rule.
