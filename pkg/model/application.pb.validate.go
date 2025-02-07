// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/model/application.proto

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

// Validate checks the field values on Application with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Application) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Application with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ApplicationMultiError, or
// nil if none found.
func (m *Application) ValidateAll() error {
	return m.validate(true)
}

func (m *Application) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := ApplicationValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := ApplicationValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetPipedId()) < 1 {
		err := ApplicationValidationError{
			field:  "PipedId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetProjectId()) < 1 {
		err := ApplicationValidationError{
			field:  "ProjectId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := ApplicationKind_name[int32(m.GetKind())]; !ok {
		err := ApplicationValidationError{
			field:  "Kind",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetGitPath() == nil {
		err := ApplicationValidationError{
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
				errors = append(errors, ApplicationValidationError{
					field:  "GitPath",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "GitPath",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGitPath()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ApplicationValidationError{
				field:  "GitPath",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for CloudProvider

	// no validation rules for PlatformProvider

	// no validation rules for Description

	// no validation rules for Labels

	if all {
		switch v := interface{}(m.GetMostRecentlySuccessfulDeployment()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "MostRecentlySuccessfulDeployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "MostRecentlySuccessfulDeployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMostRecentlySuccessfulDeployment()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ApplicationValidationError{
				field:  "MostRecentlySuccessfulDeployment",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMostRecentlyTriggeredDeployment()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "MostRecentlyTriggeredDeployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "MostRecentlyTriggeredDeployment",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMostRecentlyTriggeredDeployment()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ApplicationValidationError{
				field:  "MostRecentlyTriggeredDeployment",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetSyncState()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "SyncState",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ApplicationValidationError{
					field:  "SyncState",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSyncState()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ApplicationValidationError{
				field:  "SyncState",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Deploying

	if m.GetDeletedAt() < 0 {
		err := ApplicationValidationError{
			field:  "DeletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Deleted

	// no validation rules for Disabled

	if m.GetCreatedAt() <= 0 {
		err := ApplicationValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() <= 0 {
		err := ApplicationValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ApplicationMultiError(errors)
	}

	return nil
}

// ApplicationMultiError is an error wrapping multiple validation errors
// returned by Application.ValidateAll() if the designated constraints aren't met.
type ApplicationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ApplicationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ApplicationMultiError) AllErrors() []error { return m }

// ApplicationValidationError is the validation error returned by
// Application.Validate if the designated constraints aren't met.
type ApplicationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationValidationError) ErrorName() string { return "ApplicationValidationError" }

// Error satisfies the builtin error interface
func (e ApplicationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplication.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationValidationError{}

// Validate checks the field values on ApplicationSyncState with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ApplicationSyncState) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ApplicationSyncState with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ApplicationSyncStateMultiError, or nil if none found.
func (m *ApplicationSyncState) ValidateAll() error {
	return m.validate(true)
}

func (m *ApplicationSyncState) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if _, ok := ApplicationSyncStatus_name[int32(m.GetStatus())]; !ok {
		err := ApplicationSyncStateValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for ShortReason

	// no validation rules for Reason

	// no validation rules for HeadDeploymentId

	if m.GetTimestamp() <= 0 {
		err := ApplicationSyncStateValidationError{
			field:  "Timestamp",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ApplicationSyncStateMultiError(errors)
	}

	return nil
}

// ApplicationSyncStateMultiError is an error wrapping multiple validation
// errors returned by ApplicationSyncState.ValidateAll() if the designated
// constraints aren't met.
type ApplicationSyncStateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ApplicationSyncStateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ApplicationSyncStateMultiError) AllErrors() []error { return m }

// ApplicationSyncStateValidationError is the validation error returned by
// ApplicationSyncState.Validate if the designated constraints aren't met.
type ApplicationSyncStateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationSyncStateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationSyncStateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationSyncStateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationSyncStateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationSyncStateValidationError) ErrorName() string {
	return "ApplicationSyncStateValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationSyncStateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationSyncState.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationSyncStateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationSyncStateValidationError{}

// Validate checks the field values on ApplicationDeploymentReference with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ApplicationDeploymentReference) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ApplicationDeploymentReference with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// ApplicationDeploymentReferenceMultiError, or nil if none found.
func (m *ApplicationDeploymentReference) ValidateAll() error {
	return m.validate(true)
}

func (m *ApplicationDeploymentReference) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetDeploymentId()) < 1 {
		err := ApplicationDeploymentReferenceValidationError{
			field:  "DeploymentId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTrigger() == nil {
		err := ApplicationDeploymentReferenceValidationError{
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
				errors = append(errors, ApplicationDeploymentReferenceValidationError{
					field:  "Trigger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ApplicationDeploymentReferenceValidationError{
					field:  "Trigger",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTrigger()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ApplicationDeploymentReferenceValidationError{
				field:  "Trigger",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Summary

	// no validation rules for Version

	// no validation rules for ConfigFilename

	for idx, item := range m.GetVersions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ApplicationDeploymentReferenceValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ApplicationDeploymentReferenceValidationError{
						field:  fmt.Sprintf("Versions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ApplicationDeploymentReferenceValidationError{
					field:  fmt.Sprintf("Versions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.GetStartedAt() <= 0 {
		err := ApplicationDeploymentReferenceValidationError{
			field:  "StartedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCompletedAt() < 0 {
		err := ApplicationDeploymentReferenceValidationError{
			field:  "CompletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ApplicationDeploymentReferenceMultiError(errors)
	}

	return nil
}

// ApplicationDeploymentReferenceMultiError is an error wrapping multiple
// validation errors returned by ApplicationDeploymentReference.ValidateAll()
// if the designated constraints aren't met.
type ApplicationDeploymentReferenceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ApplicationDeploymentReferenceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ApplicationDeploymentReferenceMultiError) AllErrors() []error { return m }

// ApplicationDeploymentReferenceValidationError is the validation error
// returned by ApplicationDeploymentReference.Validate if the designated
// constraints aren't met.
type ApplicationDeploymentReferenceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationDeploymentReferenceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationDeploymentReferenceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationDeploymentReferenceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationDeploymentReferenceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationDeploymentReferenceValidationError) ErrorName() string {
	return "ApplicationDeploymentReferenceValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationDeploymentReferenceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationDeploymentReference.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationDeploymentReferenceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationDeploymentReferenceValidationError{}
