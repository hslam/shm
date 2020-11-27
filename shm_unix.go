// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin linux dragonfly freebsd netbsd openbsd

package shm

import (
	"syscall"
	"unsafe"
)

const (
	// IPC_CREAT creates if key is nonexistent
	IPC_CREAT = 00001000

	//IPC_RMID removes identifier
	IPC_RMID = 0
)

// GetAt calls the shmget and shmat system call.
func GetAt(key int, size int, shmFlg int) (uintptr, []byte, error) {
	if shmFlg == 0 {
		shmFlg = IPC_CREAT | 0600
	}
	shmid, _, errno := syscall.Syscall(SYS_SHMGET, uintptr(key), uintptr(validSize(int64(size))), uintptr(shmFlg))
	if int(shmid) < 0 {
		return 0, nil, syscall.Errno(errno)
	}
	shmaddr, _, errno := syscall.Syscall(SYS_SHMAT, shmid, 0, uintptr(shmFlg))
	if int(shmaddr) < 0 {
		Remove(shmid)
		return 0, nil, syscall.Errno(errno)
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{shmaddr, size, size}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	return shmid, b, nil
}

// Dt calls the shmdt system call.
func Dt(b []byte) error {
	r1, _, errno := syscall.Syscall(SYS_SHMDT, uintptr(unsafe.Pointer(&b[0])), 0, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Remove removes the shm with the given id.
func Remove(shmid uintptr) error {
	r1, _, errno := syscall.Syscall(SYS_SHMCTL, shmid, IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Ftruncate changes the file size of the fd.
func Ftruncate(fd int, length int64) (err error) {
	return syscall.Ftruncate(fd, length)
}

// Close closes the fd.
func Close(fd int) (err error) {
	return syscall.Close(fd)
}
