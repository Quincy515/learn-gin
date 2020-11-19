// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: Models.proto

package services

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _models_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on ProdModel with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *ProdModel) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ProdId

	// no validation rules for ProdName

	// no validation rules for ProdPrice

	return nil
}

// ProdModelValidationError is the validation error returned by
// ProdModel.Validate if the designated constraints aren't met.
type ProdModelValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ProdModelValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ProdModelValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ProdModelValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ProdModelValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ProdModelValidationError) ErrorName() string { return "ProdModelValidationError" }

// Error satisfies the builtin error interface
func (e ProdModelValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProdModel.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ProdModelValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ProdModelValidationError{}

// Validate checks the field values on OrderMain with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *OrderMain) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for OrderId

	// no validation rules for OrderNo

	// no validation rules for UserId

	if m.GetOrderMoney() <= 1 {
		return OrderMainValidationError{
			field:  "OrderMoney",
			reason: "value must be greater than 1",
		}
	}

	if v, ok := interface{}(m.GetOrderTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return OrderMainValidationError{
				field:  "OrderTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetOrderDetail() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrderMainValidationError{
					field:  fmt.Sprintf("OrderDetail[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// OrderMainValidationError is the validation error returned by
// OrderMain.Validate if the designated constraints aren't met.
type OrderMainValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderMainValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderMainValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderMainValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderMainValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderMainValidationError) ErrorName() string { return "OrderMainValidationError" }

// Error satisfies the builtin error interface
func (e OrderMainValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderMain.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderMainValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderMainValidationError{}

// Validate checks the field values on OrderDetail with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *OrderDetail) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for DetailId

	// no validation rules for OrderNo

	// no validation rules for ProdId

	// no validation rules for ProdPrice

	// no validation rules for ProdNum

	return nil
}

// OrderDetailValidationError is the validation error returned by
// OrderDetail.Validate if the designated constraints aren't met.
type OrderDetailValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderDetailValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderDetailValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderDetailValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderDetailValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderDetailValidationError) ErrorName() string { return "OrderDetailValidationError" }

// Error satisfies the builtin error interface
func (e OrderDetailValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderDetail.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderDetailValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderDetailValidationError{}

// Validate checks the field values on UserInfo with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *UserInfo) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for UserId

	// no validation rules for UserScore

	return nil
}

// UserInfoValidationError is the validation error returned by
// UserInfo.Validate if the designated constraints aren't met.
type UserInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserInfoValidationError) ErrorName() string { return "UserInfoValidationError" }

// Error satisfies the builtin error interface
func (e UserInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserInfoValidationError{}
