// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package modeutil

import (
	"os"
)

func IsNamedPipe(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeNamedPipe != 0
}
