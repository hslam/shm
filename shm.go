// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package shm provides a way to use shared memory.
package shm

import (
	"os"
)

// validSize returns the valid size.
func validSize(size int64) int64 {
	pageSize := int64(os.Getpagesize())
	if size%pageSize == 0 {
		return size
	}
	return (size/pageSize + 1) * pageSize
}
