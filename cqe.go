package goliburing

/*
#include "liburing.h"
*/
import "C"

// CQE Represents a completion queue entry.
type CQE struct {
	ring *Ring
	cqe  *C.struct_io_uring_cqe
}

func newCQE(ring *Ring) *CQE {
	return &CQE{
		ring: ring,
	}
}

// Seen sets the CQE as seen.
func (c *CQE) Seen() {
	C.io_uring_cqe_seen(c.ring.ring, c.cqe)
}
