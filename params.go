package goliburing

/*
#include <stdlib.h>
#include <string.h>
#include "liburing.h"

// Represents the result of creating an `io_uring_params`.
// Ret values:
//  = 0  	success
//  = 1 	indicates a malloc failure.
struct io_uring_params_create_result {
	int ret;
	struct io_uring_params* params;
};

// Allocates an io_uring_params on the heap and initializes parameters.
struct io_uring_params_create_result create_params(unsigned flags, unsigned sq_thread_cpu,
	unsigned sq_thread_idle, unsigned cq_entries) {
	struct io_uring_params_create_result res;
	res.params = malloc(sizeof(*res.params));
	if (!res.params) {
		res.ret = 1;
    return res;
	}

	memset(res.params, 0, sizeof(*res.params));

	res.params->flags = flags;
	res.params->sq_thread_cpu = sq_thread_cpu;
	res.params->sq_thread_idle = sq_thread_idle;
	res.params->cq_entries = cq_entries;
	return res;
}

void destroy_params(struct io_uring_params* params) {
	if (params != NULL) {
		free(params);
		params = NULL;
	}
}
*/
import "C"
import "unsafe"

// Params Wrapper around io_uring_params.
type Params struct {
	params *C.struct_io_uring_params
}

// NewParams Create new params.
// If IORingSetupSQPoll is set, sqThreadCPU and sqThreadIdle are used.
// If IORingSetupCQSize is set, cqEntries is used.
// See documentation on IORingSetupSQPoll for usage.
// Error will nil on success and of type `NewErrCreateParams` on failure.
func NewParams(flags SetupFlag, sqThreadCPU uint32, sqThreadIdle uint32, cqEntries uint32) (*Params, error) {
	// If IORingSetupSQPoll is not set, make sure the controlling
	// variables are set to 0.
	if flags&IORingSetupSQPoll == 0 {
		sqThreadCPU = 0
		sqThreadIdle = 0
	}

	// If IORingSetupCQSize is not set, make sure the controlling
	// variables are set to 0.
	if flags&IORingSetupCQSize == 0 {
		cqEntries = 0
	}

	res := C.create_params(C.uint(flags), C.uint(sqThreadCPU), C.uint(sqThreadIdle), C.uint(cqEntries))
	if res.ret != 0 {
		return nil, NewErrCreateParams(res.ret)
	}

	return &Params{
		params: res.params,
	}, nil
}

// NewDefaultParams Create new params with default values.
// Error will nil on success and of type `NewErrCreateParams` on failure.
func NewDefaultParams() (*Params, error) {
	return NewParams(IORingSetupDefault, 0, 0, 0)
}

// Destroy Destroy the params (free memory).
func (p *Params) Destroy() {
	C.destroy_params(p.params)
}

// Flags Gets the parameters flags. See `SetupFlag`
func (p *Params) Flags() SetupFlag {
	if unsafe.Pointer(p.params) == unsafe.Pointer(C.NULL) {
		return 0
	}

	return SetupFlag(p.params.flags)
}

// Features Gets the parameters features. See `FeatureFlag`.
func (p *Params) Features() FeatureFlag {
	if unsafe.Pointer(p.params) == unsafe.Pointer(C.NULL) {
		return 0
	}

	return FeatureFlag(p.params.features)
}
