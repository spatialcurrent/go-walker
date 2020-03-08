// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package pathutil

import (
	"path/filepath"
	"strings"
)

// IsCycle returns true if a link targets itself or a parent directory,
// therefore creating a infinite loop or cycle.
func IsCycle(path string, target string) bool {
	targetCleaned := filepath.Clean(target)
	if target == "." || target == ".." {
		//fmt.Println("Is Cycle:", target)
		return true
	}
	pathCleaned := filepath.Clean(path)
	if strings.HasPrefix(pathCleaned, targetCleaned) {
		//fmt.Println("Is Cycle:", path, target)
		return true
	}
	if strings.HasPrefix(target, "../") {
		//fmt.Println("Is Cycle:", target)
		return true
	}
	return false
}
