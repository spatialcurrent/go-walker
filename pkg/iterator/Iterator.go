// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package iterator

import (
	"io"
	"os"
	"syscall"
)

type Iterator struct {
	fd    int
	data  []byte
	buf   []byte
	names []string
	eof   bool
	err   error
}

func New(fd int) *Iterator {
	return &Iterator{
		fd:    fd,
		data:  make([]byte, 0),
		buf:   make([]byte, os.Getpagesize()),
		names: make([]string, 0),
		eof:   false,
		err:   nil,
	}
}

func (it *Iterator) Reset(fd int) {
	it.fd = fd
	it.data = make([]byte, 0)
	it.buf = make([]byte, os.Getpagesize())
	it.names = make([]string, 0)
	it.eof = false
	it.err = nil
}

func (it *Iterator) Next() (string, error) {
	if it.err != nil {
		return "", it.err
	}
	if it.eof {
		return "", io.EOF
	}
	if len(it.names) > 0 {
		name := it.names[0]
		it.names = it.names[1:]
		return name, nil
	}
	for {
		consumed, count, names := syscall.ParseDirent(it.data, 1, it.names[:0])
		it.data = it.data[consumed:]
		if count == 0 {
			for {
				n, err := syscall.ReadDirent(it.fd, it.buf)
				if err != nil {
					it.err = err
					return "", err
				}
				if n == 0 {
					it.eof = true
					return "", io.EOF
				}
				it.data = append(it.data, it.buf[0:n]...)
				break
			}
			continue
		}
		if count == 1 {
			return names[0], nil
		}
		it.names = names[1:]
		return names[0], nil
	}
	it.eof = true
	return "", io.EOF
}
