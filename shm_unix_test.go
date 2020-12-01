// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package shm

import (
	"github.com/hslam/ftok"
	"github.com/hslam/mmap"
	"os"
	"testing"
	"time"
)

func TestGetAt(t *testing.T) {
	context := "Hello World"
	done := make(chan struct{})
	go func() {
		key, err := ftok.Ftok("/tmp", 0x22)
		if err != nil {
			panic(err)
		}
		shmid, data, err := GetAttach(key, 128, IPC_CREAT|0600)
		if err != nil {
			t.Error(err)
		}
		defer Remove(shmid)
		defer Detach(data)
		copy(data, context)
		time.Sleep(time.Millisecond * 200)
		close(done)
	}()
	time.Sleep(time.Millisecond * 100)
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	_, data, err := GetAttach(key, 128, 0600)
	if err != nil {
		t.Error(err)
	}
	defer Detach(data)
	if context != string(data[:11]) {
		t.Error(context, string(data[:11]))
	}
	<-done
}

func TestGetAtZeroFlag(t *testing.T) {
	context := "Hello World"
	done := make(chan struct{})
	go func() {
		key, err := ftok.Ftok("/tmp", 0x22)
		if err != nil {
			panic(err)
		}
		shmid, data, err := GetAttach(key, 128, 0)
		if err != nil {
			t.Error(err)
		}
		defer Remove(shmid)
		defer Detach(data)
		copy(data, context)
		time.Sleep(time.Millisecond * 200)
		close(done)
	}()
	time.Sleep(time.Millisecond * 100)
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	_, data, err := GetAttach(key, 128, 0600)
	if err != nil {
		t.Error(err)
	}
	defer Detach(data)
	if context != string(data[:11]) {
		t.Error(context, string(data[:11]))
	}
	<-done
}

func TestOpen(t *testing.T) {
	context := "Hello World"
	done := make(chan struct{})
	go func() {
		name := "shared"
		fd, err := Open(name, O_RDWR|O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer Unlink(name)
		defer Close(fd)
		length := 128
		Ftruncate(fd, int64(length))
		data, err := mmap.Open(fd, 0, length, mmap.READ|mmap.WRITE)
		if err != nil {
			panic(err)
		}
		defer mmap.Munmap(data)
		copy(data, []byte("Hello World"))
		time.Sleep(time.Millisecond * 200)
		close(done)
	}()
	time.Sleep(time.Millisecond * 100)
	fd, err := Open("shared", O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer Close(fd)
	data, err := mmap.Open(fd, 0, 128, mmap.READ)
	if err != nil {
		panic(err)
	}
	defer mmap.Munmap(data)
	if context != string(data[:11]) {
		t.Error(context, string(data[:11]))
	}
	<-done
}

func TestValidSize(t *testing.T) {
	size := int64(os.Getpagesize()) * 2
	if validSize(size) != size {
		t.Error(validSize(size), size)
	}
}
