// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package shm provides a way to use System V shared memory.
package shm

import (
	"os"
	"syscall"
)

const (
	// O_RDONLY opens the file read-only.
	O_RDONLY int = syscall.O_RDONLY
	// O_WRONLY opens the file write-only.
	O_WRONLY int = syscall.O_WRONLY
	// O_RDWR opens the file read-write.
	O_RDWR int = syscall.O_RDWR

	// O_CREATE creates a new file if none exists.
	O_CREATE int = syscall.O_CREAT
)

// validSize returns the valid size.
func validSize(size int64) int64 {
	pageSize := int64(os.Getpagesize())
	if size%pageSize == 0 {
		return size
	}
	return (size/pageSize + 1) * pageSize
}
