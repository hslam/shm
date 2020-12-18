# shm
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/shm)](https://pkg.go.dev/github.com/hslam/shm)
[![Build Status](https://github.com/hslam/shm/workflows/build/badge.svg)](https://github.com/hslam/shm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/shm)](https://goreportcard.com/report/github.com/hslam/shm)
[![LICENSE](https://img.shields.io/github/license/hslam/shm.svg?style=flat-square)](https://github.com/hslam/shm/blob/master/LICENSE)

Package shm provides a way to use System V shared memory.

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
**Writer**
```go
package main

import (
	"fmt"
	"github.com/hslam/ftok"
	"github.com/hslam/shm"
	"time"
)

func main() {
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	shmid, data, err := shm.GetAttach(key, 128, shm.IPC_CREAT|0600)
	if err != nil {
		panic(err)
	}
	defer shm.Remove(shmid)
	defer shm.Detach(data)
	context := []byte("Hello World")
	copy(data, context)
	fmt.Println(string(data[:11]))
	time.Sleep(time.Second * 10)
}
```
**Reader**
```go
package main

import (
	"fmt"
	"github.com/hslam/ftok"
	"github.com/hslam/shm"
)

func main() {
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	_, data, err := shm.GetAttach(key, 128, 0600)
	if err != nil {
		panic(err)
	}
	defer shm.Detach(data)
	fmt.Println(string(data[:11]))
}
```
#### Output
```
Hello World
```

#### SHM OPEN Example
**Writer**
```go
package main

import (
	"fmt"
	"github.com/hslam/mmap"
	"github.com/hslam/shm"
	"time"
)

func main() {
	name := "shared"
	fd, err := shm.Open(name, shm.O_RDWR|shm.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer shm.Unlink(name)
	defer shm.Close(fd)
	length := 128
	shm.Ftruncate(fd, int64(length))
	data, err := mmap.Open(fd, 0, length, mmap.READ|mmap.WRITE)
	if err != nil {
		panic(err)
	}
	defer mmap.Munmap(data)
	context := []byte("Hello World")
	copy(data, context)
	fmt.Println(string(data[:11]))
	time.Sleep(time.Second * 10)
}
```
**Reader**
```go
package main

import (
	"fmt"
	"github.com/hslam/mmap"
	"github.com/hslam/shm"
)

func main() {
	name := "shared"
	fd, err := shm.Open(name, shm.O_RDONLY, 0600)
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


