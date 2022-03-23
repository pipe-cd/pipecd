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

// Validate checks the field values on InsightDeploymentSubset with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightDeploymentSubset) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDeploymentSubset with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightDeploymentSubsetMultiError, or nil if none found.
func (m *InsightDeploymentSubset) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentSubset) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := InsightDeploymentSubsetValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreatedAt() < 0 {
		err := InsightDeploymentSubsetValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() < 0 {
		err := InsightDeploymentSubsetValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return InsightDeploymentSubsetMultiError(errors)
	}

	return nil
}

// InsightDeploymentSubsetMultiError is an error wrapping multiple validation
// errors returned by InsightDeploymentSubset.ValidateAll() if the designated
// constraints aren't met.
type InsightDeploymentSubsetMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentSubsetMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentSubsetMultiError) AllErrors() []error { return m }

// InsightDeploymentSubsetValidationError is the validation error returned by
// InsightDeploymentSubset.Validate if the designated constraints aren't met.
type InsightDeploymentSubsetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentSubsetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentSubsetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentSubsetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentSubsetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentSubsetValidationError) ErrorName() string {
	return "InsightDeploymentSubsetValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentSubsetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentSubset.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentSubsetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentSubsetValidationError{}

// Validate checks the field values on InsightDailyDeployment with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightDailyDeployment) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDailyDeployment with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightDailyDeploymentMultiError, or nil if none found.
func (m *InsightDailyDeployment) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDailyDeployment) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetDate() < 0 {
		err := InsightDailyDeploymentValidationError{
			field:  "Date",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetCreatedAt() < 0 {
		err := InsightDailyDeploymentValidationError{
			field:  "CreatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetUpdatedAt() < 0 {
		err := InsightDailyDeploymentValidationError{
			field:  "UpdatedAt",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetDailyDeployments() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, InsightDailyDeploymentValidationError{
						field:  fmt.Sprintf("DailyDeployments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, InsightDailyDeploymentValidationError{
						field:  fmt.Sprintf("DailyDeployments[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InsightDailyDeploymentValidationError{
					field:  fmt.Sprintf("DailyDeployments[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return InsightDailyDeploymentMultiError(errors)
	}

	return nil
}

// InsightDailyDeploymentMultiError is an error wrapping multiple validation
// errors returned by InsightDailyDeployment.ValidateAll() if the designated
// constraints aren't met.
type InsightDailyDeploymentMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDailyDeploymentMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDailyDeploymentMultiError) AllErrors() []error { return m }

// InsightDailyDeploymentValidationError is the validation error returned by
// InsightDailyDeployment.Validate if the designated constraints aren't met.
type InsightDailyDeploymentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDailyDeploymentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDailyDeploymentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDailyDeploymentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDailyDeploymentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDailyDeploymentValidationError) ErrorName() string {
	return "InsightDailyDeploymentValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDailyDeploymentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDailyDeployment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDailyDeploymentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDailyDeploymentValidationError{}

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

	if all {
		switch v := interface{}(m.GetDateRange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, InsightDeploymentChunkValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, InsightDeploymentChunkValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDateRange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InsightDeploymentChunkValidationError{
				field:  "DateRange",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
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

// Validate checks the field values on InsightDeploymentChunkMetaData with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightDeploymentChunkMetaData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightDeploymentChunkMetaData with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// InsightDeploymentChunkMetaDataMultiError, or nil if none found.
func (m *InsightDeploymentChunkMetaData) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentChunkMetaData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, InsightDeploymentChunkMetaDataValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, InsightDeploymentChunkMetaDataValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InsightDeploymentChunkMetaDataValidationError{
					field:  fmt.Sprintf("Data[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return InsightDeploymentChunkMetaDataMultiError(errors)
	}

	return nil
}

// InsightDeploymentChunkMetaDataMultiError is an error wrapping multiple
// validation errors returned by InsightDeploymentChunkMetaData.ValidateAll()
// if the designated constraints aren't met.
type InsightDeploymentChunkMetaDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentChunkMetaDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentChunkMetaDataMultiError) AllErrors() []error { return m }

// InsightDeploymentChunkMetaDataValidationError is the validation error
// returned by InsightDeploymentChunkMetaData.Validate if the designated
// constraints aren't met.
type InsightDeploymentChunkMetaDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentChunkMetaDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightDeploymentChunkMetaDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightDeploymentChunkMetaDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentChunkMetaDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentChunkMetaDataValidationError) ErrorName() string {
	return "InsightDeploymentChunkMetaDataValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentChunkMetaDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentChunkMetaData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentChunkMetaDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentChunkMetaDataValidationError{}

// Validate checks the field values on InsightChunkDateRange with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *InsightChunkDateRange) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on InsightChunkDateRange with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// InsightChunkDateRangeMultiError, or nil if none found.
func (m *InsightChunkDateRange) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightChunkDateRange) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetFrom() < 0 {
		err := InsightChunkDateRangeValidationError{
			field:  "From",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTo() < 0 {
		err := InsightChunkDateRangeValidationError{
			field:  "To",
			reason: "value must be greater than or equal to 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return InsightChunkDateRangeMultiError(errors)
	}

	return nil
}

// InsightChunkDateRangeMultiError is an error wrapping multiple validation
// errors returned by InsightChunkDateRange.ValidateAll() if the designated
// constraints aren't met.
type InsightChunkDateRangeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightChunkDateRangeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightChunkDateRangeMultiError) AllErrors() []error { return m }

// InsightChunkDateRangeValidationError is the validation error returned by
// InsightChunkDateRange.Validate if the designated constraints aren't met.
type InsightChunkDateRangeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightChunkDateRangeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InsightChunkDateRangeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InsightChunkDateRangeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightChunkDateRangeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightChunkDateRangeValidationError) ErrorName() string {
	return "InsightChunkDateRangeValidationError"
}

// Error satisfies the builtin error interface
func (e InsightChunkDateRangeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightChunkDateRange.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightChunkDateRangeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightChunkDateRangeValidationError{}

// Validate checks the field values on
// InsightDeploymentChunkMetaData_InsightChunkData with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *InsightDeploymentChunkMetaData_InsightChunkData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on
// InsightDeploymentChunkMetaData_InsightChunkData with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in
// InsightDeploymentChunkMetaData_InsightChunkDataMultiError, or nil if none found.
func (m *InsightDeploymentChunkMetaData_InsightChunkData) ValidateAll() error {
	return m.validate(true)
}

func (m *InsightDeploymentChunkMetaData_InsightChunkData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetDateRange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, InsightDeploymentChunkMetaData_InsightChunkDataValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, InsightDeploymentChunkMetaData_InsightChunkDataValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDateRange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InsightDeploymentChunkMetaData_InsightChunkDataValidationError{
				field:  "DateRange",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ChunkKey

	// no validation rules for ChunkSize

	if len(errors) > 0 {
		return InsightDeploymentChunkMetaData_InsightChunkDataMultiError(errors)
	}

	return nil
}

// InsightDeploymentChunkMetaData_InsightChunkDataMultiError is an error
// wrapping multiple validation errors returned by
// InsightDeploymentChunkMetaData_InsightChunkData.ValidateAll() if the
// designated constraints aren't met.
type InsightDeploymentChunkMetaData_InsightChunkDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m InsightDeploymentChunkMetaData_InsightChunkDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m InsightDeploymentChunkMetaData_InsightChunkDataMultiError) AllErrors() []error { return m }

// InsightDeploymentChunkMetaData_InsightChunkDataValidationError is the
// validation error returned by
// InsightDeploymentChunkMetaData_InsightChunkData.Validate if the designated
// constraints aren't met.
type InsightDeploymentChunkMetaData_InsightChunkDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) Field() string {
	return e.field
}

// Reason function returns reason value.
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) Reason() string {
	return e.reason
}

// Cause function returns cause value.
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) ErrorName() string {
	return "InsightDeploymentChunkMetaData_InsightChunkDataValidationError"
}

// Error satisfies the builtin error interface
func (e InsightDeploymentChunkMetaData_InsightChunkDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInsightDeploymentChunkMetaData_InsightChunkData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InsightDeploymentChunkMetaData_InsightChunkDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InsightDeploymentChunkMetaData_InsightChunkDataValidationError{}
