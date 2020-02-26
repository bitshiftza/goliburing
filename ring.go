package goliburing

/*
#cgo pkg-config: liburing

#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include "liburing.h"

// Represents the result of creating an `io_uring`.
// Ret values:
//  = 0  	success
//  = 1 	indicates a malloc failure.
//  < 0 	failure from io_uring_queue_init.
struct io_uring_create_result {
	int ret;
	struct io_uring* ring;
};

// Allocates an io_uring on the heap and initializes the queue.
struct io_uring_create_result create_ring(unsigned queueDepth, struct io_uring_params* params) {
	struct io_uring_create_result res;
	res.ring = malloc(sizeof(*res.ring));
	if (!res.ring) {
		res.ret = 1;
    return res;
	}

	res.ret = io_uring_queue_init_params(queueDepth, res.ring, params);
	if (res.ret < 0) {
		free(res.ring);
	}

	return res;
}

// Destroys and frees an allocated io_uring.
void destroy_ring(struct io_uring* ring) {
	if(ring != NULL) {
		io_uring_queue_exit(ring);
		free(ring);
		ring = 0;
	}
}

// static void msec_to_ts(struct __kernel_timespec *ts, unsigned int msec)
// {
// 	ts->tv_sec = msec / 1000;
// 	ts->tv_nsec = (msec % 1000) * 1000000;
// }

// int wait_cqes(struct io_uring *ring, struct io_uring_cqe **cqe_ptr, unsigned wait_nr) {
// 	struct __kernel_timespec ts;
// 	msec_to_ts(&ts, 1000);

// 	return io_uring_wait_cqes(ring, cqe_ptr, wait_nr, &ts, 0);
// }
*/
import "C"
import "fmt"

// Ring Wrapper around io_uring.
type Ring struct {
	ring       *C.struct_io_uring
	Params     *Params
	queueDepth uint32
	sqes       []*SQE
	sqeIndex   uint32
	cqe        *CQE
}

// NewRing Create a new Ring.
// If params == nil, default params are used.
// If params != nil, they will be owned and destoyed by the Ring.
// Error will nil on success and of type `NewErrCreateRing` on failure.
func NewRing(queueDepth uint32, params *Params) (*Ring, error) {
	if params == nil {
		// Create default parameters.
		var err error
		params, err = NewDefaultParams()
		if err != nil {
			return nil, NewErrCreateRing(0, err)
		}
	}

	// Create ring.
	res := C.create_ring(C.uint(queueDepth), params.params)
	if res.ret != 0 {
		// Since the ring owns any params, destroy them.
		params.Destroy()
		return nil, NewErrCreateRing(res.ret, nil)
	}

	ring := &Ring{
		ring:       res.ring,
		Params:     params,
		queueDepth: queueDepth,
		sqes:       make([]*SQE, queueDepth),
		sqeIndex:   0,
	}
	ring.cqe = newCQE(ring)

	for i := 0; i < int(queueDepth); i++ {
		ring.sqes[i], _ = newSQE(ring)
	}

	return ring, nil
}

// QueueDepth Get the queue depth.
func (r *Ring) QueueDepth() uint32 {
	return r.queueDepth
}

// Destroy Destroy the ring.
func (r *Ring) Destroy() {
	C.destroy_ring(r.ring)
	r.Params.Destroy()
}

// GetEmptySQE Gets an empty submission queue entry.
func (r *Ring) GetEmptySQE() (*SQE, error) {
	sqe := r.sqes[r.sqeIndex]
	r.sqeIndex++
	if r.sqeIndex == r.queueDepth {
		r.sqeIndex = 0
	}
	err := sqe.init()
	return sqe, err
}

// Submit Submits SQEs.
func (r *Ring) Submit() {
	C.io_uring_submit(r.ring)
}

// WaitCQE Wait for a completion queue event.
func (r *Ring) WaitCQE() (*CQE, error) {
	ret := int(C.io_uring_wait_cqe(r.ring, &r.cqe.cqe))
	if ret != 0 {
		return nil, fmt.Errorf("WaitCQE failed with %d", ret)
	}

	return r.cqe, nil
}
