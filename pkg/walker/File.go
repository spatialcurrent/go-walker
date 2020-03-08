// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package walker

import (
	"os"
)

type File interface {
	Read(b []byte) (n int, err error)
	Stat() (os.FileInfo, error)
	Close() error
}
