package goliburing

/*
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include "liburing.h"

struct io_data {
	off_t  offset;
	size_t len;
	struct iovec iov;
};

struct io_data data;
struct io_uring_cqe *cqe;
struct io_uring_sqe *sqe;
int write_and_submit_and_wait(struct io_uring *ring, void* buf, size_t len, off_t offset, int fd) {
	data.offset = offset;
	data.len = len;
	data.iov.iov_base = buf;
	data.iov.iov_len = len;

	sqe = io_uring_get_sqe(ring);
	if (!sqe) {
    return -2; // Submission queue full.
	}

	io_uring_prep_writev(sqe, fd, &data.iov, 1, data.offset);

	// io_uring_prep_rw(IORING_OP_WRITEV, sqe, fd, &data->iov, 1, data->offset);
	// sqe->opcode = IORING_OP_WRITEV;
	// sqe->flags = 0;
	// sqe->ioprio = 0;
	// sqe->fd = fd;
	// sqe->off = data->offset;
	// sqe->addr = (unsigned long)  &data->iov;
	// sqe->len = 1;
	// sqe->rw_flags = 0;
	// sqe->user_data = (unsigned long) data;
	// sqe->__pad2[0] = sqe->__pad2[1] = sqe->__pad2[2] = 0;

	io_uring_sqe_set_data(sqe, &data);
	io_uring_submit(ring);

	while(1) {
		int ret = io_uring_peek_cqe(ring, &cqe);
		if (ret == 0) {
			io_uring_cq_advance(ring, 1);
			break;
		}
	}

	return cqe->res;
}
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// FileWriterSync synchronous file writer using an io_uring.
type FileWriterSync struct {
	ring   *Ring
	file   *os.File
	offset int
}

// NewFileWriterSync Create a new FileWriterSync.
func NewFileWriterSync(ring *Ring, file *os.File) (*FileWriterSync, error) {
	if ring == nil {
		return nil, fmt.Errorf("ring not provided")
	}

	return &FileWriterSync{
		ring:   ring,
		file:   file,
		offset: 0,
	}, nil
}

// Write the data to the file. Returns number of bytes written or an error.
func (r *FileWriterSync) Write(data []byte) (int, error) {
	p := unsafe.Pointer(&data[0])
	l := C.size_t(len(data))
	ret := C.write_and_submit_and_wait(r.ring.ring, p, l, C.off_t(r.offset), C.int(r.file.Fd()))
	if ret >= 0 {
		return int(ret), nil
	}

	return 0, fmt.Errorf("Write failed %d", ret)
}
