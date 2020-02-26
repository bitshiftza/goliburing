package goliburing

/*
#include "liburing.h"
*/
import "C"

// SetupFlag Setup flags for io_uring.
type SetupFlag = uint32

const (
	// IORingSetupDefault no specififc flag
	IORingSetupDefault SetupFlag = 0
	// IORingSetupIOPoll io_context is polled
	// Perform busy-waiting for an I/O completion, as opposed to getting notifications via an asynchronous IRQ
	// (Interrupt Request). The file system (if any) and block device must support polling in order for this
	// to work. Busy-waiting provides lower latency, but may consume more CPU resources than interrupt driven
	// I/O. Currently, this feature is usable only on a file descriptor opened using the O_DIRECT flag. When
	// a read or write is submitted to a polled context, the application must poll for completions on the CQ
	// ring by calling io_uring_enter(2). It is illegal to mix and match polled and non-polled I/O on an
	// io_uring instance.
	IORingSetupIOPoll = C.IORING_SETUP_IOPOLL
	// IORingSetupSQPoll SQ poll thread
	// When this flag is specified, a kernel thread is created to perform submission queue polling. An io_ur‐
	// ing instance configured in this way enables an application to issue I/O without ever context switching
	// into the kernel. By using the submission queue to fill in new submission queue entries and watching
	// for completions on the completion queue, the application can submit and reap I/Os without doing a sin‐
	// gle system call.
	// If the kernel thread is idle for more than sq_thread_idle milliseconds, it will set the IOR‐
	// ING_SQ_NEED_WAKEUP bit in the flags field of the struct io_sq_ring. When this happens, the application
	// must call io_uring_enter(2) to wake the kernel thread. If I/O is kept busy, the kernel thread will
	// never sleep. An application making use of this feature will need to guard the io_uring_enter(2) call
	// with the following code sequence:
	// 		/*
	// 			* Ensure that the wakeup flag is read after the tail pointer has been
	// 			* written.
	// 			*/
	// 		smp_mb();
	// 		if (*sq_ring->flags & IORING_SQ_NEED_WAKEUP)
	// 				io_uring_enter(fd, 0, 0, IORING_ENTER_SQ_WAKEUP);
	// where sq_ring is a submission queue ring setup using the struct io_sqring_offsets described below.
	// To successfully use this feature, the application must register a set of files to be used for IO
	// through io_uring_register(2) using the IORING_REGISTER_FILES opcode. Failure to do so will result in
	// submitted IO being errored with EBADF.
	IORingSetupSQPoll = C.IORING_SETUP_SQPOLL
	// IORingSetupSQAff sq_thread_cpu is valid
	// If this flag is specified, then the poll thread will be bound to the cpu set in the sq_thread_cpu field
	// of the struct io_uring_params. This flag is only meaningful when IORING_SETUP_SQPOLL is specified.
	IORingSetupSQAff = C.IORING_SETUP_SQ_AFF
	// IORingSetupCQSize App defines CQ size
	// Create the completion queue with struct io_uring_params.cq_entries entries. The value must be greater
	// than entries, and may be rounded up to the next power-of-two.
	IORingSetupCQSize = C.IORING_SETUP_CQSIZE
	// IORingSetupClamp Clamp SQ/CQ ring sizes
	IORingSetupClamp = C.IORING_SETUP_CLAMP
	// IORingSetupAttachWQ Attach to existing wq
	IORingSetupAttachWQ = C.IORING_SETUP_ATTACH_WQ
)

// FeatureFlag Feature flags for io_uring.
// io_uring_params->features is filled in by the kernel, which
// specifies various features supported by current kernel version.
type FeatureFlag = uint32

const (
	// IORingFeatSingleMMap If this flag is set, the two SQ and CQ rings can be mapped with a single mmap(2) call.
	// The SQEs must still be allocated separately. This brings the necessary mmap(2) calls down from three to two.
	IORingFeatSingleMMap = C.IORING_FEAT_SINGLE_MMAP
	// IORingFeatNoDrop If this flag is set, io_uring supports never dropping completion events.
	// If a completion event occurs and the CQ ring is full, the kernel stores the event internally until
	// such a time that the CQ ring has room for more entries. If this overflow condition is entered, attempting
	// to submit more IO with fail with the -EBUSY error value, if it can't flush the overflown events
	// to the CQ ring. If this happens, the application must reap events from the CQ ring and attempt the submit again.
	IORingFeatNoDrop = C.IORING_FEAT_NODROP
	// IORingFeatSubmitStable If this flag is set, applications can be certain that any data for async offload has
	// been consumed when the kernel has consumed the SQE.
	IORingFeatSubmitStable = C.IORING_FEAT_SUBMIT_STABLE
	// IORingFeatRWCurPos If this flag is set, applications can specify offset == -1 with IORING_OP_{READV,WRITEV} , IOR‐
	// ING_OP_{READ,WRITE}_FIXED , and IORING_OP_{READ,WRITE} to mean current file position, which behaves
	// like preadv2(2) and pwritev2(2) with offset == -1. It'll use (and update) the current file position.
	// This obviously comes with the caveat that if the application has multiple reads or writes in flight,
	// then the end result will not be as expected. This is similar to threads sharing a file descriptor and
	// doing IO using the current file position.
	IORingFeatRWCurPos = C.IORING_FEAT_RW_CUR_POS
	// IORingFeatCurPersonality If this flag is set, then io_uring guarantees that both sync and async execution of a
	// request assumes the credentials of the task that called io_uring_enter(2) to queue the requests. If this flag isn't
	// set, then requests are issued with the credentials of the task that originally registered the io_uring.
	// If only one task is using a ring, then this flag doesn't matter as the credentials will always be the
	// same. Note that this is the default behavior, tasks can still register different personalities through
	// io_uring_register(2) with IORING_REGISTER_PERSONALITY and specify the personality to use in the sqe.
	IORingFeatCurPersonality = C.IORING_FEAT_CUR_PERSONALITY
)

// OpFlag operation flag.
type OpFlag = uint32

const (
	// IORingOpNOP Do not perform any I/O. This is useful for testing the performance of the io_uring implementation itself.
	IORingOpNOP = C.IORING_OP_NOP
	// IORingOpReadV Vectored read operation, similar to preadv2(2).
	IORingOpReadV = C.IORING_OP_READV
	// IORingOpWriteV Vectored write operation, similar to pwritev2(2).
	IORingOpWriteV = C.IORING_OP_WRITEV
	// IORingOpFSync File sync.  See also fsync(2). Note that, while I/O is initiated in the order in which it appears in the submission queue, completions are unordered.
	// For example, an application which places a write I/O followed by an fsync in the submission queue cannot expect the fsync to apply to the write. The two
	// operations execute in parallel, so the fsync may complete before the write is issued to the storage. The same is also true for previously issued writes
	// that have not completed prior to the fsync.
	IORingOpFSync = C.IORING_OP_FSYNC
	// IORingOpReadFixed Read from pre-mapped buffers. See io_uring_register(2) for details on how to setup a context for fixed reads.
	IORingOpReadFixed = C.IORING_OP_READ_FIXED
	// IORingOpWriteFixed Write to pre-mapped buffers. See io_uring_register(2) for details on how to setup a context for fixed writes.
	IORingOpWriteFixed = C.IORING_OP_WRITE_FIXED
	// IORingOpPollAdd Poll the fd specified in the submission queue entry for the events specified in the poll_events field.  Unlike poll or epoll  without  EPOLLONESHOT,  this
	// interface always works in one shot mode.  That is, once the poll operation is completed, it will have to be resubmitted.
	IORingOpPollAdd = C.IORING_OP_POLL_ADD
	// IORingOpPolLRemove Remove an existing poll request.  If found, the res field of the struct io_uring_cqe will contain 0.  If not found, res will contain -ENOENT.
	IORingOpPolLRemove = C.IORING_OP_POLL_REMOVE
	// IORingOpSyncFileRange Issue  the  equivalent  of  a  sync_file_range  (2) on the file descriptor. The fd field is the file descriptor to sync, the off field holds the offset in
	// bytes, the len field holds the length in bytes, and the flags field holds the flags for the command. See also sync_file_range(2).   for  the  general  de‐
	// scription of the related system call. Available since 5.2.
	IORingOpSyncFileRange = C.IORING_OP_SYNC_FILE_RANGE
	// IORingOpSendMsg Issue  the equivalent of a sendmsg(2) system call.  fd must be set to the socket file descriptor, addr must contain a pointer to the msghdr structure, and
	// flags holds the flags associated with the system call. See also sendmsg(2).  for the general description of the related system call. Available since 5.3.
	IORingOpSendMsg = C.IORING_OP_SENDMSG
	// IORingOpRecvMsg Works just like IORING_OP_SENDMSG, except for recvmsg(2) instead. See the description of IORING_OP_SENDMSG. Available since 5.3.
	IORingOpRecvMsg = C.IORING_OP_RECVMSG
	// IORingOpTimeout This command will register a timeout operation. The addr field must contain a pointer to a struct timespec64 structure, len must contain 1 to signify  one
	// timespec64 structure, timeout_flags may contain IORING_TIMEOUT_ABS for an absolutel timeout value, or 0 for a relative timeout.  off may contain a comple‐
	// tion event count. If not set, this defaults to 1. A timeout will trigger a wakeup event on the completion ring for anyone waiting for  events.  A  timeout
	// condition  is  met  when  either the specified timeout expires, or the specified number of events have completed. Either condition will trigger the event.
	// io_uring timeouts use the CLOCK_MONOTONIC clock source. The request will complete with -ETIME if the timeout  got  completed  through  expiration  of  the
	// timer,  or  0 if the timeout got completed through requests completing on their own. If the timeout was cancelled before it expired, the request will com‐
	// plete with -ECANCELED.  Available since 5.4.
	IORingOpTimeout = C.IORING_OP_TIMEOUT
	// IORingOpTimeoutRemove Attempt to remove an existing timeout operation.  addr must contain the user_data field of the previously issued timeout operation. If the specified time‐
	// out  request is found and cancelled successfully, this request will terminate with a result value of 0 If the timeout request was found but expiration was
	// already in progress, this request will terminate with a result value of -EBUSY If the timeout request wasn't found, the request will terminate with a  re‐
	// sult value of -ENOENT Available since 5.5.
	IORingOpTimeoutRemove = C.IORING_OP_TIMEOUT_REMOVE
	// IORingOpAccept Issue the equivalent of an accept4(2) system call.  fd must be set to the socket file descriptor, addr must contain the pointer to the sockaddr structure,
	// and addr2 must contain a pointer to the socklen_t addrlen field. See also accept4(2) for the general description of the  related  system  call.  Available
	// since 5.5.
	IORingOpAccept = C.IORING_OP_ACCEPT
	// IORingOpAsyncCancel Attempt to cancel an already issued request.  addr must contain the user_data field of the request that should be cancelled. The cancellation request will
	// complete with one of the following results codes. I found, the res field of the cqe will contain 0. If not found, res will contain -ENOENT. If  found  and
	// attempted cancelled, the res field will contain -EALREADY. In this case, the request may or may not terminate. In general, requests that are interruptible
	// (like socket IO) will get cancelled, while disk IO requests cannot be cancelled if already started.  Available since 5.5.
	IORingOpAsyncCancel = C.IORING_OP_ASYNC_CANCEL
	// IORingOpLinkTimeout This request must be linked with another request through IOSQE_IO_LINK which is described below. Unlike IORING_OP_TIMEOUT, IORING_OP_LINK_TIMEOUT acts  on
	// the  linked  request, not the completion queue. The format of the command is otherwise like IORING_OP_TIMEOUT, except there's no completion event count as
	// it's tied to a specific request.  If used, the timeout specified in the command will cancel the linked command, unless the linked command completes before
	// the  timeout.  The  timeout  will complete with -ETIME if the timer expired and the linked request was attempted cancelled, or -ECANCELED if the timer got
	// cancelled because of completion of the linked request. Like IORING_OP_TIMEOUT the clock source used is CLOCK_MONOTONIC Available since 5.5.
	IORingOpLinkTimeout = C.IORING_OP_LINK_TIMEOUT
	// IORingOpConnect Issue the equivalent of a connect(2) system call.  fd must be set to the socket file descriptor, addr must contain the pointer to the sockaddr  structure,
	// and off must contain the socklen_t addrlen field. See also connect(2) for the general description of the related system call. Available since 5.5.
	IORingOpConnect = C.IORING_OP_CONNECT
	// IORingOpFAllocate Issue  the equivalent of a fallocate(2) system call.  fd must be set to the file descriptor, off must contain the offset on which to operate, and len must
	// contain the length. See also fallocate(2) for the general description of the related system call. Available since 5.6.
	IORingOpFAllocate = C.IORING_OP_FALLOCATE
	// IORingOpOpenAt Issue the equivalent of a openat(2) system call.  fd is the dirfd argument, addr must contain a pointer to the *pathname argument, open_flags should  con‐
	// tain  any flags passed in, and mode is access mode of the file. See also openat(2) for the general description of the related system call. Available since
	// 5.6.
	IORingOpOpenAt = C.IORING_OP_OPENAT
	// IORingOpCLOSE Issue the equivalent of a close(2) system call.  fd is the file descriptor to be closed. See also close(2) for the general description of the related sys‐
	// tem call. Available since 5.6.
	IORingOpCLOSE = C.IORING_OP_CLOSE
	// IORingOpFilesUpdate This  command  is  an alternative to using IORING_REGISTER_FILES_UPDATE which then works in an async fashion, like the rest of the io_uring commands.  The
	// arguments passed in are the same.  addr must contain a pointer to the array of file descriptors, len must contain the length of the array,  and  off  must
	// contain  the  offset at which to operate. Note that the array of file descriptors pointed to in addr must remain valid until this operation has completed.
	// Available since 5.6.
	IORingOpFilesUpdate = C.IORING_OP_FILES_UPDATE
	// IORingOpStatX Issue  the  equivalent of a statx(2) system call.  fd is the dirfd argument, addr must contain a pointer to the *pathname string, statx_flags is the flags
	// argument, len should be the mask argument, and off must contain a pointer to the statxbuf to be filled in. See also statx(2) for the  general  description
	// of the related system call. Available since 5.6.
	IORingOpStatX = C.IORING_OP_STATX
	// IORingOpRead Issue  the  equivalent  of  a read(2) or write(2) system call.  fd is the file descriptor to be operated on, addr contains the buffer in question, and len
	// contains the length of the IO operation. These are non-vectored versions of the  IORING_OP_READV  and  IORING_OP_WRITEV  opcodes.  See  also  read(2)  and
	// write(2) for the general description of the related system call. Available since 5.6.
	IORingOpRead = C.IORING_OP_READ
	// IORingOpWrite Issue  the  equivalent  of  a read(2) or write(2) system call.  fd is the file descriptor to be operated on, addr contains the buffer in question, and len
	// contains the length of the IO operation. These are non-vectored versions of the  IORING_OP_READV  and  IORING_OP_WRITEV  opcodes.  See  also  read(2)  and
	// write(2) for the general description of the related system call. Available since 5.6.
	IORingOpWrite = C.IORING_OP_WRITE
	// IORingOpFAdvise Issue the equivalent of a posix_fadvise(2) system call.  fd must be set to the file descriptor, off must contain the offset on which to operate, len  must
	// contain the length, and fadvise_advice must contain the advice associated with the operation. See also posix_fadvise(2) for the general description of the
	// related system call. Available since 5.6.
	IORingOpFAdvise = C.IORING_OP_FADVISE
	// IORingOpMAdvise Issue the equivalent of a madvise(2) system call.  addr must contain the address to operate on, len must contain the length on which to operate, and  fad‐
	// vise_advice  must  contain the advice associated with the operation. See also madvise(2) for the general description of the related system call. Available
	// since 5.6.
	IORingOpMAdvise = C.IORING_OP_MADVISE
	// IORingOpSend todo
	IORingOpSend = C.IORING_OP_SEND
	// IORingOpRecv todo
	IORingOpRecv = C.IORING_OP_RECV
	// IORingOpOpenAt2 todo
	IORingOpOpenAt2 = C.IORING_OP_OPENAT2
	// IORingOpEPollCtl todo
	IORingOpEPollCtl = C.IORING_OP_EPOLL_CTL
)
