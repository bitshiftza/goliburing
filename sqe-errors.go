package goliburing

import "fmt"

// ErrGetSQECode Error codes returned when fetching an empty SQE.
type ErrGetSQECode = int

const (
	// ErrSQEMalloc Could not allocate memory for SQE.
	ErrSQEMalloc ErrGetSQECode = 1
	// ErrSQEFull Submission queue is full.
	ErrSQEFull = 2
)

// ErrGetSQE Represents an error while fetching an empty SQE.
type ErrGetSQE struct {
	Code    ErrGetSQECode
	Message string
}

// NewErrGetSQE Creates a new get SQE error.
func NewErrGetSQE(code ErrGetSQECode) *ErrGetSQE {
	var message string
	if code == ErrSQEMalloc {
		message = "could not allocate memory"
	} else if code == ErrSQEFull {
		message = "submission queue is full"
	} else {
		message = fmt.Sprintf("unknown error code %d", code)
	}

	return &ErrGetSQE{
		Code:    code,
		Message: message,
	}
}

func (e *ErrGetSQE) Error() string {
	return fmt.Sprintf("submission queue entry: %s", e.Message)
}
