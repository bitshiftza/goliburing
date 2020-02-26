package goliburing

/*
#include <string.h>
#include <errno.h>
*/
import "C"
import "fmt"

// ErrCreateParamsCode Error codes returned from NewParams
type ErrCreateParamsCode = int

const (
	// ErrCreateParamsEMALLOC Could not allocate memory for params.
	ErrCreateParamsEMALLOC ErrCreateParamsCode = 0
)

// ErrCreateParams Represents an error while creating io_uring_params.
type ErrCreateParams struct {
	Code    ErrCreateParamsCode
	Message string
}

// NewErrCreateParams Creates a new create params error.
func NewErrCreateParams(code C.int) *ErrCreateParams {
	var message string
	if code == 1 {
		message = "could not allocate memory"
	} else {
		message = "unknown error occurred."
	}

	return &ErrCreateParams{
		Code:    ErrCreateParamsCode(code),
		Message: message,
	}
}

func (e *ErrCreateParams) Error() string {
	return fmt.Sprintf("create params failed: %s", e.Message)
}
