// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/model/insight.proto

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

// Validate checks the field values on InsightSample with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *InsightSample) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightSample with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in InsightSampleMultiError, or
// nil if none found.
func (m *InsightSample) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightSample) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Labels

	if all {
		switch v := interface{}(m.GetDataPoint()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, InsightSampleValidationError{
					field:  "DataPoint",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, InsightSampleValidationError{
					field:  "DataPoint",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDataPoint()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InsightSampleValidationError{
				field:  "DataPoint",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return InsightSampleMultiError(errors)
	}

	return nil
}

// InsightSampleMultiError is an error wrapping multiple validation errors
// returned by InsightSample.ValidateAll() if the designated constraints
// aren't met.
type InsightSampleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightSampleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightSampleMultiError) AllErrors() []error { return m }

// InsightSampleValidationError is the validation error returned by
// InsightSample.Validate if the designated constraints aren't met.
type InsightSampleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightSampleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightSampleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightSampleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightSampleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightSampleValidationError) ErrorName() string { return "InsightSampleValidationError" }

// Error satisfies the builtin error interface
func (e InsightSampleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightSample.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightSampleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightSampleValidationError{}

// Validate checks the field values on InsightSampleStream with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightSampleStream) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightSampleStream with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightSampleStreamMultiError, or nil if none found.
func (m *InsightSampleStream) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightSampleStream) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Labels

	for idx, item := range m.GetDataPoints() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, InsightSampleStreamValidationError{
						field:  fmt.Sprintf("DataPoints[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, InsightSampleStreamValidationError{
						field:  fmt.Sprintf("DataPoints[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InsightSampleStreamValidationError{
					field:  fmt.Sprintf("DataPoints[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return InsightSampleStreamMultiError(errors)
	}

	return nil
}

// InsightSampleStreamMultiError is an error wrapping multiple validation
// errors returned by InsightSampleStream.ValidateAll() if the designated
// constraints aren't met.
type InsightSampleStreamMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightSampleStreamMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightSampleStreamMultiError) AllErrors() []error { return m }

// InsightSampleStreamValidationError is the validation error returned by
// InsightSampleStream.Validate if the designated constraints aren't met.
type InsightSampleStreamValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightSampleStreamValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightSampleStreamValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightSampleStreamValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightSampleStreamValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightSampleStreamValidationError) ErrorName() string {
	return "InsightSampleStreamValidationError"
}

// Error satisfies the builtin error interface
func (e InsightSampleStreamValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightSampleStream.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightSampleStreamValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightSampleStreamValidationError{}

// Validate checks the field values on InsightDataPoint with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *InsightDataPoint) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDataPoint with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightDataPointMultiError, or nil if none found.
func (m *InsightDataPoint) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDataPoint) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetTimestamp() <= 0 {
		err := InsightDataPointValidationError{
			field:  "Timestamp",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetValue() <= 0 {
		err := InsightDataPointValidationError{
			field:  "Value",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return InsightDataPointMultiError(errors)
	}

	return nil
}

// InsightDataPointMultiError is an error wrapping multiple validation errors
// returned by InsightDataPoint.ValidateAll() if the designated constraints
// aren't met.
type InsightDataPointMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDataPointMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDataPointMultiError) AllErrors() []error { return m }

// InsightDataPointValidationError is the validation error returned by
// InsightDataPoint.Validate if the designated constraints aren't met.
type InsightDataPointValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDataPointValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDataPointValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDataPointValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDataPointValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDataPointValidationError) ErrorName() string { return "InsightDataPointValidationError" }

// Error satisfies the builtin error interface
func (e InsightDataPointValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDataPoint.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDataPointValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDataPointValidationError{}

// Validate checks the field values on InsightApplicationCount with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightApplicationCount) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightApplicationCount with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightApplicationCountMultiError, or nil if none found.
func (m *InsightApplicationCount) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightApplicationCount) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Labels

	// no validation rules for Count

	if len(errors) > 0 {
		return InsightApplicationCountMultiError(errors)
	}

	return nil
}

// InsightApplicationCountMultiError is an error wrapping multiple validation
// errors returned by InsightApplicationCount.ValidateAll() if the designated
// constraints aren't met.
type InsightApplicationCountMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightApplicationCountMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightApplicationCountMultiError) AllErrors() []error { return m }

// InsightApplicationCountValidationError is the validation error returned by
// InsightApplicationCount.Validate if the designated constraints aren't met.
type InsightApplicationCountValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightApplicationCountValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightApplicationCountValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightApplicationCountValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightApplicationCountValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightApplicationCountValidationError) ErrorName() string {
	return "InsightApplicationCountValidationError"
}

// Error satisfies the builtin error interface
func (e InsightApplicationCountValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightApplicationCount.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightApplicationCountValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightApplicationCountValidationError{}

// Validate checks the field values on InsightDeploymentChunk with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightDeploymentChunk) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDeploymentChunk with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightDeploymentChunkMultiError, or nil if none found.
func (m *InsightDeploymentChunk) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentChunk) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Version

	if m.GetFrom() < 0 {
		err := InsightDeploymentChunkValidationError{
			field:  "From",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTo() < 0 {
		err := InsightDeploymentChunkValidationError{
			field:  "To",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetDeployments() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, InsightDeploymentChunkValidationError{
						field:  fmt.Sprintf("Deployments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, InsightDeploymentChunkValidationError{
						field:  fmt.Sprintf("Deployments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InsightDeploymentChunkValidationError{
					field:  fmt.Sprintf("Deployments[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return InsightDeploymentChunkMultiError(errors)
	}

	return nil
}

// InsightDeploymentChunkMultiError is an error wrapping multiple validation
// errors returned by InsightDeploymentChunk.ValidateAll() if the designated
// constraints aren't met.
type InsightDeploymentChunkMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentChunkMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentChunkMultiError) AllErrors() []error { return m }

// InsightDeploymentChunkValidationError is the validation error returned by
// InsightDeploymentChunk.Validate if the designated constraints aren't met.
type InsightDeploymentChunkValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentChunkValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentChunkValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentChunkValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentChunkValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentChunkValidationError) ErrorName() string {
	return "InsightDeploymentChunkValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentChunkValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentChunk.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentChunkValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentChunkValidationError{}

// Validate checks the field values on InsightDeployment with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *InsightDeployment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDeployment with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightDeploymentMultiError, or nil if none found.
func (m *InsightDeployment) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeployment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for AppId

	// no validation rules for Labels

	// no validation rules for StartedAt

	// no validation rules for CompletedAt

	// no validation rules for RollbackStartedAt

	// no validation rules for CompleteStatus

	if len(errors) > 0 {
		return InsightDeploymentMultiError(errors)
	}

	return nil
}

// InsightDeploymentMultiError is an error wrapping multiple validation errors
// returned by InsightDeployment.ValidateAll() if the designated constraints
// aren't met.
type InsightDeploymentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentMultiError) AllErrors() []error { return m }

// InsightDeploymentValidationError is the validation error returned by
// InsightDeployment.Validate if the designated constraints aren't met.
type InsightDeploymentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentValidationError) ErrorName() string {
	return "InsightDeploymentValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeployment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentValidationError{}

// Validate checks the field values on InsightDeploymentChunkMetadata with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightDeploymentChunkMetadata) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDeploymentChunkMetadata with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// InsightDeploymentChunkMetadataMultiError, or nil if none found.
func (m *InsightDeploymentChunkMetadata) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentChunkMetadata) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetChunks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, InsightDeploymentChunkMetadataValidationError{
						field:  fmt.Sprintf("Chunks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, InsightDeploymentChunkMetadataValidationError{
						field:  fmt.Sprintf("Chunks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InsightDeploymentChunkMetadataValidationError{
					field:  fmt.Sprintf("Chunks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return InsightDeploymentChunkMetadataMultiError(errors)
	}

	return nil
}

// InsightDeploymentChunkMetadataMultiError is an error wrapping multiple
// validation errors returned by InsightDeploymentChunkMetadata.ValidateAll()
// if the designated constraints aren't met.
type InsightDeploymentChunkMetadataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentChunkMetadataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentChunkMetadataMultiError) AllErrors() []error { return m }

// InsightDeploymentChunkMetadataValidationError is the validation error
// returned by InsightDeploymentChunkMetadata.Validate if the designated
// constraints aren't met.
type InsightDeploymentChunkMetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentChunkMetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentChunkMetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentChunkMetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentChunkMetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentChunkMetadataValidationError) ErrorName() string {
	return "InsightDeploymentChunkMetadataValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentChunkMetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentChunkMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentChunkMetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentChunkMetadataValidationError{}

// Validate checks the field values on InsightDeploymentChunkMetadata_ChunkMeta
// with the rules defined in the proto definition for this message. If any
// rules are violated, the first error encountered is returned, or nil if
// there are no violations.
func (m *InsightDeploymentChunkMetadata_ChunkMeta) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on
// InsightDeploymentChunkMetadata_ChunkMeta with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in
// InsightDeploymentChunkMetadata_ChunkMetaMultiError, or nil if none found.
func (m *InsightDeploymentChunkMetadata_ChunkMeta) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentChunkMetadata_ChunkMeta) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetFrom() < 0 {
		err := InsightDeploymentChunkMetadata_ChunkMetaValidationError{
			field:  "From",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTo() < 0 {
		err := InsightDeploymentChunkMetadata_ChunkMetaValidationError{
			field:  "To",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Name

	// no validation rules for Size

	// no validation rules for Count

	if len(errors) > 0 {
		return InsightDeploymentChunkMetadata_ChunkMetaMultiError(errors)
	}

	return nil
}

// InsightDeploymentChunkMetadata_ChunkMetaMultiError is an error wrapping
// multiple validation errors returned by
// InsightDeploymentChunkMetadata_ChunkMeta.ValidateAll() if the designated
// constraints aren't met.
type InsightDeploymentChunkMetadata_ChunkMetaMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentChunkMetadata_ChunkMetaMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentChunkMetadata_ChunkMetaMultiError) AllErrors() []error { return m }

// InsightDeploymentChunkMetadata_ChunkMetaValidationError is the validation
// error returned by InsightDeploymentChunkMetadata_ChunkMeta.Validate if the
// designated constraints aren't met.
type InsightDeploymentChunkMetadata_ChunkMetaValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) ErrorName() string {
	return "InsightDeploymentChunkMetadata_ChunkMetaValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentChunkMetadata_ChunkMetaValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentChunkMetadata_ChunkMeta.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentChunkMetadata_ChunkMetaValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentChunkMetadata_ChunkMetaValidationError{}
