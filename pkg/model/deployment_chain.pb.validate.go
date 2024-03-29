// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/model/deployment_chain.proto

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

// Validate checks the field values on DeploymentChain with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeploymentChain) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeploymentChain with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeploymentChainMultiError, or nil if none found.
func (m *DeploymentChain) ValidateAll() error {
	return m.validate(true)
}

func (m *DeploymentChain) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := DeploymentChainValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetProjectId()) < 1 {
		err := DeploymentChainValidationError{
			field:  "ProjectId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := ChainStatus_name[int32(m.GetStatus())]; !ok {
		err := DeploymentChainValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetBlocks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DeploymentChainValidationError{
						field:  fmt.Sprintf("Blocks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DeploymentChainValidationError{
						field:  fmt.Sprintf("Blocks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DeploymentChainValidationError{
					field:  fmt.Sprintf("Blocks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.GetCompletedAt() < 0 {
		err := DeploymentChainValidationError{
			field:  "CompletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreatedAt() <= 0 {
		err := DeploymentChainValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() <= 0 {
		err := DeploymentChainValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeploymentChainMultiError(errors)
	}

	return nil
}

// DeploymentChainMultiError is an error wrapping multiple validation errors
// returned by DeploymentChain.ValidateAll() if the designated constraints
// aren't met.
type DeploymentChainMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeploymentChainMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeploymentChainMultiError) AllErrors() []error { return m }

// DeploymentChainValidationError is the validation error returned by
// DeploymentChain.Validate if the designated constraints aren't met.
type DeploymentChainValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeploymentChainValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeploymentChainValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeploymentChainValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeploymentChainValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeploymentChainValidationError) ErrorName() string { return "DeploymentChainValidationError" }

// Error satisfies the builtin error interface
func (e DeploymentChainValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeploymentChain.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeploymentChainValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeploymentChainValidationError{}

// Validate checks the field values on ChainApplicationRef with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ChainApplicationRef) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChainApplicationRef with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ChainApplicationRefMultiError, or nil if none found.
func (m *ChainApplicationRef) ValidateAll() error {
	return m.validate(true)
}

func (m *ChainApplicationRef) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetApplicationId()) < 1 {
		err := ChainApplicationRefValidationError{
			field:  "ApplicationId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for ApplicationName

	if len(errors) > 0 {
		return ChainApplicationRefMultiError(errors)
	}

	return nil
}

// ChainApplicationRefMultiError is an error wrapping multiple validation
// errors returned by ChainApplicationRef.ValidateAll() if the designated
// constraints aren't met.
type ChainApplicationRefMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChainApplicationRefMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChainApplicationRefMultiError) AllErrors() []error { return m }

// ChainApplicationRefValidationError is the validation error returned by
// ChainApplicationRef.Validate if the designated constraints aren't met.
type ChainApplicationRefValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChainApplicationRefValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChainApplicationRefValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChainApplicationRefValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChainApplicationRefValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChainApplicationRefValidationError) ErrorName() string {
	return "ChainApplicationRefValidationError"
}

// Error satisfies the builtin error interface
func (e ChainApplicationRefValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChainApplicationRef.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChainApplicationRefValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChainApplicationRefValidationError{}

// Validate checks the field values on ChainDeploymentRef with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ChainDeploymentRef) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChainDeploymentRef with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ChainDeploymentRefMultiError, or nil if none found.
func (m *ChainDeploymentRef) ValidateAll() error {
	return m.validate(true)
}

func (m *ChainDeploymentRef) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetDeploymentId()) < 1 {
		err := ChainDeploymentRefValidationError{
			field:  "DeploymentId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := DeploymentStatus_name[int32(m.GetStatus())]; !ok {
		err := ChainDeploymentRefValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for StatusReason

	if len(errors) > 0 {
		return ChainDeploymentRefMultiError(errors)
	}

	return nil
}

// ChainDeploymentRefMultiError is an error wrapping multiple validation errors
// returned by ChainDeploymentRef.ValidateAll() if the designated constraints
// aren't met.
type ChainDeploymentRefMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChainDeploymentRefMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChainDeploymentRefMultiError) AllErrors() []error { return m }

// ChainDeploymentRefValidationError is the validation error returned by
// ChainDeploymentRef.Validate if the designated constraints aren't met.
type ChainDeploymentRefValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChainDeploymentRefValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChainDeploymentRefValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChainDeploymentRefValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChainDeploymentRefValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChainDeploymentRefValidationError) ErrorName() string {
	return "ChainDeploymentRefValidationError"
}

// Error satisfies the builtin error interface
func (e ChainDeploymentRefValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChainDeploymentRef.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChainDeploymentRefValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChainDeploymentRefValidationError{}

// Validate checks the field values on ChainNode with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ChainNode) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChainNode with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ChainNodeMultiError, or nil
// if none found.
func (m *ChainNode) ValidateAll() error {
	return m.validate(true)
}

func (m *ChainNode) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetApplicationRef() == nil {
		err := ChainNodeValidationError{
			field:  "ApplicationRef",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetApplicationRef()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ChainNodeValidationError{
					field:  "ApplicationRef",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ChainNodeValidationError{
					field:  "ApplicationRef",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetApplicationRef()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ChainNodeValidationError{
				field:  "ApplicationRef",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDeploymentRef()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ChainNodeValidationError{
					field:  "DeploymentRef",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ChainNodeValidationError{
					field:  "DeploymentRef",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDeploymentRef()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ChainNodeValidationError{
				field:  "DeploymentRef",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ChainNodeMultiError(errors)
	}

	return nil
}

// ChainNodeMultiError is an error wrapping multiple validation errors returned
// by ChainNode.ValidateAll() if the designated constraints aren't met.
type ChainNodeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChainNodeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChainNodeMultiError) AllErrors() []error { return m }

// ChainNodeValidationError is the validation error returned by
// ChainNode.Validate if the designated constraints aren't met.
type ChainNodeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChainNodeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChainNodeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChainNodeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChainNodeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChainNodeValidationError) ErrorName() string { return "ChainNodeValidationError" }

// Error satisfies the builtin error interface
func (e ChainNodeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChainNode.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChainNodeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChainNodeValidationError{}

// Validate checks the field values on ChainBlock with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ChainBlock) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChainBlock with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ChainBlockMultiError, or
// nil if none found.
func (m *ChainBlock) ValidateAll() error {
	return m.validate(true)
}

func (m *ChainBlock) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetNodes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ChainBlockValidationError{
						field:  fmt.Sprintf("Nodes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ChainBlockValidationError{
						field:  fmt.Sprintf("Nodes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ChainBlockValidationError{
					field:  fmt.Sprintf("Nodes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if _, ok := ChainBlockStatus_name[int32(m.GetStatus())]; !ok {
		err := ChainBlockValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetStartedAt() < 0 {
		err := ChainBlockValidationError{
			field:  "StartedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCompletedAt() < 0 {
		err := ChainBlockValidationError{
			field:  "CompletedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ChainBlockMultiError(errors)
	}

	return nil
}

// ChainBlockMultiError is an error wrapping multiple validation errors
// returned by ChainBlock.ValidateAll() if the designated constraints aren't met.
type ChainBlockMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChainBlockMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChainBlockMultiError) AllErrors() []error { return m }

// ChainBlockValidationError is the validation error returned by
// ChainBlock.Validate if the designated constraints aren't met.
type ChainBlockValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChainBlockValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChainBlockValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChainBlockValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChainBlockValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChainBlockValidationError) ErrorName() string { return "ChainBlockValidationError" }

// Error satisfies the builtin error interface
func (e ChainBlockValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChainBlock.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChainBlockValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChainBlockValidationError{}
