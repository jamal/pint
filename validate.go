package pint

import (
	"fmt"
	"strconv"
)

// ErrValidate is a custom error type returned by validation functions, the
// Error() format should only be used internally. Use ErrValidate.String() for
// an error that can be returned to the client.
type ErrValidate struct {
	s string
}

func (e *ErrValidate) Error() string {
	return "pint.Validate: " + e.s
}

func (e *ErrValidate) String() string {
	return e.s
}

func validateInt(val int64, field field) error {
	if minVal, ok := field.options.get("min"); ok {
		min, err := strconv.ParseInt(minVal, 10, 64)
		if err != nil {
			return err
		}
		if val < min {
			return &ErrValidate{fmt.Sprintf("%s must be greater than %d", field.name, min)}
		}
	}
	if maxVal, ok := field.options.get("max"); ok {
		max, err := strconv.ParseInt(maxVal, 10, 64)
		if err != nil {
			return err
		}
		if val > max {
			return &ErrValidate{fmt.Sprintf("%s must be less than %d", field.name, max)}
		}
	}
	return nil
}
