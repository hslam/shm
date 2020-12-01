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
	IPC_CREAT = 01000

	// IPC_EXCL fails if key exists.
	IPC_EXCL = 02000

	// IPC_NOWAIT returns error no wait.
	IPC_NOWAIT = 04000

	// IPC_PRIVATE is private key
	IPC_PRIVATE = 00000

	// SEM_UNDO sets up adjust on exit entry
	SEM_UNDO = 010000

	// IPC_RMID removes identifier
	IPC_RMID = 0
	// IPC_SET sets ipc_perm options.
	IPC_SET = 1
	// IPC_STAT gets ipc_perm options.
	IPC_STAT = 2
)

// Get calls the shmget system call.
func Get(key int, size int, shmFlg int) (int, error) {
	if shmFlg == 0 {
		shmFlg = IPC_CREAT | 0600
	}
	r1, _, errno := syscall.Syscall(SYS_SHMGET, uintptr(key), uintptr(validSize(int64(size))), uintptr(shmFlg))
	shmid := int(r1)
	if shmid < 0 {
		return shmid, syscall.Errno(errno)
	}
	return shmid, nil
}

// At calls the shmat system call.
func At(shmid int, shmFlg int) (uintptr, error) {
	shmaddr, _, errno := syscall.Syscall(SYS_SHMAT, uintptr(shmid), 0, uintptr(shmFlg))
	if int(shmaddr) < 0 {
		return shmaddr, syscall.Errno(errno)
	}
	return shmaddr, nil
}

// Dt calls the shmdt system call.
func Dt(addr uintptr) error {
	r1, _, errno := syscall.Syscall(SYS_SHMDT, addr, 0, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// GetAttach calls the shmget and shmat system call.
func GetAttach(key int, size int, shmFlg int) (int, []byte, error) {
	shmid, err := Get(key, size, shmFlg)
	if err != nil {
		return shmid, nil, err
	}
	shmaddr, err := At(shmid, shmFlg)
	if err != nil {
		Remove(shmid)
		return shmid, nil, err
	}
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{shmaddr, size, size}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	return shmid, b, nil
}

// Detach calls the shmdt system call with []byte b.
func Detach(b []byte) error {
	return Dt(uintptr(unsafe.Pointer(&b[0])))
}

// Remove removes the shm with the given id.
func Remove(shmid int) error {
	r1, _, errno := syscall.Syscall(SYS_SHMCTL, uintptr(shmid), IPC_RMID, 0)
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
