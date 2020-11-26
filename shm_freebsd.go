// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build freebsd

package shm

import (
	"syscall"
	"unsafe"
)

const (
	// SYS_SHMGET is syscall SYS_SHMGET constant
	SYS_SHMGET = 231
	// SYS_SHMAT is syscall SYS_SHMAT constant
	SYS_SHMAT = 228
	// SYS_SHMDT is syscall SYS_SHMDT constant
	SYS_SHMDT = 230
	// SYS_SHMCTL is syscall SYS_SHMCTL constant
	SYS_SHMCTL = 229
)

// Open returns the fd.
func Open(name string, oflag int, mode int) (int, error) {
	n, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, err
	}
	fd, _, errno := syscall.Syscall(syscall.SYS_SHM_OPEN, uintptr(unsafe.Pointer(n)), uintptr(oflag), uintptr(mode))
	if errno != 0 {
		return 0, errno
	}
	return int(fd), nil
}

// Unlink unlinks the name.
func Unlink(name string) error {
	n, err := syscall.BytePtrFromString(name)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall(syscall.SYS_SHM_UNLINK, uintptr(unsafe.Pointer(n)), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
