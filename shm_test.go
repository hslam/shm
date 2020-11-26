package shm

import (
	"github.com/hslam/mmap"
	"os"
	"testing"
	"time"
)

func TestGetAt(t *testing.T) {
	context := "Hello World"
	done := make(chan struct{})
	go func() {
		shmid, data, err := GetAt(2, 128, IPC_CREATE|0600)
		if err != nil {
			t.Error(err)
		}
		defer Remove(shmid)
		defer Dt(data)
		copy(data, context)
		time.Sleep(time.Second * 2)
		close(done)
	}()
	time.Sleep(time.Second)
	_, data, err := GetAt(2, 128, 0600)
	if err != nil {
		t.Error(err)
	}
	defer Dt(data)
	if context != string(data[:11]) {
		t.Error(context, string(data[:11]))
	}
	<-done
}

func TestGetAtZeroFlag(t *testing.T) {
	context := "Hello World"
	done := make(chan struct{})
	go func() {
		shmid, data, err := GetAt(2, 128, 0)
		if err != nil {
			t.Error(err)
		}
		defer Remove(shmid)
		defer Dt(data)
		copy(data, context)
		time.Sleep(time.Second * 2)
		close(done)
	}()
	time.Sleep(time.Second)
	_, data, err := GetAt(2, 128, 0600)
	if err != nil {
		t.Error(err)
	}
	defer Dt(data)
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
		length := 128
		Ftruncate(fd, int64(length))
		defer Close(fd)
		defer Unlink(name)
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
