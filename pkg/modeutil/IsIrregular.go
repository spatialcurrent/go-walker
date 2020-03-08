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

func IsIrregular(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeIrregular != 0
}
