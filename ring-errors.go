package goliburing

/*
#include <string.h>
#include <errno.h>
*/
import "C"

import "fmt"

// ErrCreateRingCode Error codes returned from NewRing
type ErrCreateRingCode = int

const (
	// ErrCreateRingEPARAMS Could not create Params.
	ErrCreateRingEPARAMS ErrCreateRingCode = 0
	// ErrCreateRingEMALLOC Could not allocate memory for ring.
	ErrCreateRingEMALLOC = 0
	// ErrCreateRingEFAULT params is outside your accessible address space.
	ErrCreateRingEFAULT = C.EFAULT
	// ErrCreateRingEINVAL The resv array contains non-zero data, p.flags contains an unsupported flag,
	// entries is out of bounds, IORING_SETUP_SQ_AFF was specified, but IORING_SETUP_SQPOLL
	// was not, or IORING_SETUP_CQSIZE was specified, but io_uring_params.cq_entries was invalid.
	ErrCreateRingEINVAL = C.EINVAL
	// ErrCreateRingEMFILE The per-process limit on the number of open file descriptors has been reached
	// (see the description of RLIMIT_NOFILE in getrlimit(2)).
	ErrCreateRingEMFILE = C.EMFILE
	// ErrCreateRingENFILE The system-wide limit on the total number of open files has been reached.
	ErrCreateRingENFILE = C.ENFILE
	// ErrCreateRingENOMEM Insufficient kernel resources are available.
	ErrCreateRingENOMEM = C.ENOMEM
	// ErrCreateRingEPERM IORING_SETUP_SQPOLL was specified, but the effective user ID of the caller
	// did not have sufficient privileges.
	ErrCreateRingEPERM = C.EPERM
)

// ErrCreateRing Represents an error while creating an io_uring.
type ErrCreateRing struct {
	Code     ErrCreateRingCode
	Message  string
	SubError error
}

// NewErrCreateRing Creates a new create ring error.
func NewErrCreateRing(code C.int, subErr error) *ErrCreateRing {
	var message string
	if code == 0 {
		message = fmt.Sprintf("could not create params: %s", subErr.Error())
	} else if code == 1 {
		message = "could not allocate memory"
	} else {
		message = fmt.Sprintf("code: %d, err: %s", -code, C.GoString((C.strerror(-code))))
	}

	return &ErrCreateRing{
		Message:  message,
		Code:     ErrCreateRingCode(-code),
		SubError: subErr,
	}
}

func (e *ErrCreateRing) Error() string {
	return fmt.Sprintf("create ring failed: %s", e.Message)
}
