// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package walker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalker(t *testing.T) {
	w, err := NewWalker(&NewWalkerInput{
		SkipPath:      nil,
		SkipFile:      nil,
		SkipLink:      nil,
		ErrorLinkStat: nil,
		ErrorStat:     nil,
		ErrorWalk:     nil,
		Limit:         NoLimit,
	})
	require.NoError(t, err)
	n, err := w.Walk(context.Background(), []string{"testdata"}, func(ctx context.Context, p string, f File) error {
		return nil
	})
	assert.Equal(t, 11, n)
	assert.NoError(t, err)
}
