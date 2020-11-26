// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin linux dragonfly freebsd netbsd openbsd

package shm

import (
	"syscall"
	"unsafe"
)

const (
	// IPC_CREATE creates if key is nonexistent
	IPC_CREATE = 00001000

	//IPC_RMID removes identifier
	IPC_RMID = 0

	// O_RDONLY opens the file read-only.
	O_RDONLY int = syscall.O_RDONLY
	// O_WRONLY opens the file write-only.
	O_WRONLY int = syscall.O_WRONLY
	// O_RDWR opens the file read-write.
	O_RDWR int = syscall.O_RDWR

	// O_CREATE creates a new file if none exists.
	O_CREATE int = syscall.O_CREAT
)

// GetAt calls the shmget and shmat system call.
func GetAt(key int, size int, shmFlg int) (uintptr, []byte, error) {
	if shmFlg == 0 {
		shmFlg = IPC_CREATE | 0600
	}
	shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(validSize(int64(size))), uintptr(shmFlg))
	if shmid <= 0 || errno != 0 {
		return 0, nil, syscall.Errno(errno)
	}
	shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if errno != 0 {
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
	_, _, errno := syscall.Syscall(syscall.SYS_SHMDT, uintptr(unsafe.Pointer(&b[0])), 0, 0)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Remove removes the shm with the given id.
func Remove(shmid uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_SHMCTL, shmid, IPC_RMID, 0)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Open returns the shm id.
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

// Ftruncate changes the file size of the fd.
func Ftruncate(fd int, length int64) (err error) {
	return syscall.Ftruncate(fd, length)
}

// Close closes the fd.
func Close(fd int) (err error) {
	return syscall.Close(fd)
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
