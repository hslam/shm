// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build windows

package shm

import (
	"errors"
	"os"
	"syscall"
)

// Open returns the fd.
func Open(name string, oflag int, perm int) (int, error) {
	file, err := os.OpenFile(name, oflag, os.FileMode(perm))
	if err != nil {
		return 0, err
	}
	return int(file.Fd()), nil
}

// Unlink unlinks the name.
func Unlink(name string) error {
	return os.Remove(name)
}

// GetAt calls the shmget and shmat system call.
func GetAt(key int, size int, shmFlg int) (uintptr, []byte, error) {
	return 0, nil, errors.New("not supported")
}

// Dt calls the shmdt system call.
func Dt(b []byte) error {
	return errors.New("not supported")
}

// Remove removes the shm with the given id.
func Remove(shmid uintptr) error {
	return errors.New("not supported")
}

// Ftruncate changes the file size of the fd.
func Ftruncate(fd int, length int64) (err error) {
	return syscall.Ftruncate(syscall.Handle(fd), length)
}

// Close closes the fd.
func Close(fd int) (err error) {
	return syscall.Close(syscall.Handle(fd))
}
