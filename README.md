# shm
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/shm)](https://pkg.go.dev/github.com/hslam/shm)
[![Build Status](https://travis-ci.org/hslam/shm.svg?branch=master)](https://travis-ci.org/hslam/shm)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/shm)](https://goreportcard.com/report/github.com/hslam/shm)
[![LICENSE](https://img.shields.io/github/license/hslam/shm.svg?style=flat-square)](https://github.com/hslam/shm/blob/master/LICENSE)

Package shm provides a way to use shared memory.

## Get started

### Install
```
go get github.com/hslam/shm
```
### Import
```
import "github.com/hslam/shm"
```
### Usage
#### SHM GET Example
```go
package main

import (
	"fmt"
	"github.com/hslam/shm"
	"log"
	"time"
)

func main() {
	done := make(chan struct{})
	go func() {
		writer()
		close(done)
	}()
	time.Sleep(time.Second)
	reader()
	<-done
}

func writer() {
	shmid, data, err := shm.GetAt(2, 128, shm.IPC_CREAT|0600)
	if err != nil {
		log.Fatal(err)
	}
	defer shm.Remove(shmid)
	defer shm.Dt(data)
	copy(data, []byte("Hello World"))
	time.Sleep(time.Second * 2)
}

func reader() {
	_, data, err := shm.GetAt(2, 128, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer shm.Dt(data)
	fmt.Println(string(data[:11]))
}
```
#### Output
```
Hello World
```

#### SHM OPEN Example
```go
package main

import (
	"fmt"
	"github.com/hslam/mmap"
	"github.com/hslam/shm"
	"time"
)

func main() {
	done := make(chan struct{})
	go func() {
		writer()
		close(done)
	}()
	time.Sleep(time.Second)
	reader()
	<-done
}

func writer() {
	name := "shared"
	fd, err := shm.Open(name, shm.O_RDWR|shm.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer shm.Close(fd)
	length := 128
	shm.Ftruncate(fd, int64(length))
	defer shm.Unlink(name)
	data, err := mmap.Open(fd, 0, length, mmap.READ|mmap.WRITE)
	if err != nil {
		panic(err)
	}
	defer mmap.Munmap(data)
	copy(data, []byte("Hello World"))
	time.Sleep(time.Second * 2)
}

func reader() {
	fd, err := shm.Open("shared", shm.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer shm.Close(fd)
	data, err := mmap.Open(fd, 0, 128, mmap.READ)
	if err != nil {
		panic(err)
	}
	defer mmap.Munmap(data)
	fmt.Println(string(data[:11]))
}
```
#### Output
```
Hello World
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
shm was written by Meng Huang.


