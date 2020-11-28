// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package shm

import (
	"github.com/hslam/mmap"
	"os"
	"testing"
	"time"
)

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
		time.Sleep(time.Second * 2)
		close(done)
	}()
	time.Sleep(time.Second)
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
