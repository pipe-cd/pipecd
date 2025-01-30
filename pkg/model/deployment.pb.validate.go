// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/model/deployment.proto

package model

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
)

// Validate checks the field values on Deployment with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Deployment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Deployment with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeploymentMultiError, or
// nil if none found.
func (m *Deployment) ValidateAll() error {
	return m.validate(true)
}

func (m *Deployment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := DeploymentValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetApplicationId()) < 1 {
		err := DeploymentValidationError{
			field:  "ApplicationId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetApplicationName()) < 1 {
		err := DeploymentValidationError{
			field:  "ApplicationName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetPipedId()) < 1 {
		err := DeploymentValidationError{
			field:  "PipedId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetProjectId()) < 1 {
		err := DeploymentValidationError{
			field:  "ProjectId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := ApplicationKind_name[int32(m.GetKind())]; !ok {
		err := DeploymentValidationError{
			field:  "Kind",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetGitPath() == nil {
		err := DeploymentValidationError{
			field:  "GitPath",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetGitPath()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "GitPath",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "GitPath",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGitPath()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeploymentValidationError{
				field:  "GitPath",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for CloudProvider

	// no validation rules for PlatformProvider

	// no validation rules for Labels

	if m.GetTrigger() == nil {
		err := DeploymentValidationError{
			field:  "Trigger",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetTrigger()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "Trigger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "Trigger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTrigger()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeploymentValidationError{
				field:  "Trigger",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Summary

	// no validation rules for Version

	for idx, item := range m.GetVersions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DeploymentValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DeploymentValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DeploymentValidationError{
					field:  fmt.Sprintf("Versions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for RunningCommitHash

	// no validation rules for RunningConfigFilename

	if _, ok := DeploymentStatus_name[int32(m.GetStatus())]; !ok {
		err := DeploymentValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for StatusReason

	for idx, item := range m.GetStages() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DeploymentValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DeploymentValidationError{
						field:  fmt.Sprintf("Stages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DeploymentValidationError{
					field:  fmt.Sprintf("Stages[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Metadata

	if all {
		switch v := interface{}(m.GetMetadataV2()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "MetadataV2",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeploymentValidationError{
					field:  "MetadataV2",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetadataV2()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeploymentValidationError{
				field:  "MetadataV2",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for DeploymentChainId

	// no validation rules for DeploymentChainBlockIndex

	if m.GetCompletedAt() < 0 {
		err := DeploymentValidationError{
			field:  "CompletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreatedAt() < 0 {
		err := DeploymentValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() < 0 {
		err := DeploymentValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeploymentMultiError(errors)
	}

	return nil
}

// DeploymentMultiError is an error wrapping multiple validation errors
// returned by Deployment.ValidateAll() if the designated constraints aren't met.
type DeploymentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeploymentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeploymentMultiError) AllErrors() []error { return m }

// DeploymentValidationError is the validation error returned by
// Deployment.Validate if the designated constraints aren't met.
type DeploymentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeploymentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeploymentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeploymentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeploymentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeploymentValidationError) ErrorName() string { return "DeploymentValidationError" }

// Error satisfies the builtin error interface
func (e DeploymentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeployment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeploymentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeploymentValidationError{}

// Validate checks the field values on DeploymentTrigger with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeploymentTrigger) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeploymentTrigger with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeploymentTriggerMultiError, or nil if none found.
func (m *DeploymentTrigger) ValidateAll() error {
	return m.validate(true)
}

func (m *DeploymentTrigger) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetCommit() == nil {
		err := DeploymentTriggerValidationError{
			field:  "Commit",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCommit()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeploymentTriggerValidationError{
					field:  "Commit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeploymentTriggerValidationError{
					field:  "Commit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommit()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeploymentTriggerValidationError{
				field:  "Commit",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Commander

	if m.GetTimestamp() <= 0 {
		err := DeploymentTriggerValidationError{
			field:  "Timestamp",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for SyncStrategy

	// no validation rules for StrategySummary

	if len(errors) > 0 {
		return DeploymentTriggerMultiError(errors)
	}

	return nil
}

// DeploymentTriggerMultiError is an error wrapping multiple validation errors
// returned by DeploymentTrigger.ValidateAll() if the designated constraints
// aren't met.
type DeploymentTriggerMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeploymentTriggerMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeploymentTriggerMultiError) AllErrors() []error { return m }

// DeploymentTriggerValidationError is the validation error returned by
// DeploymentTrigger.Validate if the designated constraints aren't met.
type DeploymentTriggerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeploymentTriggerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeploymentTriggerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeploymentTriggerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeploymentTriggerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeploymentTriggerValidationError) ErrorName() string {
	return "DeploymentTriggerValidationError"
}

// Error satisfies the builtin error interface
func (e DeploymentTriggerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeploymentTrigger.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeploymentTriggerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeploymentTriggerValidationError{}

// Validate checks the field values on PipelineStage with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PipelineStage) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PipelineStage with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PipelineStageMultiError, or
// nil if none found.
func (m *PipelineStage) ValidateAll() error {
	return m.validate(true)
}

func (m *PipelineStage) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := PipelineStageValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := PipelineStageValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Desc

	// no validation rules for Index

	// no validation rules for Predefined

	// no validation rules for Visible

	if _, ok := StageStatus_name[int32(m.GetStatus())]; !ok {
		err := PipelineStageValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for StatusReason

	// no validation rules for Metadata

	// no validation rules for RetriedCount

	// no validation rules for Rollback

	if m.GetCompletedAt() < 0 {
		err := PipelineStageValidationError{
			field:  "CompletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreatedAt() <= 0 {
		err := PipelineStageValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() <= 0 {
		err := PipelineStageValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Skippable

	// no validation rules for Approvable

	if len(errors) > 0 {
		return PipelineStageMultiError(errors)
	}

	return nil
}

// PipelineStageMultiError is an error wrapping multiple validation errors
// returned by PipelineStage.ValidateAll() if the designated constraints
// aren't met.
type PipelineStageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PipelineStageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PipelineStageMultiError) AllErrors() []error { return m }

// PipelineStageValidationError is the validation error returned by
// PipelineStage.Validate if the designated constraints aren't met.
type PipelineStageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PipelineStageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PipelineStageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PipelineStageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PipelineStageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PipelineStageValidationError) ErrorName() string { return "PipelineStageValidationError" }

// Error satisfies the builtin error interface
func (e PipelineStageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPipelineStage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PipelineStageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PipelineStageValidationError{}

// Validate checks the field values on Commit with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Commit) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Commit with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CommitMultiError, or nil if none found.
func (m *Commit) ValidateAll() error {
	return m.validate(true)
}

func (m *Commit) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetHash()) < 1 {
		err := CommitValidationError{
			field:  "Hash",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetMessage()) < 1 {
		err := CommitValidationError{
			field:  "Message",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetAuthor()) < 1 {
		err := CommitValidationError{
			field:  "Author",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetBranch()) < 1 {
		err := CommitValidationError{
			field:  "Branch",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for PullRequest

	// no validation rules for Url

	if m.GetCreatedAt() <= 0 {
		err := CommitValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CommitMultiError(errors)
	}

	return nil
}

// CommitMultiError is an error wrapping multiple validation errors returned by
// Commit.ValidateAll() if the designated constraints aren't met.
type CommitMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CommitMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CommitMultiError) AllErrors() []error { return m }

// CommitValidationError is the validation error returned by Commit.Validate if
// the designated constraints aren't met.
type CommitValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CommitValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CommitValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CommitValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CommitValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CommitValidationError) ErrorName() string { return "CommitValidationError" }

// Error satisfies the builtin error interface
func (e CommitValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCommit.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CommitValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CommitValidationError{}

// Validate checks the field values on DeploymentMetadata with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeploymentMetadata) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeploymentMetadata with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeploymentMetadataMultiError, or nil if none found.
func (m *DeploymentMetadata) ValidateAll() error {
	return m.validate(true)
}

func (m *DeploymentMetadata) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetShared()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DeploymentMetadataValidationError{
					field:  "Shared",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DeploymentMetadataValidationError{
					field:  "Shared",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetShared()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DeploymentMetadataValidationError{
				field:  "Shared",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	{
		sorted_keys := make([]string, len(m.GetPlugins()))
		i := 0
		for key := range m.GetPlugins() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetPlugins()[key]
			_ = val

			// no validation rules for Plugins[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, DeploymentMetadataValidationError{
							field:  fmt.Sprintf("Plugins[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, DeploymentMetadataValidationError{
							field:  fmt.Sprintf("Plugins[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return DeploymentMetadataValidationError{
						field:  fmt.Sprintf("Plugins[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return DeploymentMetadataMultiError(errors)
	}

	return nil
}

// DeploymentMetadataMultiError is an error wrapping multiple validation errors
// returned by DeploymentMetadata.ValidateAll() if the designated constraints
// aren't met.
type DeploymentMetadataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeploymentMetadataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeploymentMetadataMultiError) AllErrors() []error { return m }

// DeploymentMetadataValidationError is the validation error returned by
// DeploymentMetadata.Validate if the designated constraints aren't met.
type DeploymentMetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeploymentMetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeploymentMetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeploymentMetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeploymentMetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeploymentMetadataValidationError) ErrorName() string {
	return "DeploymentMetadataValidationError"
}

// Error satisfies the builtin error interface
func (e DeploymentMetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeploymentMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeploymentMetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeploymentMetadataValidationError{}

// Validate checks the field values on DeploymentMetadata_KeyValues with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeploymentMetadata_KeyValues) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeploymentMetadata_KeyValues with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeploymentMetadata_KeyValuesMultiError, or nil if none found.
func (m *DeploymentMetadata_KeyValues) ValidateAll() error {
	return m.validate(true)
}

func (m *DeploymentMetadata_KeyValues) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for KeyValues

	if len(errors) > 0 {
		return DeploymentMetadata_KeyValuesMultiError(errors)
	}

	return nil
}

// DeploymentMetadata_KeyValuesMultiError is an error wrapping multiple
// validation errors returned by DeploymentMetadata_KeyValues.ValidateAll() if
// the designated constraints aren't met.
type DeploymentMetadata_KeyValuesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeploymentMetadata_KeyValuesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeploymentMetadata_KeyValuesMultiError) AllErrors() []error { return m }

// DeploymentMetadata_KeyValuesValidationError is the validation error returned
// by DeploymentMetadata_KeyValues.Validate if the designated constraints
// aren't met.
type DeploymentMetadata_KeyValuesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeploymentMetadata_KeyValuesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeploymentMetadata_KeyValuesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeploymentMetadata_KeyValuesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeploymentMetadata_KeyValuesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeploymentMetadata_KeyValuesValidationError) ErrorName() string {
	return "DeploymentMetadata_KeyValuesValidationError"
}

// Error satisfies the builtin error interface
func (e DeploymentMetadata_KeyValuesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeploymentMetadata_KeyValues.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeploymentMetadata_KeyValuesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeploymentMetadata_KeyValuesValidationError{}
