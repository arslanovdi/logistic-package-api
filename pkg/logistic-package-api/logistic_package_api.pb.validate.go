// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: logistic_package_api.proto

package logistic_package_api_v1

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

// Validate checks the field values on Package with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Package) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Package with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PackageMultiError, or nil if none found.
func (m *Package) ValidateAll() error {
	return m.validate(true)
}

func (m *Package) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Title

	if all {
		switch v := interface{}(m.GetCreated()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PackageValidationError{
					field:  "Created",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PackageValidationError{
					field:  "Created",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreated()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PackageValidationError{
				field:  "Created",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return PackageMultiError(errors)
	}

	return nil
}

// PackageMultiError is an error wrapping multiple validation errors returned
// by Package.ValidateAll() if the designated constraints aren't met.
type PackageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PackageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PackageMultiError) AllErrors() []error { return m }

// PackageValidationError is the validation error returned by Package.Validate
// if the designated constraints aren't met.
type PackageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PackageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PackageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PackageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PackageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PackageValidationError) ErrorName() string { return "PackageValidationError" }

// Error satisfies the builtin error interface
func (e PackageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPackage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PackageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PackageValidationError{}

// Validate checks the field values on CreatePackageRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreatePackageRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePackageRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePackageRequestV1MultiError, or nil if none found.
func (m *CreatePackageRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePackageRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetValue()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreatePackageRequestV1ValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreatePackageRequestV1ValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreatePackageRequestV1ValidationError{
				field:  "Value",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreatePackageRequestV1MultiError(errors)
	}

	return nil
}

// CreatePackageRequestV1MultiError is an error wrapping multiple validation
// errors returned by CreatePackageRequestV1.ValidateAll() if the designated
// constraints aren't met.
type CreatePackageRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePackageRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePackageRequestV1MultiError) AllErrors() []error { return m }

// CreatePackageRequestV1ValidationError is the validation error returned by
// CreatePackageRequestV1.Validate if the designated constraints aren't met.
type CreatePackageRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePackageRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePackageRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePackageRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePackageRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePackageRequestV1ValidationError) ErrorName() string {
	return "CreatePackageRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePackageRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePackageRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePackageRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePackageRequestV1ValidationError{}

// Validate checks the field values on CreatePackageResponseV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreatePackageResponseV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePackageResponseV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePackageResponseV1MultiError, or nil if none found.
func (m *CreatePackageResponseV1) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePackageResponseV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PackageId

	if len(errors) > 0 {
		return CreatePackageResponseV1MultiError(errors)
	}

	return nil
}

// CreatePackageResponseV1MultiError is an error wrapping multiple validation
// errors returned by CreatePackageResponseV1.ValidateAll() if the designated
// constraints aren't met.
type CreatePackageResponseV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePackageResponseV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePackageResponseV1MultiError) AllErrors() []error { return m }

// CreatePackageResponseV1ValidationError is the validation error returned by
// CreatePackageResponseV1.Validate if the designated constraints aren't met.
type CreatePackageResponseV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePackageResponseV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePackageResponseV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePackageResponseV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePackageResponseV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePackageResponseV1ValidationError) ErrorName() string {
	return "CreatePackageResponseV1ValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePackageResponseV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePackageResponseV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePackageResponseV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePackageResponseV1ValidationError{}

// Validate checks the field values on DescribePackageV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DescribePackageV1Request) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DescribePackageV1Request with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DescribePackageV1RequestMultiError, or nil if none found.
func (m *DescribePackageV1Request) ValidateAll() error {
	return m.validate(true)
}

func (m *DescribePackageV1Request) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPackageId() <= 0 {
		err := DescribePackageV1RequestValidationError{
			field:  "PackageId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DescribePackageV1RequestMultiError(errors)
	}

	return nil
}

// DescribePackageV1RequestMultiError is an error wrapping multiple validation
// errors returned by DescribePackageV1Request.ValidateAll() if the designated
// constraints aren't met.
type DescribePackageV1RequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DescribePackageV1RequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DescribePackageV1RequestMultiError) AllErrors() []error { return m }

// DescribePackageV1RequestValidationError is the validation error returned by
// DescribePackageV1Request.Validate if the designated constraints aren't met.
type DescribePackageV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribePackageV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribePackageV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribePackageV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribePackageV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribePackageV1RequestValidationError) ErrorName() string {
	return "DescribePackageV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e DescribePackageV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribePackageV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribePackageV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribePackageV1RequestValidationError{}

// Validate checks the field values on DescribePackageV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DescribePackageV1Response) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DescribePackageV1Response with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DescribePackageV1ResponseMultiError, or nil if none found.
func (m *DescribePackageV1Response) ValidateAll() error {
	return m.validate(true)
}

func (m *DescribePackageV1Response) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetValue()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DescribePackageV1ResponseValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DescribePackageV1ResponseValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DescribePackageV1ResponseValidationError{
				field:  "Value",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return DescribePackageV1ResponseMultiError(errors)
	}

	return nil
}

// DescribePackageV1ResponseMultiError is an error wrapping multiple validation
// errors returned by DescribePackageV1Response.ValidateAll() if the
// designated constraints aren't met.
type DescribePackageV1ResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DescribePackageV1ResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DescribePackageV1ResponseMultiError) AllErrors() []error { return m }

// DescribePackageV1ResponseValidationError is the validation error returned by
// DescribePackageV1Response.Validate if the designated constraints aren't met.
type DescribePackageV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribePackageV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribePackageV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribePackageV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribePackageV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribePackageV1ResponseValidationError) ErrorName() string {
	return "DescribePackageV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DescribePackageV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribePackageV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribePackageV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribePackageV1ResponseValidationError{}

// Validate checks the field values on ListPackagesV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListPackagesV1Request) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListPackagesV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListPackagesV1RequestMultiError, or nil if none found.
func (m *ListPackagesV1Request) ValidateAll() error {
	return m.validate(true)
}

func (m *ListPackagesV1Request) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetOffset() <= 0 {
		err := ListPackagesV1RequestValidationError{
			field:  "Offset",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetLimit() <= 0 {
		err := ListPackagesV1RequestValidationError{
			field:  "Limit",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ListPackagesV1RequestMultiError(errors)
	}

	return nil
}

// ListPackagesV1RequestMultiError is an error wrapping multiple validation
// errors returned by ListPackagesV1Request.ValidateAll() if the designated
// constraints aren't met.
type ListPackagesV1RequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListPackagesV1RequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListPackagesV1RequestMultiError) AllErrors() []error { return m }

// ListPackagesV1RequestValidationError is the validation error returned by
// ListPackagesV1Request.Validate if the designated constraints aren't met.
type ListPackagesV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPackagesV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPackagesV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPackagesV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPackagesV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPackagesV1RequestValidationError) ErrorName() string {
	return "ListPackagesV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListPackagesV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPackagesV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPackagesV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPackagesV1RequestValidationError{}

// Validate checks the field values on ListPackagesV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListPackagesV1Response) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListPackagesV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListPackagesV1ResponseMultiError, or nil if none found.
func (m *ListPackagesV1Response) ValidateAll() error {
	return m.validate(true)
}

func (m *ListPackagesV1Response) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPackages() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListPackagesV1ResponseValidationError{
						field:  fmt.Sprintf("Packages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListPackagesV1ResponseValidationError{
						field:  fmt.Sprintf("Packages[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListPackagesV1ResponseValidationError{
					field:  fmt.Sprintf("Packages[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListPackagesV1ResponseMultiError(errors)
	}

	return nil
}

// ListPackagesV1ResponseMultiError is an error wrapping multiple validation
// errors returned by ListPackagesV1Response.ValidateAll() if the designated
// constraints aren't met.
type ListPackagesV1ResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListPackagesV1ResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListPackagesV1ResponseMultiError) AllErrors() []error { return m }

// ListPackagesV1ResponseValidationError is the validation error returned by
// ListPackagesV1Response.Validate if the designated constraints aren't met.
type ListPackagesV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPackagesV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPackagesV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPackagesV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPackagesV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPackagesV1ResponseValidationError) ErrorName() string {
	return "ListPackagesV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListPackagesV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPackagesV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPackagesV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPackagesV1ResponseValidationError{}

// Validate checks the field values on RemovePackageV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RemovePackageV1Request) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RemovePackageV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RemovePackageV1RequestMultiError, or nil if none found.
func (m *RemovePackageV1Request) ValidateAll() error {
	return m.validate(true)
}

func (m *RemovePackageV1Request) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetPackageId() <= 0 {
		err := RemovePackageV1RequestValidationError{
			field:  "PackageId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return RemovePackageV1RequestMultiError(errors)
	}

	return nil
}

// RemovePackageV1RequestMultiError is an error wrapping multiple validation
// errors returned by RemovePackageV1Request.ValidateAll() if the designated
// constraints aren't met.
type RemovePackageV1RequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RemovePackageV1RequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RemovePackageV1RequestMultiError) AllErrors() []error { return m }

// RemovePackageV1RequestValidationError is the validation error returned by
// RemovePackageV1Request.Validate if the designated constraints aren't met.
type RemovePackageV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemovePackageV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemovePackageV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemovePackageV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemovePackageV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemovePackageV1RequestValidationError) ErrorName() string {
	return "RemovePackageV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e RemovePackageV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemovePackageV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemovePackageV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemovePackageV1RequestValidationError{}

// Validate checks the field values on RemovePackageV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RemovePackageV1Response) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RemovePackageV1Response with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RemovePackageV1ResponseMultiError, or nil if none found.
func (m *RemovePackageV1Response) ValidateAll() error {
	return m.validate(true)
}

func (m *RemovePackageV1Response) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Removed

	if len(errors) > 0 {
		return RemovePackageV1ResponseMultiError(errors)
	}

	return nil
}

// RemovePackageV1ResponseMultiError is an error wrapping multiple validation
// errors returned by RemovePackageV1Response.ValidateAll() if the designated
// constraints aren't met.
type RemovePackageV1ResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RemovePackageV1ResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RemovePackageV1ResponseMultiError) AllErrors() []error { return m }

// RemovePackageV1ResponseValidationError is the validation error returned by
// RemovePackageV1Response.Validate if the designated constraints aren't met.
type RemovePackageV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemovePackageV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemovePackageV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemovePackageV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemovePackageV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemovePackageV1ResponseValidationError) ErrorName() string {
	return "RemovePackageV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RemovePackageV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemovePackageV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemovePackageV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemovePackageV1ResponseValidationError{}
