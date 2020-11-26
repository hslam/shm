// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build linux

package shm

import (
	"os"
	"syscall"
)

const (
	// SYS_SHMGET is syscall SYS_SHMGET constant
	SYS_SHMGET = syscall.SYS_SHMGET
	// SYS_SHMAT is syscall SYS_SHMAT constant
	SYS_SHMAT = syscall.SYS_SHMAT
	// SYS_SHMDT is syscall SYS_SHMDT constant
	SYS_SHMDT = syscall.SYS_SHMDT
	// SYS_SHMCTL is syscall SYS_SHMCTL constant
	SYS_SHMCTL = syscall.SYS_SHMCTL
)

// Open returns the fd.
func Open(name string, oflag int, perm int) (int, error) {
	return syscall.Open("/dev/shm/"+name, oflag, uint32(perm))
}

// Unlink unlinks the name.
func Unlink(name string) error {
	return os.Remove("/dev/shm/" + name)
}
