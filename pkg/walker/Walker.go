// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package walker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spatialcurrent/go-lazy/pkg/lazy"
	"github.com/spatialcurrent/go-walker/pkg/iterator"
	"github.com/spatialcurrent/go-walker/pkg/modeutil"
	"github.com/spatialcurrent/go-walker/pkg/oserror"
	"github.com/spatialcurrent/go-walker/pkg/pathutil"
)

const (
	NoLimit = -1
)

type cacheFileInfo struct {
	*os.File
	fileInfo os.FileInfo
	err      error
}

func (f *cacheFileInfo) Stat() (os.FileInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.fileInfo != nil {
		return f.fileInfo, nil
	}
	fi, err := f.File.Stat()
	if err != nil {
		f.err = err
		return nil, err
	}
	f.fileInfo = fi
	return fi, nil
}

type Walker struct {
	skipPathFn      func(path string) (bool, error)
	skipFileFn      func(path string, file File) (bool, error)
	skipLinkFn      func(path string, target string) (bool, error)
	errorLinkStatFn func(path string, err error) (bool, bool, error)
	errorStatFn     func(path string, file File, err error) (bool, error)
	errorWalkFn     func(path string, file File, err error) (bool, error)
	limit           int
}

type NewWalkerInput struct {
	SkipPath      func(path string) (bool, error)
	SkipFile      func(path string, file File) (bool, error)
	SkipLink      func(path string, target string) (bool, error)
	ErrorLinkStat func(path string, err error) (bool, bool, error)
	ErrorStat     func(path string, file File, err error) (bool, error)
	ErrorWalk     func(path string, file File, err error) (bool, error)
	Limit         int
}

func NewWalker(input *NewWalkerInput) (*Walker, error) {
	w := &Walker{
		skipPathFn:      input.SkipPath,
		skipFileFn:      input.SkipFile,
		skipLinkFn:      input.SkipLink,
		errorLinkStatFn: input.ErrorLinkStat,
		errorStatFn:     input.ErrorStat,
		errorWalkFn:     input.ErrorWalk,
		limit:           input.Limit,
	}
	return w, nil
}

func (w *Walker) skipPath(path string) (bool, error) {
	if w.skipPathFn == nil {
		return false, nil
	}
	return w.skipPathFn(path)
}

func (w *Walker) skipFile(path string, file File) (bool, error) {
	if w.skipFileFn == nil {
		return false, nil
	}
	return w.skipFileFn(path, file)
}

func (w *Walker) skipLink(path string) (bool, error) {
	if w.skipLinkFn == nil {
		return false, nil
	}
	target, err := os.Readlink(path)
	if err != nil {
		if oserror.IsNotExist(err) {
			return w.skipLinkFn(path, "") // if not exist then pass with target as a blank string.
		}
		if oserror.IsDeniedPermission(err) {
			return w.skipLinkFn(path, "") // if denied permission then pass with target as a blank string.
		}
		return false, fmt.Errorf("error reading link %q: %w", path, err)
	}
	return w.skipLinkFn(path, target)
}

func (w *Walker) handleErrorLinkStat(path string, err error) (bool, bool, error) {
	if w.errorLinkStatFn == nil {
		return false, true, err // skip but do not abort
	}
	return w.errorLinkStatFn(path, err)
}

func (w *Walker) handleErrorStat(path string, file File, err error) (bool, error) {
	if w.errorStatFn == nil {
		return true, err
	}
	return w.errorStatFn(path, file, err)
}

func (w *Walker) handleErrorWalk(path string, file File, err error) (bool, error) {
	if w.errorWalkFn == nil {
		return true, err
	}
	return w.errorWalkFn(path, file, err)
}

func (w *Walker) WalkBucket(ctx context.Context, bucket string, prefix string, f func(ctx context.Context, path string, file File) error) (int, error) {
	return 0, nil
}

func (w *Walker) walkDirectory(ctx context.Context, dir string, d File, fd int, fn func(ctx context.Context, path string, file File) error, limit int) (int, error) {

	if limit == 0 {
		return 0, nil
	}

	it := iterator.New(fd)

	count := 0
	for {

		if limit == 0 {
			break
		}

		name, err := it.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return count, fmt.Errorf("error reading next directory entry for %q: %w", dir, err)
		}

		path := filepath.Join(dir, name)

		skip, err := w.skipPath(path)
		if err != nil {
			return count, err
		}
		if skip {
			continue
		}

		linkFileInfo, err := os.Lstat(path)
		if err != nil {
			abort, skip, err := w.handleErrorLinkStat(path, err)
			if abort {
				return 0, fmt.Errorf("error link stating file %q: %w", path, err)
			}
			if skip {
				continue
			}
			linkFileInfo = nil
		}

		if modeutil.IsNamedPipe(linkFileInfo) || modeutil.IsIrregular(linkFileInfo) || modeutil.IsSocket(linkFileInfo) {
			continue
		}

		file := NewCachedFile(lazy.NewLazyFile(path, os.O_RDONLY, 0))

		skip, err = w.skipFile(path, file)
		if err != nil {
			_ = file.Close()
			return count, err
		}
		if skip {
			continue
		}

		if w.skipLinkFn != nil {
			if linkFileInfo == nil {
				panic(path)
			}
			if modeutil.IsLink(linkFileInfo) {
				skip, err := w.skipLink(path)
				if err != nil {
					_ = file.Close()
					if !oserror.IsNotExist(err) {
						return count, err
					}
				}
				if skip {
					continue
				}
			}
		}

		// increment counter
		count += 1

		// decrement limit
		if limit > 0 {
			limit -= 1
		}

		err = fn(ctx, path, file)
		if err != nil {
			abort, err := w.handleErrorWalk(path, file, err)
			if abort {
				_ = file.Close()
				return count, fmt.Errorf("error calling walk function for file %q: %w", name, err)
			}
		}

		if limit == 0 {
			break
		}

		/* if we received an error when stating before, then do not try again
		if fileInfoError != nil {
			_ = file.Close()
			continue
		}
		*/

		fileInfo, fileInfoError := file.Stat()
		if fileInfoError != nil {
			_ = file.Close()
			abort, err := w.handleErrorStat(path, file, fileInfoError)
			if abort {
				return count, fmt.Errorf("error stating file %q: %w", name, err)
			} else {
				continue
			}
		}

		if fileInfo.IsDir() {
			if linkFileInfo == nil {
				panic(path)
			}
			if modeutil.IsLink(linkFileInfo) {
				target, err := os.Readlink(path)
				if err != nil {
					_ = file.Close()
					return count, fmt.Errorf("error reading link to directory %q: %w", path, err)
				}
				if !pathutil.IsCycle(path, target) {
					fd, err := file.Fd()
					if err != nil {
						return count, fmt.Errorf("error opening directory %q: %w", dir, err)
					}
					n, err := w.walkDirectory(ctx, target, file, int(fd), fn, limit)
					if err != nil {
						_ = file.Close()
						return count, fmt.Errorf("error walking directory %q: %w", path, err)
					}
					count += n
				}
			} else {
				fd, err := file.Fd()
				if err != nil {
					return count, fmt.Errorf("error opening directory %q: %w", dir, err)
				}
				n, err := w.walkDirectory(ctx, path, file, int(fd), fn, limit)
				if err != nil {
					_ = file.Close()
					return count, fmt.Errorf("error walking directory %q: %w", path, err)
				}
				count += n
			}
		}

		_ = file.Close()

	}

	return count, nil
}

func (w *Walker) WalkFileSystem(ctx context.Context, root string, fn func(ctx context.Context, path string, file File) error, limit int) (int, error) {

	if limit == 0 {
		return 0, nil
	}

	if skip, err := w.skipPath(root); err != nil || skip {
		return 0, err
	}

	linkFileInfo, err := os.Lstat(root)
	if err != nil {
		abort, skip, err := w.handleErrorLinkStat(root, err)
		if abort || skip {
			return 0, fmt.Errorf("error link stating root %q: %w", root, err)
		}
	}

	r, err := os.Open(root)
	if err != nil {
		return 0, fmt.Errorf("error opening root %q: %w", root, err)
	}

	cacheFileInfo := &cacheFileInfo{File: r, fileInfo: nil}

	if skip, err := w.skipFile(root, cacheFileInfo); err != nil || skip {
		_ = r.Close()
		return 0, err
	}

	if w.skipLinkFn != nil {
		if modeutil.IsLink(linkFileInfo) {
			if skip, err := w.skipLink(root); err != nil || skip {
				return 0, err
			}
		}
	}

	err = fn(ctx, root, r)
	if err != nil {
		abort, err := w.handleErrorWalk(root, r, err)
		if abort {
			_ = r.Close()
			return 1, fmt.Errorf("error calling walk function for root %q: %w", root, err)
		}
	}

	// decrement limit
	if limit > 0 {
		limit -= 1
		if limit == 0 {
			return 1, nil
		}
	}

	fi, err := cacheFileInfo.Stat()
	if err != nil {
		_, err := w.handleErrorStat(root, r, err)
		if err != nil {
			_ = r.Close()
			return 1, fmt.Errorf("error stating root %q: %w", root, err)
		}
	}

	if fi.IsDir() {
		n, err := w.walkDirectory(ctx, root, cacheFileInfo, int(r.Fd()), fn, limit)
		if err != nil {
			return 1 + n, fmt.Errorf("error walking root directory %q: %w", root, err)
		}
		_ = r.Close()
		return 1 + n, nil
	}

	_ = r.Close()
	return 1, nil
}

func (w *Walker) splitUri(uri string) (string, string) {
	if i := strings.Index(uri, "://"); i != -1 {
		return uri[0:i], uri[i+3:]
	}
	return "", uri
}

func (w *Walker) Walk(ctx context.Context, uris []string, f func(ctx context.Context, path string, file File) error) (int, error) {
	if w.limit == 0 {
		return 0, nil
	}
	count := 0
	for i, uri := range uris {
		scheme, path := w.splitUri(uri)
		if len(scheme) > 0 {
			switch scheme {
			case "s3":
				if i := strings.Index(uri, "/"); i != -1 {
					n, err := w.WalkBucket(ctx, uri[0:i], uri[i+1:], f)
					if err != nil {
						return n + count, fmt.Errorf("error walking bucket %q: %w", uri, err)
					}
					count += n
				} else {
					n, err := w.WalkBucket(ctx, path, "", f)
					if err != nil {
						return n + count, fmt.Errorf("error walking bucket %q: %w", uri, err)
					}
					count += n
				}
			case "file":
				n, err := w.WalkFileSystem(ctx, path, f, w.limit)
				if err != nil {
					return n + count, fmt.Errorf("error walking file system %q: %w", uri, err)
				}
				count += n
			default:
				return 0, fmt.Errorf("error walking uri %q (%d): unknown scheme %q", uri, i, scheme)
			}
		} else {
			n, err := w.WalkFileSystem(ctx, uri, f, w.limit)
			if err != nil {
				return n + count, fmt.Errorf("error walking file system %q: %w", uri, err)
			}
			count += n
		}
	}
	return count, nil
}
