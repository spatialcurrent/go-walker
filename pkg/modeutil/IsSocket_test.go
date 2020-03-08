// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package modeutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSocketLstatRegular(t *testing.T) {
	fi, err := os.Lstat("testdata/doc.1.txt")
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.False(t, IsSocket(fi))
}

func TestIsSocketLstatLink(t *testing.T) {
	fi, err := os.Lstat("testdata/doc.2.txt")
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.False(t, IsSocket(fi))
}

func TestIsSocketLstatNamedPipe(t *testing.T) {
	createNamedPipeIfNotExist("testdata/doc.3.txt")
	fi, err := os.Lstat("testdata/doc.3.txt")
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.False(t, IsSocket(fi))
}

func TestIsSocketStatLink(t *testing.T) {
	fi, err := os.Stat("testdata/doc.2.txt")
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.False(t, IsSocket(fi))
}