package goliburing

/*
#include <stdlib.h>
#include <string.h>
#include <sys/uio.h>
#include "liburing.h"

struct sqe {
	off_t offset;
	struct iovec *iov;
	struct io_uring_sqe *sqe;
};

struct sqe_create_result {
	struct sqe *s;
	int ret;
};

struct sqe_create_result get_sqe(struct io_uring *ring) {
	struct sqe_create_result res;
	res.s = malloc(sizeof(*res.s));
	if(!res.s) {
		res.ret = 1;
		return res;
	}

	res.s->iov = malloc(sizeof(*res.s->iov));
	if(!res.s->iov) {
		res.ret = 1;
		free(res.s);
		return res;
	}

	// res.s->sqe = io_uring_get_sqe(ring);
	// if(res.s->sqe == NULL) {
	// 	free(res.s);
	// 	res.ret = 2;
	// 	return res;
	// }

	res.ret = 0;
	return res;
}
*/
import "C"
import (
	"unsafe"
)

// SQE Represents a submission queue entry.
type SQE struct {
	sqe  *C.struct_sqe
	ring *Ring
}

func newSQE(ring *Ring) (*SQE, error) {
	res := C.get_sqe(ring.ring)
	if int(res.ret) != 0 {
		return nil, NewErrGetSQE(ErrGetSQECode(res.ret))
	}

	return &SQE{
		ring: ring,
		sqe:  res.s,
	}, nil
}

func (s *SQE) init() error {
	s.sqe.sqe = C.io_uring_get_sqe(s.ring.ring)
	if unsafe.Pointer(s.sqe.sqe) == unsafe.Pointer(C.NULL) {
		return NewErrGetSQE(ErrSQEFull)
	}
	return nil
}

// PrepWriteV Prepare a vectored write.
func (s *SQE) PrepWriteV(fd int, data []byte, offset uint64) {
	s.sqe.offset = C.off_t(offset)
	s.sqe.iov.iov_len = C.size_t(len(data))
	s.sqe.iov.iov_base = unsafe.Pointer(&data[0])

	C.io_uring_prep_rw(IORingOpWriteV, s.sqe.sqe,
		C.int(fd),
		unsafe.Pointer(s.sqe.iov),
		C.uint(1),
		C.ulonglong(offset))
}

// SetData Set the data for a submission queue event.
func (s *SQE) SetData(data []byte) {
	C.io_uring_sqe_set_data(s.sqe.sqe, unsafe.Pointer(&data[0]))
}
