// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package oserror

import (
	"errors"
	"os"
)

func IsDeniedPermission(err error) bool {
	var pathError *os.PathError
	if errors.As(err, &pathError) && os.IsPermission(pathError) {
		return true
	}
	return false
}
