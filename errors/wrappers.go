package errors

import "errors"

// These variables are used to give us access to existing
// functions in the stdlib errors pkg.
var (
	As = errors.As
	Is = errors.Is
)
