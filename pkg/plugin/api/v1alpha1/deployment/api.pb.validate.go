// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/plugin/api/v1alpha1/deployment/api.proto

package deployment

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"

	model "github.com/pipe-cd/pipecd/pkg/model"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort

	_ = model.SyncStrategy(0)
)

// Validate checks the field values on DetermineVersionsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DetermineVersionsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DetermineVersionsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DetermineVersionsRequestMultiError, or nil if none found.
func (m *DetermineVersionsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DetermineVersionsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetInput() == nil {
		err := DetermineVersionsRequestValidationError{
			field:  "Input",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetInput()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DetermineVersionsRequestValidationError{
					field:  "Input",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DetermineVersionsRequestValidationError{
					field:  "Input",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInput()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DetermineVersionsRequestValidationError{
				field:  "Input",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return DetermineVersionsRequestMultiError(errors)
	}

	return nil
}

// DetermineVersionsRequestMultiError is an error wrapping multiple validation
// errors returned by DetermineVersionsRequest.ValidateAll() if the designated
// constraints aren't met.
type DetermineVersionsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DetermineVersionsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DetermineVersionsRequestMultiError) AllErrors() []error { return m }

// DetermineVersionsRequestValidationError is the validation error returned by
// DetermineVersionsRequest.Validate if the designated constraints aren't met.
type DetermineVersionsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DetermineVersionsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DetermineVersionsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DetermineVersionsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DetermineVersionsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DetermineVersionsRequestValidationError) ErrorName() string {
	return "DetermineVersionsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DetermineVersionsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDetermineVersionsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DetermineVersionsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DetermineVersionsRequestValidationError{}

// Validate checks the field values on DetermineVersionsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DetermineVersionsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DetermineVersionsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DetermineVersionsResponseMultiError, or nil if none found.
func (m *DetermineVersionsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DetermineVersionsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetVersions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DetermineVersionsResponseValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DetermineVersionsResponseValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DetermineVersionsResponseValidationError{
					field:  fmt.Sprintf("Versions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return DetermineVersionsResponseMultiError(errors)
	}

	return nil
}

// DetermineVersionsResponseMultiError is an error wrapping multiple validation
// errors returned by DetermineVersionsResponse.ValidateAll() if the
// designated constraints aren't met.
type DetermineVersionsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DetermineVersionsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DetermineVersionsResponseMultiError) AllErrors() []error { return m }

// DetermineVersionsResponseValidationError is the validation error returned by
// DetermineVersionsResponse.Validate if the designated constraints aren't met.
type DetermineVersionsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DetermineVersionsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DetermineVersionsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DetermineVersionsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DetermineVersionsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DetermineVersionsResponseValidationError) ErrorName() string {
	return "DetermineVersionsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DetermineVersionsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDetermineVersionsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DetermineVersionsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DetermineVersionsResponseValidationError{}

// Validate checks the field values on DetermineStrategyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DetermineStrategyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DetermineStrategyRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DetermineStrategyRequestMultiError, or nil if none found.
func (m *DetermineStrategyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DetermineStrategyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetInput() == nil {
		err := DetermineStrategyRequestValidationError{
			field:  "Input",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetInput()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DetermineStrategyRequestValidationError{
					field:  "Input",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DetermineStrategyRequestValidationError{
					field:  "Input",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetInput()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DetermineStrategyRequestValidationError{
				field:  "Input",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return DetermineStrategyRequestMultiError(errors)
	}

	return nil
}

// DetermineStrategyRequestMultiError is an error wrapping multiple validation
// errors returned by DetermineStrategyRequest.ValidateAll() if the designated
// constraints aren't met.
type DetermineStrategyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DetermineStrategyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DetermineStrategyRequestMultiError) AllErrors() []error { return m }

// DetermineStrategyRequestValidationError is the validation error returned by
// DetermineStrategyRequest.Validate if the designated constraints aren't met.
type DetermineStrategyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DetermineStrategyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DetermineStrategyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DetermineStrategyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DetermineStrategyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DetermineStrategyRequestValidationError) ErrorName() string {
	return "DetermineStrategyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DetermineStrategyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDetermineStrategyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DetermineStrategyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DetermineStrategyRequestValidationError{}

// Validate checks the field values on DetermineStrategyResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DetermineStrategyResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DetermineStrategyResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DetermineStrategyResponseMultiError, or nil if none found.
func (m *DetermineStrategyResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DetermineStrategyResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for SyncStrategy

	// no validation rules for Summary

	if len(errors) > 0 {
		return DetermineStrategyResponseMultiError(errors)
	}

	return nil
}

// DetermineStrategyResponseMultiError is an error wrapping multiple validation
// errors returned by DetermineStrategyResponse.ValidateAll() if the
// designated constraints aren't met.
type DetermineStrategyResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DetermineStrategyResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DetermineStrategyResponseMultiError) AllErrors() []error { return m }

// DetermineStrategyResponseValidationError is the validation error returned by
// DetermineStrategyResponse.Validate if the designated constraints aren't met.
type DetermineStrategyResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DetermineStrategyResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DetermineStrategyResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DetermineStrategyResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DetermineStrategyResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DetermineStrategyResponseValidationError) ErrorName() string {
	return "DetermineStrategyResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DetermineStrategyResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDetermineStrategyResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DetermineStrategyResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DetermineStrategyResponseValidationError{}

// Validate checks the field values on BuildStagesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *BuildStagesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BuildStagesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// BuildStagesRequestMultiError, or nil if none found.
func (m *BuildStagesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *BuildStagesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetStages() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, BuildStagesRequestValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, BuildStagesRequestValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return BuildStagesRequestValidationError{
					field:  fmt.Sprintf("Stages[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return BuildStagesRequestMultiError(errors)
	}

	return nil
}

// BuildStagesRequestMultiError is an error wrapping multiple validation errors
// returned by BuildStagesRequest.ValidateAll() if the designated constraints
// aren't met.
type BuildStagesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BuildStagesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BuildStagesRequestMultiError) AllErrors() []error { return m }

// BuildStagesRequestValidationError is the validation error returned by
// BuildStagesRequest.Validate if the designated constraints aren't met.
type BuildStagesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BuildStagesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BuildStagesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BuildStagesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BuildStagesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BuildStagesRequestValidationError) ErrorName() string {
	return "BuildStagesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e BuildStagesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBuildStagesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BuildStagesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BuildStagesRequestValidationError{}

// Validate checks the field values on BuildStagesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *BuildStagesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BuildStagesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// BuildStagesResponseMultiError, or nil if none found.
func (m *BuildStagesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *BuildStagesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetStages() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, BuildStagesResponseValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, BuildStagesResponseValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return BuildStagesResponseValidationError{
					field:  fmt.Sprintf("Stages[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return BuildStagesResponseMultiError(errors)
	}

	return nil
}

// BuildStagesResponseMultiError is an error wrapping multiple validation
// errors returned by BuildStagesResponse.ValidateAll() if the designated
// constraints aren't met.
type BuildStagesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BuildStagesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BuildStagesResponseMultiError) AllErrors() []error { return m }

// BuildStagesResponseValidationError is the validation error returned by
// BuildStagesResponse.Validate if the designated constraints aren't met.
type BuildStagesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BuildStagesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BuildStagesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BuildStagesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BuildStagesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BuildStagesResponseValidationError) ErrorName() string {
	return "BuildStagesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e BuildStagesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBuildStagesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BuildStagesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BuildStagesResponseValidationError{}

// Validate checks the field values on FetchDefinedStagesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FetchDefinedStagesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FetchDefinedStagesRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FetchDefinedStagesRequestMultiError, or nil if none found.
func (m *FetchDefinedStagesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *FetchDefinedStagesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return FetchDefinedStagesRequestMultiError(errors)
	}

	return nil
}

// FetchDefinedStagesRequestMultiError is an error wrapping multiple validation
// errors returned by FetchDefinedStagesRequest.ValidateAll() if the
// designated constraints aren't met.
type FetchDefinedStagesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FetchDefinedStagesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FetchDefinedStagesRequestMultiError) AllErrors() []error { return m }

// FetchDefinedStagesRequestValidationError is the validation error returned by
// FetchDefinedStagesRequest.Validate if the designated constraints aren't met.
type FetchDefinedStagesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FetchDefinedStagesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FetchDefinedStagesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FetchDefinedStagesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FetchDefinedStagesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FetchDefinedStagesRequestValidationError) ErrorName() string {
	return "FetchDefinedStagesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e FetchDefinedStagesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFetchDefinedStagesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FetchDefinedStagesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FetchDefinedStagesRequestValidationError{}

// Validate checks the field values on FetchDefinedStagesResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *FetchDefinedStagesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FetchDefinedStagesResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FetchDefinedStagesResponseMultiError, or nil if none found.
func (m *FetchDefinedStagesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *FetchDefinedStagesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return FetchDefinedStagesResponseMultiError(errors)
	}

	return nil
}

// FetchDefinedStagesResponseMultiError is an error wrapping multiple
// validation errors returned by FetchDefinedStagesResponse.ValidateAll() if
// the designated constraints aren't met.
type FetchDefinedStagesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FetchDefinedStagesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FetchDefinedStagesResponseMultiError) AllErrors() []error { return m }

// FetchDefinedStagesResponseValidationError is the validation error returned
// by FetchDefinedStagesResponse.Validate if the designated constraints aren't met.
type FetchDefinedStagesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FetchDefinedStagesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FetchDefinedStagesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FetchDefinedStagesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FetchDefinedStagesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FetchDefinedStagesResponseValidationError) ErrorName() string {
	return "FetchDefinedStagesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e FetchDefinedStagesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFetchDefinedStagesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FetchDefinedStagesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FetchDefinedStagesResponseValidationError{}

// Validate checks the field values on PlanPluginInput with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PlanPluginInput) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PlanPluginInput with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PlanPluginInputMultiError, or nil if none found.
func (m *PlanPluginInput) ValidateAll() error {
	return m.validate(true)
}

func (m *PlanPluginInput) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetDeployment() == nil {
		err := PlanPluginInputValidationError{
			field:  "Deployment",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetDeployment()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "Deployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "Deployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDeployment()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PlanPluginInputValidationError{
				field:  "Deployment",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for PluginConfig

	if all {
		switch v := interface{}(m.GetRunningDeploymentSource()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "RunningDeploymentSource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "RunningDeploymentSource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRunningDeploymentSource()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PlanPluginInputValidationError{
				field:  "RunningDeploymentSource",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetTargetDeploymentSource()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "TargetDeploymentSource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PlanPluginInputValidationError{
					field:  "TargetDeploymentSource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTargetDeploymentSource()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PlanPluginInputValidationError{
				field:  "TargetDeploymentSource",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PlanPluginInputMultiError(errors)
	}

	return nil
}

// PlanPluginInputMultiError is an error wrapping multiple validation errors
// returned by PlanPluginInput.ValidateAll() if the designated constraints
// aren't met.
type PlanPluginInputMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PlanPluginInputMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PlanPluginInputMultiError) AllErrors() []error { return m }

// PlanPluginInputValidationError is the validation error returned by
// PlanPluginInput.Validate if the designated constraints aren't met.
type PlanPluginInputValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PlanPluginInputValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PlanPluginInputValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PlanPluginInputValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PlanPluginInputValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PlanPluginInputValidationError) ErrorName() string { return "PlanPluginInputValidationError" }

// Error satisfies the builtin error interface
func (e PlanPluginInputValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPlanPluginInput.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PlanPluginInputValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PlanPluginInputValidationError{}

// Validate checks the field values on BuildStagesRequest_StageConfig with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *BuildStagesRequest_StageConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on BuildStagesRequest_StageConfig with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// BuildStagesRequest_StageConfigMultiError, or nil if none found.
func (m *BuildStagesRequest_StageConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *BuildStagesRequest_StageConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := BuildStagesRequest_StageConfigValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Desc

	// no validation rules for Timeout

	// no validation rules for Config

	if len(errors) > 0 {
		return BuildStagesRequest_StageConfigMultiError(errors)
	}

	return nil
}

// BuildStagesRequest_StageConfigMultiError is an error wrapping multiple
// validation errors returned by BuildStagesRequest_StageConfig.ValidateAll()
// if the designated constraints aren't met.
type BuildStagesRequest_StageConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m BuildStagesRequest_StageConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m BuildStagesRequest_StageConfigMultiError) AllErrors() []error { return m }

// BuildStagesRequest_StageConfigValidationError is the validation error
// returned by BuildStagesRequest_StageConfig.Validate if the designated
// constraints aren't met.
type BuildStagesRequest_StageConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e BuildStagesRequest_StageConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e BuildStagesRequest_StageConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e BuildStagesRequest_StageConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e BuildStagesRequest_StageConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e BuildStagesRequest_StageConfigValidationError) ErrorName() string {
	return "BuildStagesRequest_StageConfigValidationError"
}

// Error satisfies the builtin error interface
func (e BuildStagesRequest_StageConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sBuildStagesRequest_StageConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = BuildStagesRequest_StageConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = BuildStagesRequest_StageConfigValidationError{}
