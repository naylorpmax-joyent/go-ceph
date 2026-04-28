//go:build ceph_preview

package rados

// #cgo LDFLAGS: -lrados
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
//
// #if __APPLE__
// #define ceph_time_t __darwin_time_t
// #define ceph_suseconds_t __darwin_suseconds_t
// #elif __GLIBC__
// #define ceph_time_t __time_t
// #define ceph_suseconds_t __suseconds_t
// #else
// #define ceph_time_t time_t
// #define ceph_suseconds_t suseconds_t
// #endif
import "C"

import "unsafe"

// Checksum calculates the checksum of the given object data, using one of the supported checksum algorithms.
//
// Implements:
//
//	int rados_checksum(rados_ioctx_t io,
//	                   const char *oid,
//	                   rados_checksum_type_t type,
//	                   const char *init_value,
//	                   size_t init_value_len,
//	                   size_t len,
//	                   uint64_t off,
//	                   size_t chunk_size,
//	                   char *pchecksum,
//	                   size_t checksum_len);
func (ioctx *IOContext) Checksum(oid string, checksumType ChecksumType, dst []byte, opts ...func(*ChecksumOptions)) error {
	// apply defaults and options
	copts := &ChecksumOptions{initValue: make([]byte, 8)}
	if checksumType == ChecksumTypeXXHash64 {
		copts.initValue = make([]byte, 16)
	}
	for _, o := range opts {
		o(copts)
	}

	// call library
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	return getError(C.rados_checksum(
		ioctx.ioctx,
		coid,
		C.rados_checksum_type_t(checksumType),
		(*C.char)(unsafe.Pointer(&copts.initValue[0])),
		C.size_t(len(copts.initValue)),
		C.size_t(copts.len),
		C.uint64_t(copts.off),
		C.size_t(copts.chunkSize),
		(*C.char)(unsafe.Pointer(&dst[0])),
		C.size_t(len(dst)),
	))
}

// ChecksumType indicates checksum algorithm types supported by the IOContext.Checksum method.
// Equivalent to the rados_checksum_type_t enum.
type ChecksumType uint32

const (
	// ChecksumTypeXXHash32 produces an encoded le32 checksum of the given object.
	ChecksumTypeXXHash32 = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_XXHASH32)
	// ChecksumTypeXXHash64 produces an encoded le64 checksum of the given object.
	ChecksumTypeXXHash64 = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_XXHASH64)
	// ChecksumTypeCRC32C produces an encoded le32 checksum of the given object.
	ChecksumTypeCRC32C = ChecksumType(C.LIBRADOS_CHECKSUM_TYPE_CRC32C)
)

// ChecksumOptions exposes non-required parameters for the Checksum method.
type ChecksumOptions struct {
	off       uint64
	len       uint64
	chunkSize uint64
	initValue []byte
}

// ChecksumOff sets the object offset to start checksumming in the object.
// By default, the entire object will be checksummed.
func ChecksumOff(v uint64) func(*ChecksumOptions) {
	return func(copts *ChecksumOptions) {
		copts.off = v
	}
}

// ChecksumLen sets the the number of bytes to checksum in the object.
// By default, the entire object will be checksummed.
func ChecksumLen(v uint64) func(*ChecksumOptions) {
	return func(copts *ChecksumOptions) {
		copts.len = v
	}
}

// ChecksumChunkSize sets the length-aligned chunk size for the checksum calculation.
// By default, the entire object will be checksummed as a single chunk.
func ChecksumChunkSize(v uint64) func(*ChecksumOptions) {
	return func(copts *ChecksumOptions) {
		copts.chunkSize = v
	}
}

// ChecksumInitValue sets the initial value for the checksum calculation.
func ChecksumInitValue(v []byte) func(*ChecksumOptions) {
	return func(copts *ChecksumOptions) {
		copts.initValue = v
	}
}
