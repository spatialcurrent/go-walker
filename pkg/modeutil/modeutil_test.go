package modeutil

import (
	"os"
	"syscall"
)

func createNamedPipeIfNotExist(p string) {
	_, err := os.Lstat(p)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		panic(err)
	}
	err = syscall.Mkfifo(p, 0660)
	if err != nil {
		panic(err)
	}
}
