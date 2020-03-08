// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package walker

import (
	"os"

	"github.com/spatialcurrent/go-lazy/pkg/lazy"
)

type CachedFile struct {
	*lazy.LazyFile
	fileInfo os.FileInfo
	err      error
}

func NewCachedFile(lazyFile *lazy.LazyFile) *CachedFile {
	return &CachedFile{
		LazyFile: lazyFile,
		fileInfo: nil,
		err:      nil,
	}
}

func (f *CachedFile) Stat() (os.FileInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.fileInfo != nil {
		return f.fileInfo, nil
	}
	fileInfo, err := f.LazyFile.Stat()
	if err != nil {
		f.err = err
		return nil, err
	}
	f.fileInfo = fileInfo
	return fileInfo, nil
}
